package handler

import (
	"net/http"

	"github.com/franchise-system/backend/internal/entity"
	"github.com/franchise-system/backend/internal/service/chatbot"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ChatHandler struct {
	service *chatbot.Service
}

func NewChatHandler(service *chatbot.Service) *ChatHandler {
	return &ChatHandler{service: service}
}

// Chat handles the main chat endpoint
func (h *ChatHandler) Chat(c *gin.Context) {
	userID := c.MustGet("user_id").(uuid.UUID)

	var req entity.ChatRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Pesan tidak boleh kosong"})
		return
	}

	// Validate message length
	if len(req.Message) > 1000 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Pesan terlalu panjang (maksimal 1000 karakter)"})
		return
	}

	response, err := h.service.Chat(userID, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Maaf, terjadi kesalahan. Silakan coba lagi."})
		return
	}

	c.JSON(http.StatusOK, response)
}

// GetConversations returns all conversations for the user
func (h *ChatHandler) GetConversations(c *gin.Context) {
	userID := c.MustGet("user_id").(uuid.UUID)

	convs, err := h.service.GetConversations(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal memuat percakapan"})
		return
	}

	if convs == nil {
		convs = []entity.ChatConversation{}
	}

	c.JSON(http.StatusOK, gin.H{"conversations": convs})
}

// GetMessages returns all messages for a conversation
func (h *ChatHandler) GetMessages(c *gin.Context) {
	userID := c.MustGet("user_id").(uuid.UUID)

	convID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID percakapan tidak valid"})
		return
	}

	msgs, err := h.service.GetMessages(convID, userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Percakapan tidak ditemukan"})
		return
	}

	if msgs == nil {
		msgs = []entity.ChatMessage{}
	}

	c.JSON(http.StatusOK, gin.H{"messages": msgs})
}

// DeleteConversation removes a conversation
func (h *ChatHandler) DeleteConversation(c *gin.Context) {
	userID := c.MustGet("user_id").(uuid.UUID)

	convID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID percakapan tidak valid"})
		return
	}

	if err := h.service.DeleteConversation(convID, userID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menghapus percakapan"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Percakapan berhasil dihapus"})
}
