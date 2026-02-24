package postgres

import (
	"database/sql"
	"time"

	"github.com/franchise-system/backend/internal/entity"
	"github.com/google/uuid"
	"github.com/lib/pq"
)

type AIKnowledgeRepo struct {
	db *sql.DB
}

func NewAIKnowledgeRepo(db *sql.DB) *AIKnowledgeRepo {
	return &AIKnowledgeRepo{db: db}
}

// ──────────────────────────────────────────────────────────────
// KNOWLEDGE BASE
// ──────────────────────────────────────────────────────────────

// GetAllActiveKnowledge returns all active knowledge entries
func (r *AIKnowledgeRepo) GetAllActiveKnowledge() ([]entity.AIKnowledgeBase, error) {
	rows, err := r.db.Query(`
		SELECT kb.id, kb.category_id, kb.title, kb.slug, kb.content, kb.keywords,
		       kb.priority, kb.is_active, kb.created_at, kb.updated_at,
		       COALESCE(c.name, ''), COALESCE(c.slug, '')
		FROM ai_knowledge_base kb
		LEFT JOIN ai_kb_categories c ON kb.category_id = c.id
		WHERE kb.is_active = true
		ORDER BY kb.priority DESC, kb.created_at ASC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []entity.AIKnowledgeBase
	for rows.Next() {
		var kb entity.AIKnowledgeBase
		var catName, catSlug string
		err := rows.Scan(
			&kb.ID, &kb.CategoryID, &kb.Title, &kb.Slug, &kb.Content, &kb.Keywords,
			&kb.Priority, &kb.IsActive, &kb.CreatedAt, &kb.UpdatedAt,
			&catName, &catSlug,
		)
		if err != nil {
			return nil, err
		}
		if kb.CategoryID != nil {
			kb.Category = &entity.AIKBCategory{Name: catName, Slug: catSlug}
		}
		items = append(items, kb)
	}
	return items, nil
}

// GetKnowledgeByCategory returns knowledge entries filtered by category slug
func (r *AIKnowledgeRepo) GetKnowledgeByCategory(categorySlug string) ([]entity.AIKnowledgeBase, error) {
	rows, err := r.db.Query(`
		SELECT kb.id, kb.category_id, kb.title, kb.slug, kb.content, kb.keywords,
		       kb.priority, kb.is_active, kb.created_at, kb.updated_at
		FROM ai_knowledge_base kb
		JOIN ai_kb_categories c ON kb.category_id = c.id
		WHERE kb.is_active = true AND c.slug = $1
		ORDER BY kb.priority DESC
	`, categorySlug)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []entity.AIKnowledgeBase
	for rows.Next() {
		var kb entity.AIKnowledgeBase
		err := rows.Scan(
			&kb.ID, &kb.CategoryID, &kb.Title, &kb.Slug, &kb.Content, &kb.Keywords,
			&kb.Priority, &kb.IsActive, &kb.CreatedAt, &kb.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		items = append(items, kb)
	}
	return items, nil
}

// SearchKnowledge searches knowledge by keyword matching
func (r *AIKnowledgeRepo) SearchKnowledge(query string) ([]entity.AIKnowledgeBase, error) {
	rows, err := r.db.Query(`
		SELECT kb.id, kb.category_id, kb.title, kb.slug, kb.content, kb.keywords,
		       kb.priority, kb.is_active, kb.created_at, kb.updated_at
		FROM ai_knowledge_base kb
		WHERE kb.is_active = true
		  AND (
		      LOWER(kb.title) LIKE '%' || LOWER($1) || '%'
		      OR LOWER(kb.content) LIKE '%' || LOWER($1) || '%'
		      OR $1 = ANY(kb.keywords)
		  )
		ORDER BY kb.priority DESC
		LIMIT 10
	`, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []entity.AIKnowledgeBase
	for rows.Next() {
		var kb entity.AIKnowledgeBase
		err := rows.Scan(
			&kb.ID, &kb.CategoryID, &kb.Title, &kb.Slug, &kb.Content, &kb.Keywords,
			&kb.Priority, &kb.IsActive, &kb.CreatedAt, &kb.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		items = append(items, kb)
	}
	return items, nil
}

// CreateKnowledge inserts a new knowledge entry
func (r *AIKnowledgeRepo) CreateKnowledge(kb *entity.AIKnowledgeBase) error {
	kb.ID = uuid.New()
	kb.CreatedAt = time.Now()
	kb.UpdatedAt = time.Now()
	_, err := r.db.Exec(`
		INSERT INTO ai_knowledge_base (id, category_id, title, slug, content, keywords, priority, is_active, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
	`, kb.ID, kb.CategoryID, kb.Title, kb.Slug, kb.Content, pq.Array(kb.Keywords),
		kb.Priority, kb.IsActive, kb.CreatedAt, kb.UpdatedAt,
	)
	return err
}

// UpdateKnowledge updates a knowledge entry
func (r *AIKnowledgeRepo) UpdateKnowledge(kb *entity.AIKnowledgeBase) error {
	kb.UpdatedAt = time.Now()
	_, err := r.db.Exec(`
		UPDATE ai_knowledge_base SET
			category_id = $1, title = $2, slug = $3, content = $4, keywords = $5,
			priority = $6, is_active = $7, updated_at = $8
		WHERE id = $9
	`, kb.CategoryID, kb.Title, kb.Slug, kb.Content, pq.Array(kb.Keywords),
		kb.Priority, kb.IsActive, kb.UpdatedAt, kb.ID,
	)
	return err
}

// DeleteKnowledge deletes a knowledge entry
func (r *AIKnowledgeRepo) DeleteKnowledge(id uuid.UUID) error {
	_, err := r.db.Exec(`DELETE FROM ai_knowledge_base WHERE id = $1`, id)
	return err
}

// ──────────────────────────────────────────────────────────────
// CATEGORIES
// ──────────────────────────────────────────────────────────────

func (r *AIKnowledgeRepo) GetAllCategories() ([]entity.AIKBCategory, error) {
	rows, err := r.db.Query(`
		SELECT id, name, slug, description, sort_order, is_active, created_at, updated_at
		FROM ai_kb_categories ORDER BY sort_order ASC, name ASC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var cats []entity.AIKBCategory
	for rows.Next() {
		var c entity.AIKBCategory
		err := rows.Scan(&c.ID, &c.Name, &c.Slug, &c.Description,
			&c.SortOrder, &c.IsActive, &c.CreatedAt, &c.UpdatedAt)
		if err != nil {
			return nil, err
		}
		cats = append(cats, c)
	}
	return cats, nil
}

func (r *AIKnowledgeRepo) CreateCategory(cat *entity.AIKBCategory) error {
	cat.ID = uuid.New()
	cat.CreatedAt = time.Now()
	cat.UpdatedAt = time.Now()
	_, err := r.db.Exec(`
		INSERT INTO ai_kb_categories (id, name, slug, description, sort_order, is_active, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`, cat.ID, cat.Name, cat.Slug, cat.Description, cat.SortOrder, cat.IsActive, cat.CreatedAt, cat.UpdatedAt)
	return err
}

// ──────────────────────────────────────────────────────────────
// SYSTEM PROMPTS
// ──────────────────────────────────────────────────────────────

func (r *AIKnowledgeRepo) GetActiveSystemPrompt() (*entity.AISystemPrompt, error) {
	sp := &entity.AISystemPrompt{}
	err := r.db.QueryRow(`
		SELECT id, name, prompt, is_active, created_at, updated_at
		FROM ai_system_prompts WHERE is_active = true
		ORDER BY updated_at DESC LIMIT 1
	`).Scan(&sp.ID, &sp.Name, &sp.Prompt, &sp.IsActive, &sp.CreatedAt, &sp.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return sp, nil
}

func (r *AIKnowledgeRepo) CreateSystemPrompt(sp *entity.AISystemPrompt) error {
	sp.ID = uuid.New()
	sp.CreatedAt = time.Now()
	sp.UpdatedAt = time.Now()
	_, err := r.db.Exec(`
		INSERT INTO ai_system_prompts (id, name, prompt, is_active, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)
	`, sp.ID, sp.Name, sp.Prompt, sp.IsActive, sp.CreatedAt, sp.UpdatedAt)
	return err
}

func (r *AIKnowledgeRepo) UpdateSystemPrompt(sp *entity.AISystemPrompt) error {
	sp.UpdatedAt = time.Now()
	_, err := r.db.Exec(`
		UPDATE ai_system_prompts SET name = $1, prompt = $2, is_active = $3, updated_at = $4
		WHERE id = $5
	`, sp.Name, sp.Prompt, sp.IsActive, sp.UpdatedAt, sp.ID)
	return err
}

func (r *AIKnowledgeRepo) GetAllSystemPrompts() ([]entity.AISystemPrompt, error) {
	rows, err := r.db.Query(`
		SELECT id, name, prompt, is_active, created_at, updated_at
		FROM ai_system_prompts ORDER BY created_at DESC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var prompts []entity.AISystemPrompt
	for rows.Next() {
		var sp entity.AISystemPrompt
		if err := rows.Scan(&sp.ID, &sp.Name, &sp.Prompt, &sp.IsActive, &sp.CreatedAt, &sp.UpdatedAt); err != nil {
			return nil, err
		}
		prompts = append(prompts, sp)
	}
	return prompts, nil
}

// ──────────────────────────────────────────────────────────────
// AI CONFIG
// ──────────────────────────────────────────────────────────────

func (r *AIKnowledgeRepo) GetConfig(key string) (string, error) {
	var value string
	err := r.db.QueryRow(`SELECT value FROM ai_config WHERE key = $1`, key).Scan(&value)
	return value, err
}

func (r *AIKnowledgeRepo) SetConfig(key, value, description string) error {
	_, err := r.db.Exec(`
		INSERT INTO ai_config (id, key, value, description, created_at, updated_at)
		VALUES ($1, $2, $3, $4, NOW(), NOW())
		ON CONFLICT (key) DO UPDATE SET value = $3, description = $4, updated_at = NOW()
	`, uuid.New(), key, value, description)
	return err
}

func (r *AIKnowledgeRepo) GetAllConfig() ([]entity.AIConfig, error) {
	rows, err := r.db.Query(`
		SELECT id, key, value, COALESCE(description, ''), created_at, updated_at
		FROM ai_config ORDER BY key ASC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var configs []entity.AIConfig
	for rows.Next() {
		var c entity.AIConfig
		if err := rows.Scan(&c.ID, &c.Key, &c.Value, &c.Description, &c.CreatedAt, &c.UpdatedAt); err != nil {
			return nil, err
		}
		configs = append(configs, c)
	}
	return configs, nil
}
