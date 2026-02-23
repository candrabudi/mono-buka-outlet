package router

import (
	"github.com/franchise-system/backend/internal/middleware"
	"github.com/gin-gonic/gin"
)

// RegisterPublicRoutes sets up public routes, legacy auth, and health check.
func RegisterPublicRoutes(r *gin.Engine, h Handlers, jwtSecret string) {
	// ═══════════════════════════════════════════
	// PUBLIC routes (no auth)
	// ═══════════════════════════════════════════
	pub := r.Group("/api/v1/public")
	{
		pub.GET("/outlets", h.Outlet.PublicList)
		pub.GET("/outlets/:id", h.Outlet.PublicDetail)
		pub.GET("/outlet-categories", h.OutletCategory.PublicList)
		pub.POST("/consultation", h.Lead.PublicConsultation)

		// Midtrans webhook (no auth)
		pub.POST("/midtrans/webhook", h.Invoice.MidtransWebhook)
		pub.POST("/midtrans/ebook-webhook", h.Ebook.EbookWebhook)
	}

	// ═══════════════════════════════════════════
	// LEGACY AUTH (kept for backward compat)
	// ═══════════════════════════════════════════
	legacyAuth := r.Group("/api/v1/auth")
	{
		legacyAuth.POST("/login", h.Auth.Login)
		legacyAuth.POST("/register", h.Auth.Register)
	}

	legacyAPI := r.Group("/api/v1")
	legacyAPI.Use(middleware.JWTAuth(jwtSecret))
	{
		legacyAPI.GET("/profile", h.Auth.Profile)
	}

	// ═══════════════════════════════════════════
	// HEALTH
	// ═══════════════════════════════════════════
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok", "service": "bukaoutlet"})
	})
}
