package chatbot

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

// ──────────────────────────────────────────────────────────────
// WEB SEARCH INTEGRATION
// Detects when web search is needed & auto-learns from results
// ──────────────────────────────────────────────────────────────

// needsWebSearch determines if the user's question requires web search
// beyond the internal knowledge base
func (s *Service) needsWebSearch(message string) bool {
	lower := strings.ToLower(message)

	// Direct triggers — user explicitly asking for external info
	directTriggers := []string{
		"cari di internet", "cari di web", "cari online", "search",
		"info terbaru", "berita terbaru", "update terbaru", "tren terbaru",
		"data terbaru", "statistik", "riset", "penelitian",
		"menurut", "berdasarkan data", "sumber",
		"di indonesia", "di dunia", "di asia", "secara global",
		"tahun 2024", "tahun 2025", "tahun 2026",
	}
	for _, trigger := range directTriggers {
		if strings.Contains(lower, trigger) {
			return true
		}
	}

	// Topic-based triggers — franchise/business questions likely needing external info
	topicTriggers := []string{
		"tren franchise", "franchise terbaik", "franchise terlaris",
		"franchise murah", "franchise modal kecil",
		"tips franchise", "tips bisnis", "tips usaha",
		"cara memulai franchise", "cara sukses franchise",
		"keuntungan franchise", "risiko franchise", "kelebihan franchise",
		"regulasi franchise", "hukum franchise", "peraturan franchise",
		"ijin franchise", "izin franchise", "legalitas franchise",
		"franchise makanan", "franchise minuman", "franchise fnb",
		"franchise retail", "franchise jasa",
		"roi franchise", "bep franchise",
		"waralaba", "kemitraan terbaru",
		"contoh franchise", "franchise populer", "franchise sukses",
		"apa itu franchise", "pengertian franchise", "definisi franchise",
		"perbedaan franchise", "jenis franchise", "tipe franchise",
		"strategi pemasaran", "marketing outlet",
		"manajemen outlet", "kelola outlet",
	}
	for _, topic := range topicTriggers {
		if strings.Contains(lower, topic) {
			return true
		}
	}

	// Question patterns that suggest need for broader knowledge
	questionPatterns := []string{
		"bagaimana cara", "apa saja", "berapa biaya",
		"siapa yang", "dimana bisa", "kapan sebaiknya",
		"apakah worth it", "apakah menguntungkan",
		"rekomendasi", "saran untuk", "contoh",
		"perbandingan", "dibandingkan", "vs",
	}
	franchiseContext := []string{
		"franchise", "bisnis", "usaha", "outlet", "kemitraan",
		"waralaba", "investasi", "modal",
	}

	hasQuestion := false
	hasFranchiseContext := false
	for _, q := range questionPatterns {
		if strings.Contains(lower, q) {
			hasQuestion = true
			break
		}
	}
	for _, f := range franchiseContext {
		if strings.Contains(lower, f) {
			hasFranchiseContext = true
			break
		}
	}

	// If it's a franchise-related question, search the web
	if hasQuestion && hasFranchiseContext {
		return true
	}

	return false
}

// autoLearnFromSearch saves useful web search results to the knowledge base
// so the AI learns from previous searches and doesn't need to search again
func (s *Service) autoLearnFromSearch(query string, results []WebSearchResult) {
	if len(results) == 0 {
		return
	}

	ctx := context.Background()

	// Get or create "Web Search Learning" category
	categoryID, err := s.getOrCreateSearchCategory(ctx)
	if err != nil {
		log.Printf("[AI] Failed to get search category: %v", err)
		return
	}

	// Compile results into a single knowledge entry
	var sb strings.Builder
	for i, r := range results {
		if i >= 5 {
			break // Max 5 results per entry
		}
		sb.WriteString(fmt.Sprintf("**%s**\n", r.Title))
		sb.WriteString(fmt.Sprintf("%s\n", r.Snippet))
		if r.URL != "" {
			sb.WriteString(fmt.Sprintf("Sumber: %s\n", r.URL))
		}
		sb.WriteString("\n")
	}

	content := sb.String()
	if content == "" {
		return
	}

	// Check if similar knowledge already exists (by slug match)
	title := fmt.Sprintf("Web: %s", truncateString(query, 80))
	slug := slugify(title)

	var existingID uuid.UUID
	err = s.db.QueryRowContext(ctx,
		`SELECT id FROM ai_knowledge_base WHERE slug = $1 LIMIT 1`, slug).Scan(&existingID)
	if err == nil {
		// Update existing entry with fresh data
		_, err = s.db.ExecContext(ctx,
			`UPDATE ai_knowledge_base SET content = $1, updated_at = $2 WHERE id = $3`,
			content, time.Now(), existingID)
		if err != nil {
			log.Printf("[AI] Failed to update learned knowledge: %v", err)
		} else {
			log.Printf("[AI] 📚 Updated learned knowledge: %s", title)
			s.InvalidateCache()
		}
		return
	}

	// Create new knowledge entry
	keywords := extractKeywords(query)
	kb := &KnowledgeEntry{
		ID:         uuid.New(),
		CategoryID: categoryID,
		Title:      title,
		Slug:       slug,
		Content:    content,
		Keywords:   keywords,
		Priority:   1, // Low priority (internal KB takes precedence)
		IsActive:   true,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	if err := s.saveLearnedKnowledge(ctx, kb); err != nil {
		log.Printf("[AI] Failed to save learned knowledge: %v", err)
	} else {
		log.Printf("[AI] 📚 Saved new learned knowledge: %s", title)
		s.InvalidateCache() // Refresh cache
	}
}

// KnowledgeEntry for auto-learning (simplified)
type KnowledgeEntry struct {
	ID         uuid.UUID
	CategoryID *uuid.UUID
	Title      string
	Slug       string
	Content    string
	Keywords   []string
	Priority   int
	IsActive   bool
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

// getOrCreateSearchCategory ensures the "Web Search" category exists
func (s *Service) getOrCreateSearchCategory(ctx context.Context) (*uuid.UUID, error) {
	// Try to find existing category
	row := s.db.QueryRowContext(ctx,
		`SELECT id FROM ai_kb_categories WHERE slug = 'web-search-learning' LIMIT 1`)
	var catID uuid.UUID
	if err := row.Scan(&catID); err == nil {
		return &catID, nil
	}

	// Create category
	catID = uuid.New()
	_, err := s.db.ExecContext(ctx,
		`INSERT INTO ai_kb_categories (id, name, slug, description, sort_order, is_active, created_at, updated_at)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		 ON CONFLICT (slug) DO NOTHING`,
		catID, "Pembelajaran Web", "web-search-learning",
		"Pengetahuan yang dipelajari otomatis dari pencarian web",
		99, true, time.Now(), time.Now(),
	)
	if err != nil {
		return nil, err
	}

	return &catID, nil
}

// saveLearnedKnowledge inserts a knowledge entry from auto-learning
func (s *Service) saveLearnedKnowledge(ctx context.Context, kb *KnowledgeEntry) error {
	_, err := s.db.ExecContext(ctx,
		`INSERT INTO ai_knowledge_base (id, category_id, title, slug, content, keywords, priority, is_active, created_at, updated_at)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		 ON CONFLICT (slug) DO UPDATE SET content = $5, updated_at = $10`,
		kb.ID, kb.CategoryID, kb.Title, kb.Slug, kb.Content,
		pq.Array(kb.Keywords), kb.Priority, kb.IsActive,
		kb.CreatedAt, kb.UpdatedAt,
	)
	return err
}

// ──────────────────────────────────────────────────────────────
// HELPERS
// ──────────────────────────────────────────────────────────────

func truncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen] + "..."
}

func slugify(s string) string {
	s = strings.ToLower(s)
	s = strings.ReplaceAll(s, " ", "-")
	// Remove non-alphanumeric chars except hyphens
	var result strings.Builder
	for _, r := range s {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') || r == '-' {
			result.WriteRune(r)
		}
	}
	return result.String()
}

func extractKeywords(query string) []string {
	// Simple keyword extraction from the query
	stopwords := map[string]bool{
		"yang": true, "dan": true, "atau": true, "di": true, "ke": true,
		"dari": true, "untuk": true, "dengan": true, "ini": true, "itu": true,
		"ada": true, "apa": true, "bagaimana": true, "berapa": true,
		"siapa": true, "dimana": true, "kapan": true, "mengapa": true,
		"bisa": true, "cara": true, "saya": true, "kita": true,
		"the": true, "a": true, "an": true, "is": true, "are": true,
	}

	words := strings.Fields(strings.ToLower(query))
	var keywords []string
	seen := map[string]bool{}

	for _, w := range words {
		w = strings.Trim(w, ".,!?;:'\"")
		if len(w) < 3 || stopwords[w] || seen[w] {
			continue
		}
		seen[w] = true
		keywords = append(keywords, w)
	}

	if len(keywords) > 10 {
		keywords = keywords[:10]
	}

	return keywords
}
