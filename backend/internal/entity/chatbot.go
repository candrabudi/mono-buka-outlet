package entity

import (
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

// ──────────────────────────────────────────────────────────────
// CHAT ENTITIES
// ──────────────────────────────────────────────────────────────

// ChatMessage represents a single message in a conversation
type ChatMessage struct {
	ID             uuid.UUID `json:"id"`
	ConversationID uuid.UUID `json:"conversation_id"`
	UserID         uuid.UUID `json:"user_id"`
	Role           string    `json:"role"` // "user", "assistant", "system"
	Content        string    `json:"content"`
	Intent         string    `json:"intent,omitempty"`
	TokensUsed     int       `json:"tokens_used,omitempty"`
	CreatedAt      time.Time `json:"created_at"`
}

// ChatConversation represents a chat session
type ChatConversation struct {
	ID        uuid.UUID     `json:"id"`
	UserID    uuid.UUID     `json:"user_id"`
	Title     string        `json:"title"`
	Messages  []ChatMessage `json:"messages,omitempty"`
	CreatedAt time.Time     `json:"created_at"`
	UpdatedAt time.Time     `json:"updated_at"`
}

// ChatRequest is the API request body
type ChatRequest struct {
	Message        string `json:"message" binding:"required"`
	ConversationID string `json:"conversation_id,omitempty"`
}

// ChatResponse is the API response body
type ChatResponse struct {
	Reply          string        `json:"reply"`
	Intent         string        `json:"intent"`
	ConversationID uuid.UUID     `json:"conversation_id"`
	QuickActions   []QuickAction `json:"quick_actions,omitempty"`
	RelatedData    interface{}   `json:"related_data,omitempty"`
}

// QuickAction represents a suggested follow-up action
type QuickAction struct {
	Label  string `json:"label"`
	Action string `json:"action"`
	Icon   string `json:"icon,omitempty"`
}

// ──────────────────────────────────────────────────────────────
// AI KNOWLEDGE BASE ENTITIES
// ──────────────────────────────────────────────────────────────

// AIKBCategory groups knowledge base entries
type AIKBCategory struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Slug        string    `json:"slug"`
	Description string    `json:"description"`
	SortOrder   int       `json:"sort_order"`
	IsActive    bool      `json:"is_active"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// AIKnowledgeBase stores a single piece of knowledge
type AIKnowledgeBase struct {
	ID         uuid.UUID      `json:"id"`
	CategoryID *uuid.UUID     `json:"category_id,omitempty"`
	Title      string         `json:"title"`
	Slug       string         `json:"slug"`
	Content    string         `json:"content"`
	Keywords   pq.StringArray `json:"keywords"`
	Priority   int            `json:"priority"`
	IsActive   bool           `json:"is_active"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`

	// Relation
	Category *AIKBCategory `json:"category,omitempty"`
}

// AISystemPrompt stores configurable system prompts
type AISystemPrompt struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	Prompt    string    `json:"prompt"`
	IsActive  bool      `json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// AIConfig stores AI configuration key-value pairs
type AIConfig struct {
	ID          uuid.UUID `json:"id"`
	Key         string    `json:"key"`
	Value       string    `json:"value"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// Constants
const (
	ChatRoleUser      = "user"
	ChatRoleAssistant = "assistant"
	ChatRoleSystem    = "system"
)
