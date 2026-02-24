package postgres

import (
	"database/sql"
	"time"

	"github.com/franchise-system/backend/internal/entity"
	"github.com/google/uuid"
)

type ChatRepo struct {
	db *sql.DB
}

func NewChatRepo(db *sql.DB) *ChatRepo {
	return &ChatRepo{db: db}
}

// CreateConversation creates a new conversation
func (r *ChatRepo) CreateConversation(userID uuid.UUID, title string) (*entity.ChatConversation, error) {
	conv := &entity.ChatConversation{
		ID:        uuid.New(),
		UserID:    userID,
		Title:     title,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	_, err := r.db.Exec(
		`INSERT INTO chat_conversations (id, user_id, title, created_at, updated_at)
		 VALUES ($1, $2, $3, $4, $5)`,
		conv.ID, conv.UserID, conv.Title, conv.CreatedAt, conv.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return conv, nil
}

// GetConversation returns a conversation by ID (with ownership check)
func (r *ChatRepo) GetConversation(convID, userID uuid.UUID) (*entity.ChatConversation, error) {
	conv := &entity.ChatConversation{}
	err := r.db.QueryRow(
		`SELECT id, user_id, title, created_at, updated_at
		 FROM chat_conversations WHERE id = $1 AND user_id = $2`,
		convID, userID,
	).Scan(&conv.ID, &conv.UserID, &conv.Title, &conv.CreatedAt, &conv.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return conv, nil
}

// ListConversations lists all conversations for a user
func (r *ChatRepo) ListConversations(userID uuid.UUID) ([]entity.ChatConversation, error) {
	rows, err := r.db.Query(
		`SELECT id, user_id, title, created_at, updated_at
		 FROM chat_conversations WHERE user_id = $1 ORDER BY updated_at DESC LIMIT 50`,
		userID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var convs []entity.ChatConversation
	for rows.Next() {
		var c entity.ChatConversation
		if err := rows.Scan(&c.ID, &c.UserID, &c.Title, &c.CreatedAt, &c.UpdatedAt); err != nil {
			return nil, err
		}
		convs = append(convs, c)
	}
	return convs, nil
}

// SaveMessage saves a chat message
func (r *ChatRepo) SaveMessage(msg *entity.ChatMessage) error {
	msg.ID = uuid.New()
	msg.CreatedAt = time.Now()

	_, err := r.db.Exec(
		`INSERT INTO chat_messages (id, conversation_id, user_id, role, content, intent, tokens_used, created_at)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`,
		msg.ID, msg.ConversationID, msg.UserID, msg.Role, msg.Content, msg.Intent, msg.TokensUsed, msg.CreatedAt,
	)

	// Update conversation updated_at
	if err == nil {
		r.db.Exec(
			`UPDATE chat_conversations SET updated_at = $1 WHERE id = $2`,
			time.Now(), msg.ConversationID,
		)
	}

	return err
}

// GetMessages returns messages for a conversation
func (r *ChatRepo) GetMessages(convID, userID uuid.UUID) ([]entity.ChatMessage, error) {
	rows, err := r.db.Query(
		`SELECT id, conversation_id, user_id, role, content, intent, tokens_used, created_at
		 FROM chat_messages WHERE conversation_id = $1 AND user_id = $2
		 ORDER BY created_at ASC LIMIT 200`,
		convID, userID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var msgs []entity.ChatMessage
	for rows.Next() {
		var m entity.ChatMessage
		if err := rows.Scan(&m.ID, &m.ConversationID, &m.UserID, &m.Role, &m.Content, &m.Intent, &m.TokensUsed, &m.CreatedAt); err != nil {
			return nil, err
		}
		msgs = append(msgs, m)
	}
	return msgs, nil
}

// DeleteConversation deletes a conversation and its messages
func (r *ChatRepo) DeleteConversation(convID, userID uuid.UUID) error {
	_, err := r.db.Exec(
		`DELETE FROM chat_conversations WHERE id = $1 AND user_id = $2`,
		convID, userID,
	)
	return err
}

// GetRecentContext returns the last N messages for context
func (r *ChatRepo) GetRecentContext(convID uuid.UUID, limit int) ([]entity.ChatMessage, error) {
	rows, err := r.db.Query(
		`SELECT id, conversation_id, user_id, role, content, intent, tokens_used, created_at
		 FROM chat_messages WHERE conversation_id = $1
		 ORDER BY created_at DESC LIMIT $2`,
		convID, limit,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var msgs []entity.ChatMessage
	for rows.Next() {
		var m entity.ChatMessage
		if err := rows.Scan(&m.ID, &m.ConversationID, &m.UserID, &m.Role, &m.Content, &m.Intent, &m.TokensUsed, &m.CreatedAt); err != nil {
			return nil, err
		}
		msgs = append(msgs, m)
	}

	// Reverse to chronological order
	for i, j := 0, len(msgs)-1; i < j; i, j = i+1, j-1 {
		msgs[i], msgs[j] = msgs[j], msgs[i]
	}
	return msgs, nil
}
