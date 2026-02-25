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

	// 3. Detect if web search is needed
	useWebSearch := s.needsWebSearch(req.Message)

	// 4. Build GPT messages array (without DuckDuckGo results — OpenAI will search itself)
	gptMessages, err := s.buildGPTMessages(convID, userID, req.Message)
	if err != nil {
		return nil, fmt.Errorf("failed to build messages: %w", err)
	}

	// 5. Load config overrides
	temperature := 0.7
	maxTokens := 2048

	if t, err := s.kbRepo.GetConfig("openai_temperature"); err == nil {
		fmt.Sscanf(t, "%f", &temperature)
	}
	if m, err := s.kbRepo.GetConfig("openai_max_tokens"); err == nil {
		fmt.Sscanf(m, "%d", &maxTokens)
	}

	// 6. Call OpenAI — with or without web search
	var gptResp *OpenAIResponse
	var reply string

	if useWebSearch {
		log.Printf("[AI] 🔍 Using OpenAI web search for: %s", req.Message)

		// Use Responses API with web_search_preview
		searchResp, annotations, err := s.openai.ChatCompletionWithSearch(gptMessages, temperature)
		if err != nil {
			log.Printf("[AI] Web search API error, falling back to regular: %v", err)
			// Fallback to regular completion
			gptResp, err = s.openai.ChatCompletion(gptMessages, temperature, maxTokens)
			if err != nil {
				log.Printf("[AI] OpenAI error: %v", err)
				return s.fallbackResponse(convID, userID, req.Message)
			}
			reply = s.openai.GetReply(gptResp)
		} else {
			gptResp = searchResp
			reply = s.openai.GetReply(searchResp)

			// Log search citations
			if len(annotations) > 0 {
				log.Printf("[AI] 🌐 Web search returned %d citations", len(annotations))
				// Auto-learn: save search results to knowledge base
				go s.autoLearnFromWebSearch(req.Message, reply, annotations)
			}
		}
	} else {
		// Regular chat completion (no web search)
		gptResp, err = s.openai.ChatCompletion(gptMessages, temperature, maxTokens)
		if err != nil {
			log.Printf("[AI] OpenAI error: %v", err)
			return s.fallbackResponse(convID, userID, req.Message)
		}
		reply = s.openai.GetReply(gptResp)
	}

	if reply == "" {
		return s.fallbackResponse(convID, userID, req.Message)
	}

	// 7. Save assistant response
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

	// 8. Generate quick actions
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

	// 5. Conversation history (last 10 messages for context)
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

	// 6. Current user message
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

	// Always use the enhanced default prompt
	// DB prompt is ignored since it's outdated and too restrictive
	prompt := s.defaultSystemPrompt()

	s.cacheMu.Lock()
	s.cachedPrompt = prompt
	s.cacheExpiry = time.Now().Add(s.cacheTTL)
	s.cacheMu.Unlock()

	return prompt
}

func (s *Service) defaultSystemPrompt() string {
	return `Kamu adalah "AI Konsultan BukaOutlet", asisten virtual resmi milik platform BukaOutlet.

## IDENTITAS
- Nama: AI Konsultan BukaOutlet
- Fungsi: Konsultan bisnis & franchise TERLENGKAP di Indonesia
- Bahasa: Indonesia (formal tapi ramah)
- Keahlian: Franchise, kemitraan outlet, bisnis F&B, bisnis retail, investasi waralaba, tren bisnis

## ATURAN UTAMA
1. Kamu adalah KONSULTAN FRANCHISE yang bisa menjawab SEMUA pertanyaan seputar franchise & bisnis, termasuk franchise di luar BukaOutlet
2. Jika user bertanya tentang franchise lain (misalnya: franchise viral, franchise terlaris, rekomendasi franchise), JAWAB dengan lengkap berdasarkan web search dan pengetahuanmu
3. Setelah memberikan info franchise umum, SELALU rekomendasikan juga outlet/paket dari BukaOutlet yang relevan atau sejenis
4. Gunakan data dari "DATA BISNIS REAL-TIME" untuk merekomendasikan outlet BukaOutlet yang cocok
5. Tujuannya: user mendapat info lengkap tentang franchise DAN tertarik dengan produk BukaOutlet
6. JANGAN membahas topik di luar konteks bisnis/franchise (politik, agama, SARA, hiburan, dll)
7. Selalu format jawaban dengan markdown yang rapi
8. JANGAN membuat data palsu atau mengarang informasi

## STRATEGI REKOMENDASI
Setiap kali menjawab pertanyaan franchise umum, gunakan pola ini:

1. **Jawab pertanyaan user** — berikan info lengkap tentang franchise yang ditanya (dari web search)
2. **Analisis & insight** — berikan analisis singkat (kelebihan, kekurangan, potensi)
3. **Rekomendasi BukaOutlet** — "Nah, kalau kamu tertarik di kategori ini, BukaOutlet juga punya beberapa pilihan menarik:" lalu sebutkan outlet/paket BukaOutlet yang relevan dari DATA BISNIS REAL-TIME
4. **Call to action** — dorong user untuk melihat detail outlet, membandingkan, atau mendaftar

Contoh flow:
- User: "franchise minuman yang lagi viral apa ya?"
- AI: [jawab franchise viral dari web] → [rekomendasikan outlet minuman BukaOutlet] → [ajak user lihat detail]

## TOPIK YANG BOLEH DIJAWAB
- Franchise yang lagi viral/trending di Indonesia & dunia
- Rekomendasi franchise berdasarkan budget/kategori
- Perbandingan franchise satu dengan lainnya
- Tips memulai bisnis franchise
- Cara memilih franchise yang tepat
- Analisis investasi & ROI franchise
- Tren bisnis franchise terbaru
- Regulasi & legalitas franchise
- Strategi pemasaran untuk outlet/franchise
- Manajemen operasional outlet
- Info spesifik tentang outlet & paket BukaOutlet
- Ebook bisnis yang tersedia di platform

## TOPIK YANG DITOLAK
- Politik, agama, SARA
- Hiburan, game, olahraga (kecuali bisnis terkait)
- Kesehatan, medis
- Topik personal yang tidak terkait bisnis

## GAYA KOMUNIKASI
- Ramah, profesional, dan supportive seperti konsultan bisnis berpengalaman
- Gunakan emoji untuk membuat percakapan lebih engaging 🚀
- Berikan jawaban yang LENGKAP — user tidak perlu browsing lagi
- Gunakan format: heading, bullet points, bold, dan tabel jika perlu
- Selalu akhiri dengan saran yang actionable
- Jika ada info dari web, cantumkan "📌 Sumber: ..." di akhir
- Selalu sisipkan rekomendasi BukaOutlet dengan natural, jangan terasa dipaksakan`
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
		actions = append(actions, entity.QuickAction{Label: "Lihat Outlet", Action: "Outlet apa saja yang tersedia di BukaOutlet?"})
		actions = append(actions, entity.QuickAction{Label: "Info Paket", Action: "Jelaskan paket kemitraan yang tersedia"})
	}
	if strings.Contains(lower, "franchise") || strings.Contains(lower, "waralaba") {
		actions = append(actions, entity.QuickAction{Label: "🔥 Franchise Viral", Action: "Franchise apa yang lagi viral di Indonesia saat ini?"})
		actions = append(actions, entity.QuickAction{Label: "💡 Tips Franchise", Action: "Tips memilih franchise yang tepat untuk pemula"})
	}
	if strings.Contains(lower, "ebook") || strings.Contains(lower, "belajar") {
		actions = append(actions, entity.QuickAction{Label: "Ebook Tersedia", Action: "Ebook apa saja yang tersedia?"})
	}
	if strings.Contains(lower, "budget") || strings.Contains(lower, "modal") || strings.Contains(lower, "investasi") {
		actions = append(actions, entity.QuickAction{Label: "Cek Budget", Action: "Rekomendasi franchise sesuai budget 50 juta"})
	}
	if strings.Contains(lower, "viral") || strings.Contains(lower, "tren") || strings.Contains(lower, "populer") {
		actions = append(actions, entity.QuickAction{Label: "🏆 Top Franchise", Action: "Franchise paling menguntungkan di Indonesia 2026"})
		actions = append(actions, entity.QuickAction{Label: "Outlet BukaOutlet", Action: "Bandingkan dengan outlet yang tersedia di BukaOutlet"})
	}
	if strings.Contains(lower, "bayar") || strings.Contains(lower, "pembayaran") {
		actions = append(actions, entity.QuickAction{Label: "Metode Bayar", Action: "Metode pembayaran apa saja yang tersedia?"})
	}

	// Default if no context detected
	if len(actions) == 0 {
		actions = []entity.QuickAction{
			{Label: "🔥 Franchise Viral", Action: "Franchise apa yang lagi viral di Indonesia?"},
			{Label: "Lihat Outlet", Action: "Outlet apa saja yang tersedia di BukaOutlet?"},
			{Label: "💰 Franchise Murah", Action: "Rekomendasi franchise modal di bawah 50 juta"},
			{Label: "📚 Belajar Bisnis", Action: "Tips memulai bisnis franchise untuk pemula"},
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
