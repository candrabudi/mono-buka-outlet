package chatbot

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/franchise-system/backend/internal/entity"
	"github.com/franchise-system/backend/internal/repository/postgres"
	"github.com/google/uuid"
)

// ──────────────────────────────────────────────────────────────
// AI KONSULTAN SERVICE
// Powered by OpenAI GPT + Database-driven Knowledge Base
// ──────────────────────────────────────────────────────────────

type Service struct {
	chatRepo *postgres.ChatRepo
	kbRepo   *postgres.AIKnowledgeRepo
	db       *sql.DB
	openai   *OpenAIClient

	// Fallback API key from .env (used if DB has no key)
	fallbackAPIKey string

	// Cache for performance
	cacheMu        sync.RWMutex
	cachedKB       string // compiled knowledge base text
	cachedPrompt   string // system prompt
	cachedLiveData string // live business data
	cacheExpiry    time.Time
	cacheTTL       time.Duration
}

func NewService(db *sql.DB, openaiKey string) *Service {
	s := &Service{
		chatRepo:       postgres.NewChatRepo(db),
		kbRepo:         postgres.NewAIKnowledgeRepo(db),
		db:             db,
		openai:         NewOpenAIClient(openaiKey, "gpt-4o"),
		fallbackAPIKey: openaiKey,
		cacheTTL:       5 * time.Minute,
	}

	// Load config from DB (API key, model, etc.)
	go s.reloadConfigFromDB()

	return s
}

// reloadConfigFromDB loads API key, model, and other settings from ai_config table
func (s *Service) reloadConfigFromDB() {
	// API Key — DB takes priority, fallback to .env
	if dbKey, err := s.kbRepo.GetConfig("openai_api_key"); err == nil && dbKey != "" {
		s.openai.apiKey = dbKey
		log.Printf("[AI] API Key loaded from database")
	} else if s.fallbackAPIKey != "" {
		s.openai.apiKey = s.fallbackAPIKey
	}

	// Model
	if model, err := s.kbRepo.GetConfig("openai_model"); err == nil && model != "" {
		s.openai.model = model
		log.Printf("[AI] Model: %s", model)
	}
}

// ──────────────────────────────────────────────────────────────
// MAIN CHAT FLOW
// ──────────────────────────────────────────────────────────────

func (s *Service) Chat(userID uuid.UUID, req entity.ChatRequest) (*entity.ChatResponse, error) {
	// 0. Reload config from DB (picks up changes from admin panel)
	s.reloadConfigFromDB()

	// 1. Get or create conversation
	convID, err := s.resolveConversation(userID, req)
	if err != nil {
		return nil, fmt.Errorf("failed to resolve conversation: %w", err)
	}

	// 2. Save user message
	userMsg := &entity.ChatMessage{
		ConversationID: convID,
		UserID:         userID,
		Role:           entity.ChatRoleUser,
		Content:        req.Message,
	}
	s.chatRepo.SaveMessage(userMsg)

	// 3. Build GPT messages array
	gptMessages, err := s.buildGPTMessages(convID, userID, req.Message)
	if err != nil {
		return nil, fmt.Errorf("failed to build messages: %w", err)
	}

	// 4. Call OpenAI
	temperature := 0.7
	maxTokens := 2048

	// Load config overrides
	if t, err := s.kbRepo.GetConfig("openai_temperature"); err == nil {
		fmt.Sscanf(t, "%f", &temperature)
	}
	if m, err := s.kbRepo.GetConfig("openai_max_tokens"); err == nil {
		fmt.Sscanf(m, "%d", &maxTokens)
	}

	gptResp, err := s.openai.ChatCompletion(gptMessages, temperature, maxTokens)
	if err != nil {
		log.Printf("[AI] OpenAI error: %v", err)
		return s.fallbackResponse(convID, userID, req.Message)
	}

	reply := s.openai.GetReply(gptResp)
	if reply == "" {
		return s.fallbackResponse(convID, userID, req.Message)
	}

	// 5. Save assistant response
	tokensUsed := 0
	if gptResp != nil {
		tokensUsed = gptResp.Usage.TotalTokens
	}
	assistantMsg := &entity.ChatMessage{
		ConversationID: convID,
		UserID:         userID,
		Role:           entity.ChatRoleAssistant,
		Content:        reply,
		TokensUsed:     tokensUsed,
	}
	s.chatRepo.SaveMessage(assistantMsg)

	// 6. Generate quick actions
	quickActions := s.generateQuickActions(reply)

	return &entity.ChatResponse{
		Reply:          reply,
		ConversationID: convID,
		QuickActions:   quickActions,
	}, nil
}

// ──────────────────────────────────────────────────────────────
// BUILD GPT MESSAGES
// ──────────────────────────────────────────────────────────────

func (s *Service) buildGPTMessages(convID, userID uuid.UUID, currentMessage string) ([]OpenAIMessage, error) {
	messages := []OpenAIMessage{}

	// 1. System prompt (from database)
	systemPrompt := s.getSystemPrompt()
	messages = append(messages, OpenAIMessage{
		Role:    "system",
		Content: systemPrompt,
	})

	// 2. Knowledge base context (from database)
	kbContext := s.getKnowledgeContext()
	if kbContext != "" {
		messages = append(messages, OpenAIMessage{
			Role:    "system",
			Content: "## KNOWLEDGE BASE (Referensi Utama)\n\n" + kbContext,
		})
	}

	// 3. Live business data (outlets, packages, ebooks from DB)
	liveData := s.getLiveBusinessData()
	if liveData != "" {
		messages = append(messages, OpenAIMessage{
			Role:    "system",
			Content: "## DATA BISNIS REAL-TIME\n\n" + liveData,
		})
	}

	// 4. Conversation history (last 10 messages for context)
	history, _ := s.chatRepo.GetRecentContext(convID, 10)
	for _, msg := range history {
		role := "user"
		if msg.Role == entity.ChatRoleAssistant {
			role = "assistant"
		}
		messages = append(messages, OpenAIMessage{
			Role:    role,
			Content: msg.Content,
		})
	}

	// 5. Current user message
	messages = append(messages, OpenAIMessage{
		Role:    "user",
		Content: currentMessage,
	})

	return messages, nil
}

// ──────────────────────────────────────────────────────────────
// SYSTEM PROMPT (from Database)
// ──────────────────────────────────────────────────────────────

func (s *Service) getSystemPrompt() string {
	s.cacheMu.RLock()
	if s.cachedPrompt != "" && time.Now().Before(s.cacheExpiry) {
		defer s.cacheMu.RUnlock()
		return s.cachedPrompt
	}
	s.cacheMu.RUnlock()

	// Load from database
	sp, err := s.kbRepo.GetActiveSystemPrompt()
	if err != nil || sp == nil {
		return s.defaultSystemPrompt()
	}

	s.cacheMu.Lock()
	s.cachedPrompt = sp.Prompt
	s.cacheExpiry = time.Now().Add(s.cacheTTL)
	s.cacheMu.Unlock()

	return sp.Prompt
}

func (s *Service) defaultSystemPrompt() string {
	return `Kamu adalah "AI Konsultan BukaOutlet", asisten virtual resmi milik platform BukaOutlet.

## IDENTITAS
- Nama: AI Konsultan BukaOutlet
- Fungsi: Konsultan bisnis khusus kemitraan outlet
- Bahasa: Indonesia (formal tapi ramah)

## ATURAN KETAT
1. HANYA menjawab pertanyaan seputar: kemitraan outlet, bisnis outlet, ebook bisnis, pembayaran, dan topik yang ada di knowledge base
2. Jika ditanya di luar topik bisnis outlet, TOLAK dengan sopan dan arahkan kembali ke topik yang relevan
3. JANGAN pernah memberikan informasi yang tidak ada di knowledge base atau data bisnis
4. JANGAN membahas politik, agama, SARA, atau topik sensitif
5. Gunakan data real dari "DATA BISNIS REAL-TIME" untuk memberikan informasi akurat
6. Jika tidak yakin, sarankan user untuk menghubungi customer service
7. Selalu format jawaban dengan markdown yang rapi (gunakan heading, list, bold)
8. Berikan jawaban yang informatif, terstruktur, dan mudah dipahami
9. Selalu akhiri dengan saran atau pertanyaan lanjutan
10. JANGAN membuat data palsu atau mengarang informasi

## GAYA KOMUNIKASI
- Ramah, profesional, dan supportive
- Gunakan emoji untuk membuat percakapan lebih friendly
- Berikan jawaban yang detail tapi tidak bertele-tele
- Selalu dorong user untuk mengambil action (daftar, lihat outlet, beli ebook, dll)`
}

// ──────────────────────────────────────────────────────────────
// KNOWLEDGE BASE CONTEXT (from Database)
// ──────────────────────────────────────────────────────────────

func (s *Service) getKnowledgeContext() string {
	s.cacheMu.RLock()
	if s.cachedKB != "" && time.Now().Before(s.cacheExpiry) {
		defer s.cacheMu.RUnlock()
		return s.cachedKB
	}
	s.cacheMu.RUnlock()

	// Load all active knowledge from database
	kbItems, err := s.kbRepo.GetAllActiveKnowledge()
	if err != nil || len(kbItems) == 0 {
		return ""
	}

	var sb strings.Builder
	currentCategory := ""
	for _, kb := range kbItems {
		catName := "Umum"
		if kb.Category != nil && kb.Category.Name != "" {
			catName = kb.Category.Name
		}

		if catName != currentCategory {
			sb.WriteString(fmt.Sprintf("\n### %s\n\n", catName))
			currentCategory = catName
		}

		sb.WriteString(fmt.Sprintf("#### %s\n%s\n\n", kb.Title, kb.Content))
	}

	result := sb.String()

	s.cacheMu.Lock()
	s.cachedKB = result
	s.cacheExpiry = time.Now().Add(s.cacheTTL)
	s.cacheMu.Unlock()

	return result
}

// ──────────────────────────────────────────────────────────────
// LIVE BUSINESS DATA (from existing tables)
// ──────────────────────────────────────────────────────────────

func (s *Service) getLiveBusinessData() string {
	s.cacheMu.RLock()
	if s.cachedLiveData != "" && time.Now().Before(s.cacheExpiry) {
		defer s.cacheMu.RUnlock()
		return s.cachedLiveData
	}
	s.cacheMu.RUnlock()

	ctx := context.Background()
	var sb strings.Builder

	// ── Outlets ──
	rows, err := s.db.QueryContext(ctx, `
		SELECT o.name, COALESCE(o.short_description,''), o.minimum_investment,
		       o.maximum_investment, o.profit_sharing_percentage,
		       COALESCE(o.estimated_roi,''), COALESCE(o.location_requirement,''),
		       COALESCE(o.contact_whatsapp,''), COALESCE(oc.name,'') as category_name
		FROM outlets o
		LEFT JOIN outlet_categories oc ON o.category_id = oc.id
		WHERE o.deleted_at IS NULL AND o.is_active = true
		ORDER BY o.is_featured DESC, o.created_at DESC LIMIT 20
	`)
	if err == nil {
		sb.WriteString("### Daftar Outlet Tersedia\n\n")
		for rows.Next() {
			var name, desc, roi, locReq, wa, catName string
			var minInv float64
			var maxInv *float64
			var profitShare float64
			rows.Scan(&name, &desc, &minInv, &maxInv, &profitShare, &roi, &locReq, &wa, &catName)

			sb.WriteString(fmt.Sprintf("**%s** (Kategori: %s)\n", name, catName))
			sb.WriteString(fmt.Sprintf("- Investasi: Rp %s", formatIDR(minInv)))
			if maxInv != nil {
				sb.WriteString(fmt.Sprintf(" - Rp %s", formatIDR(*maxInv)))
			}
			sb.WriteString(fmt.Sprintf("\n- Profit Sharing: %.0f%%\n", profitShare))
			sb.WriteString(fmt.Sprintf("- Estimasi ROI: %s\n", roi))
			if desc != "" {
				sb.WriteString(fmt.Sprintf("- Deskripsi: %s\n", desc))
			}
			if locReq != "" {
				sb.WriteString(fmt.Sprintf("- Kebutuhan Lokasi: %s\n", locReq))
			}
			if wa != "" {
				sb.WriteString(fmt.Sprintf("- WhatsApp: %s\n", wa))
			}
			sb.WriteString("\n")
		}
		rows.Close()
	}

	// ── Packages ──
	rows2, err := s.db.QueryContext(ctx, `
		SELECT p.name, p.price, p.minimum_dp, COALESCE(p.duration,''),
		       COALESCE(p.estimated_bep,''), COALESCE(p.net_profit,''),
		       COALESCE(p.description,''), p.benefits,
		       COALESCE(o.name,'') as outlet_name
		FROM outlet_packages p
		LEFT JOIN outlets o ON p.outlet_id = o.id
		WHERE p.is_active = true
		ORDER BY p.sort_order ASC, p.price ASC LIMIT 30
	`)
	if err == nil {
		sb.WriteString("### Paket Kemitraan Tersedia\n\n")
		for rows2.Next() {
			var name, duration, bep, netProfit, desc, outletName string
			var price, minDP int64
			var benefits []string
			rows2.Scan(&name, &price, &minDP, &duration, &bep, &netProfit, &desc, (*pqStringArray)(&benefits), &outletName)

			sb.WriteString(fmt.Sprintf("**%s** (Outlet: %s)\n", name, outletName))
			sb.WriteString(fmt.Sprintf("- Harga: Rp %s\n", formatIDR(float64(price))))
			sb.WriteString(fmt.Sprintf("- DP Minimum: Rp %s\n", formatIDR(float64(minDP))))
			if duration != "" {
				sb.WriteString(fmt.Sprintf("- Durasi: %s\n", duration))
			}
			if bep != "" {
				sb.WriteString(fmt.Sprintf("- Estimasi BEP: %s\n", bep))
			}
			if netProfit != "" {
				sb.WriteString(fmt.Sprintf("- Net Profit: %s\n", netProfit))
			}
			if len(benefits) > 0 {
				sb.WriteString("- Benefits: " + strings.Join(benefits, ", ") + "\n")
			}
			sb.WriteString("\n")
		}
		rows2.Close()
	}

	// ── Ebooks ──
	rows3, err := s.db.QueryContext(ctx, `
		SELECT title, COALESCE(author,''), price, COALESCE(description,''), total_sold
		FROM ebooks WHERE is_active = true
		ORDER BY total_sold DESC, created_at DESC LIMIT 20
	`)
	if err == nil {
		sb.WriteString("### Ebook Tersedia\n\n")
		for rows3.Next() {
			var title, author, desc string
			var price int64
			var totalSold int
			rows3.Scan(&title, &author, &price, &desc, &totalSold)

			priceText := "Gratis"
			if price > 0 {
				priceText = fmt.Sprintf("Rp %s", formatIDR(float64(price)))
			}
			sb.WriteString(fmt.Sprintf("**%s** oleh %s — %s (Terjual: %d)\n", title, author, priceText, totalSold))
			if desc != "" {
				truncated := desc
				if len(truncated) > 150 {
					truncated = truncated[:150] + "..."
				}
				sb.WriteString(fmt.Sprintf("  %s\n", truncated))
			}
			sb.WriteString("\n")
		}
		rows3.Close()
	}

	// ── Categories ──
	rows4, err := s.db.QueryContext(ctx, `
		SELECT name, COALESCE(description,'') FROM outlet_categories
		WHERE is_active = true ORDER BY sort_order ASC
	`)
	if err == nil {
		sb.WriteString("### Kategori Outlet\n\n")
		for rows4.Next() {
			var name, desc string
			rows4.Scan(&name, &desc)
			sb.WriteString(fmt.Sprintf("- **%s**: %s\n", name, desc))
		}
		sb.WriteString("\n")
		rows4.Close()
	}

	result := sb.String()

	s.cacheMu.Lock()
	s.cachedLiveData = result
	s.cacheExpiry = time.Now().Add(s.cacheTTL)
	s.cacheMu.Unlock()

	return result
}

// ──────────────────────────────────────────────────────────────
// CONVERSATION MANAGEMENT
// ──────────────────────────────────────────────────────────────

func (s *Service) resolveConversation(userID uuid.UUID, req entity.ChatRequest) (uuid.UUID, error) {
	if req.ConversationID != "" {
		parsed, err := uuid.Parse(req.ConversationID)
		if err == nil {
			conv, err := s.chatRepo.GetConversation(parsed, userID)
			if err == nil && conv != nil {
				return parsed, nil
			}
		}
	}

	// Create new conversation
	title := req.Message
	if len(title) > 50 {
		title = title[:50] + "..."
	}
	conv, err := s.chatRepo.CreateConversation(userID, title)
	if err != nil {
		return uuid.Nil, err
	}
	return conv.ID, nil
}

func (s *Service) GetConversations(userID uuid.UUID) ([]entity.ChatConversation, error) {
	return s.chatRepo.ListConversations(userID)
}

func (s *Service) GetMessages(convID, userID uuid.UUID) ([]entity.ChatMessage, error) {
	_, err := s.chatRepo.GetConversation(convID, userID)
	if err != nil {
		return nil, err
	}
	return s.chatRepo.GetMessages(convID, userID)
}

func (s *Service) DeleteConversation(convID, userID uuid.UUID) error {
	return s.chatRepo.DeleteConversation(convID, userID)
}

// ──────────────────────────────────────────────────────────────
// QUICK ACTIONS GENERATOR
// ──────────────────────────────────────────────────────────────

func (s *Service) generateQuickActions(reply string) []entity.QuickAction {
	lower := strings.ToLower(reply)

	actions := []entity.QuickAction{}

	// Detect context and suggest relevant follow-ups
	if strings.Contains(lower, "mitra") || strings.Contains(lower, "kemitraan") {
		actions = append(actions, entity.QuickAction{Label: "Alur Kemitraan", Action: "Bagaimana alur lengkap menjadi mitra?"})
		actions = append(actions, entity.QuickAction{Label: "Biaya Investasi", Action: "Berapa biaya investasi menjadi mitra?"})
	}
	if strings.Contains(lower, "outlet") || strings.Contains(lower, "paket") {
		actions = append(actions, entity.QuickAction{Label: "Lihat Outlet", Action: "Outlet apa saja yang tersedia?"})
		actions = append(actions, entity.QuickAction{Label: "Info Paket", Action: "Jelaskan paket kemitraan yang tersedia"})
	}
	if strings.Contains(lower, "ebook") || strings.Contains(lower, "belajar") {
		actions = append(actions, entity.QuickAction{Label: "Ebook Tersedia", Action: "Ebook apa saja yang tersedia?"})
	}
	if strings.Contains(lower, "budget") || strings.Contains(lower, "modal") || strings.Contains(lower, "investasi") {
		actions = append(actions, entity.QuickAction{Label: "Cek Budget", Action: "Rekomendasi outlet sesuai budget saya"})
	}
	if strings.Contains(lower, "bayar") || strings.Contains(lower, "pembayaran") {
		actions = append(actions, entity.QuickAction{Label: "Metode Bayar", Action: "Metode pembayaran apa saja yang tersedia?"})
	}

	// Default if no context detected
	if len(actions) == 0 {
		actions = []entity.QuickAction{
			{Label: "Info Mitra", Action: "Bagaimana cara menjadi mitra?"},
			{Label: "Lihat Outlet", Action: "Outlet apa saja yang tersedia?"},
			{Label: "Belajar Bisnis", Action: "Saya ingin belajar bisnis outlet"},
		}
	}

	// Limit to 4 actions max
	if len(actions) > 4 {
		actions = actions[:4]
	}

	return actions
}

// ──────────────────────────────────────────────────────────────
// FALLBACK RESPONSE (when OpenAI fails)
// ──────────────────────────────────────────────────────────────

func (s *Service) fallbackResponse(convID, userID uuid.UUID, message string) (*entity.ChatResponse, error) {
	reply := `Mohon maaf, saat ini AI Konsultan sedang mengalami gangguan sementara.

Berikut yang bisa Anda lakukan:
1. **Coba lagi** dalam beberapa saat
2. **Hubungi customer service** kami via WhatsApp
3. **Jelajahi langsung** menu Outlet, Ebook, atau Pengajuan di sidebar

Kami mohon maaf atas ketidaknyamanannya.`

	assistantMsg := &entity.ChatMessage{
		ConversationID: convID,
		UserID:         userID,
		Role:           entity.ChatRoleAssistant,
		Content:        reply,
	}
	s.chatRepo.SaveMessage(assistantMsg)

	return &entity.ChatResponse{
		Reply:          reply,
		ConversationID: convID,
		QuickActions: []entity.QuickAction{
			{Label: "Coba Lagi", Action: message},
			{Label: "Lihat Outlet", Action: "Outlet apa saja yang tersedia?"},
		},
	}, nil
}

// ──────────────────────────────────────────────────────────────
// CACHE INVALIDATION
// ──────────────────────────────────────────────────────────────

// InvalidateCache forces reload of knowledge base and system prompt
func (s *Service) InvalidateCache() {
	s.cacheMu.Lock()
	s.cachedKB = ""
	s.cachedPrompt = ""
	s.cachedLiveData = ""
	s.cacheExpiry = time.Time{}
	s.cacheMu.Unlock()
}

// ──────────────────────────────────────────────────────────────
// HELPERS
// ──────────────────────────────────────────────────────────────

func formatIDR(amount float64) string {
	s := fmt.Sprintf("%.0f", amount)
	n := len(s)
	if n <= 3 {
		return s
	}
	var result strings.Builder
	for i, ch := range s {
		if i > 0 && (n-i)%3 == 0 {
			result.WriteRune('.')
		}
		result.WriteRune(ch)
	}
	return result.String()
}

// pqStringArray is a helper for scanning PostgreSQL text arrays
type pqStringArray []string

func (a *pqStringArray) Scan(src interface{}) error {
	if src == nil {
		*a = []string{}
		return nil
	}
	switch v := src.(type) {
	case []byte:
		str := string(v)
		str = strings.Trim(str, "{}")
		if str == "" {
			*a = []string{}
			return nil
		}
		*a = strings.Split(str, ",")
		return nil
	default:
		return fmt.Errorf("unsupported type: %T", src)
	}
}
