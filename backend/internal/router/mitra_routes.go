package router

import (
	"time"

	"github.com/franchise-system/backend/internal/middleware"
	"github.com/gin-gonic/gin"
)

// RegisterMitraRoutes sets up all /api/v1/mitra/* routes for the mitra portal.
func RegisterMitraRoutes(r *gin.Engine, h Handlers, jwtSecret string) {
	mitra := r.Group("/api/v1/mitra")

	// Rate limiter: 5 attempts per minute, 2 minute lockout
	authLimiter := middleware.NewRateLimiter(5, 1*time.Minute, 2*time.Minute)

	{
		// Auth (public) — with rate limiting
		mitraAuth := mitra.Group("/auth")
		mitraAuth.Use(authLimiter.Middleware())
		{
			mitraAuth.POST("/login", h.AdminAuth.Login)          // Step 1: email+pass → OTP
			mitraAuth.POST("/verify-otp", h.AdminAuth.VerifyOTP) // Step 2: OTP → JWT
			mitraAuth.POST("/resend-otp", h.AdminAuth.ResendOTP)
			mitraAuth.POST("/register", h.Auth.MitraRegister)
		}

		// Protected mitra routes
		mitraProtected := mitra.Group("")
		mitraProtected.Use(middleware.JWTAuth(jwtSecret))
		{
			mitraProtected.GET("/profile", h.Auth.Profile)
			mitraProtected.PUT("/profile", h.Auth.UpdateMyProfile)
			mitraProtected.POST("/change-password", h.Auth.ChangeMyPassword)

			// Outlets browsing
			mitraProtected.GET("/outlets", h.Mitra.ListOutlets)
			mitraProtected.GET("/outlets/:id", h.Mitra.GetOutlet)
			mitraProtected.GET("/outlets/:id/packages", h.Mitra.GetOutletPackages)

			// Partnership Applications
			mitraProtected.POST("/applications", h.Mitra.Apply)
			mitraProtected.GET("/applications", h.Mitra.MyApplications)
			mitraProtected.GET("/applications/:id", h.Mitra.GetApplication)
			mitraProtected.POST("/applications/:id/cancel", h.Mitra.CancelApplication)

			// Partnership (own)
			mitraProtected.GET("/partnership", h.Partnership.GetByMitra)

			// Invoices (own)
			mitraProtected.GET("/invoices", h.Invoice.GetByMitra)
			mitraProtected.GET("/invoices/:id", h.Invoice.GetByID)
			mitraProtected.GET("/midtrans/client-key", h.Invoice.GetMidtransClientKey)

			// Agreements (own)
			mitraProtected.GET("/agreements", h.Agreement.GetByMitra)

			// Locations (own)
			mitraProtected.GET("/locations", h.LocationSub.GetByMitra)
		}
	}
}
