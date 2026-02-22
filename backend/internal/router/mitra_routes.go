package router

import (
	"github.com/franchise-system/backend/internal/middleware"
	"github.com/gin-gonic/gin"
)

// RegisterMitraRoutes sets up all /api/v1/mitra/* routes for the mitra portal.
func RegisterMitraRoutes(r *gin.Engine, h Handlers, jwtSecret string) {
	mitra := r.Group("/api/v1/mitra")
	{
		// Auth (public)
		mitraAuth := mitra.Group("/auth")
		{
			mitraAuth.POST("/login", h.Auth.Login)
			mitraAuth.POST("/register", h.Auth.MitraRegister)
		}

		// Protected mitra routes
		mitraProtected := mitra.Group("")
		mitraProtected.Use(middleware.JWTAuth(jwtSecret))
		{
			mitraProtected.GET("/profile", h.Auth.Profile)

			// Outlets browsing
			mitraProtected.GET("/outlets", h.Mitra.ListOutlets)
			mitraProtected.GET("/outlets/:id", h.Mitra.GetOutlet)
			mitraProtected.GET("/outlets/:id/packages", h.Mitra.GetOutletPackages)

			// Partnership Applications
			mitraProtected.POST("/applications", h.Mitra.Apply)
			mitraProtected.GET("/applications", h.Mitra.MyApplications)
			mitraProtected.GET("/applications/:id", h.Mitra.GetApplication)

			// Partnership (own)
			mitraProtected.GET("/partnership", h.Partnership.GetByMitra)

			// Invoices (own)
			mitraProtected.GET("/invoices", h.Invoice.GetByMitra)

			// Agreements (own)
			mitraProtected.GET("/agreements", h.Agreement.GetByMitra)

			// Locations (own)
			mitraProtected.GET("/locations", h.LocationSub.GetByMitra)
		}
	}
}
