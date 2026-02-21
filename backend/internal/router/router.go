package router

import (
	"github.com/franchise-system/backend/internal/entity"
	"github.com/franchise-system/backend/internal/handler"
	"github.com/franchise-system/backend/internal/middleware"
	"github.com/gin-gonic/gin"
)

type Handlers struct {
	Auth      *handler.AuthHandler
	AdminAuth *handler.AdminAuthHandler

	Outlet         *handler.OutletHandler
	OutletCategory *handler.OutletCategoryHandler
	Lead           *handler.LeadHandler
	Partnership    *handler.PartnershipHandler
	Payment        *handler.PaymentHandler
	Agreement      *handler.AgreementHandler
	Revenue        *handler.RevenueHandler
	Dashboard      *handler.DashboardHandler
	Upload         *handler.UploadHandler
	OutletPackage  *handler.OutletPackageHandler
	Meeting        *handler.MeetingHandler
	Setting        *handler.SettingHandler
	Invoice        *handler.InvoiceHandler
	LocationSub    *handler.LocationSubmissionHandler
}

func Setup(r *gin.Engine, h Handlers, jwtSecret, corsOrigins string) {
	r.Use(middleware.CORS(corsOrigins))

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
	}

	// ═══════════════════════════════════════════
	// ADMIN PANEL — /api/v1/admin/...
	// ═══════════════════════════════════════════
	admin := r.Group("/api/v1/admin")
	{
		// Auth (no JWT required)
		adminAuth := admin.Group("/auth")
		{
			adminAuth.POST("/login", h.AdminAuth.Login)          // Step 1: email+pass → OTP
			adminAuth.POST("/verify-otp", h.AdminAuth.VerifyOTP) // Step 2: OTP → JWT
			adminAuth.POST("/resend-otp", h.AdminAuth.ResendOTP) // Resend OTP
		}

		// Protected admin routes
		adminProtected := admin.Group("")
		adminProtected.Use(middleware.JWTAuth(jwtSecret))
		adminProtected.Use(middleware.AdminOnly())
		{
			// Profile
			adminProtected.GET("/profile", h.AdminAuth.Profile)

			// Dashboard — all admin roles
			adminProtected.GET("/dashboard", h.Dashboard.GetStats)

			// Upload — all admin roles
			adminProtected.POST("/upload", h.Upload.Upload)

			// Outlets — master, admin
			outlets := adminProtected.Group("/outlets")
			outlets.Use(middleware.RoleAuth(entity.RoleMaster, entity.RoleAdmin))
			{
				outlets.POST("", h.Outlet.Create)
				outlets.GET("", h.Outlet.GetAll)
				outlets.GET("/:id", h.Outlet.GetByID)
				outlets.PUT("/:id", h.Outlet.Update)
				outlets.DELETE("/:id", h.Outlet.Delete)
				outlets.PATCH("/:id/toggle", h.Outlet.ToggleActive)
				outlets.PATCH("/:id/featured", h.Outlet.ToggleFeatured)
			}

			// Outlet Categories — master, admin
			outletCats := adminProtected.Group("/outlet-categories")
			outletCats.Use(middleware.RoleAuth(entity.RoleMaster, entity.RoleAdmin))
			{
				outletCats.POST("", h.OutletCategory.Create)
				outletCats.GET("", h.OutletCategory.GetAll)
				outletCats.GET("/:id", h.OutletCategory.GetByID)
				outletCats.PUT("/:id", h.OutletCategory.Update)
				outletCats.PATCH("/:id/toggle", h.OutletCategory.ToggleActive)
				outletCats.DELETE("/:id", h.OutletCategory.Delete)
			}

			// Outlet Packages — master, admin
			outletPkgs := adminProtected.Group("/outlet-packages")
			outletPkgs.Use(middleware.RoleAuth(entity.RoleMaster, entity.RoleAdmin))
			{
				outletPkgs.POST("", h.OutletPackage.Create)
				outletPkgs.GET("/outlet/:outletId", h.OutletPackage.GetByOutletID)
				outletPkgs.GET("/:id", h.OutletPackage.GetByID)
				outletPkgs.PUT("/:id", h.OutletPackage.Update)
				outletPkgs.DELETE("/:id", h.OutletPackage.Delete)
			}

			// Leads — master, admin
			leads := adminProtected.Group("/leads")
			leads.Use(middleware.RoleAuth(entity.RoleMaster, entity.RoleAdmin))
			{
				leads.POST("", h.Lead.Create)
				leads.GET("", h.Lead.GetAll)
				leads.GET("/kanban", h.Lead.GetKanban)
				leads.GET("/:id", h.Lead.GetByID)
				leads.PUT("/:id", h.Lead.Update)
				leads.PATCH("/:id/status", h.Lead.UpdateStatus)
				leads.DELETE("/:id", h.Lead.Delete)
			}

			// Partnerships — master, admin (read), finance (can view)
			partnerships := adminProtected.Group("/partnerships")
			{
				partnerships.POST("", middleware.RoleAuth(entity.RoleMaster, entity.RoleAdmin), h.Partnership.Create)
				partnerships.GET("", middleware.RoleAuth(
					entity.RoleMaster, entity.RoleAdmin, entity.RoleFinance,
				), h.Partnership.GetAll)
				partnerships.GET("/:id", h.Partnership.GetByID)
			}

			// Payments — master, finance
			payments := adminProtected.Group("/payments")
			{
				payments.POST("", middleware.RoleAuth(entity.RoleMaster, entity.RoleFinance), h.Payment.Create)
				payments.PATCH("/:id/verify", middleware.RoleAuth(entity.RoleMaster, entity.RoleFinance), h.Payment.Verify)
				payments.GET("/partnership/:partnership_id", h.Payment.GetByPartnership)
			}

			// Agreements — master, admin
			agreements := adminProtected.Group("/agreements")
			{
				agreements.POST("", middleware.RoleAuth(entity.RoleMaster, entity.RoleAdmin), h.Agreement.Create)
				agreements.PATCH("/:id/sign", middleware.RoleAuth(entity.RoleMaster, entity.RoleAdmin), h.Agreement.Sign)
				agreements.GET("/partnership/:partnership_id", h.Agreement.GetByPartnership)
			}

			// Revenue — master, finance
			revenue := adminProtected.Group("/revenues")
			{
				revenue.POST("", middleware.RoleAuth(entity.RoleMaster, entity.RoleFinance), h.Revenue.Create)
				revenue.GET("/partnership/:partnership_id", h.Revenue.GetByPartnership)
			}

			// Users — master, admin
			users := adminProtected.Group("/users")
			users.Use(middleware.RoleAuth(entity.RoleMaster, entity.RoleAdmin))
			{
				users.GET("", h.Auth.ListUsers)
				users.POST("", h.Auth.Register)
				users.GET("/:id", h.Auth.GetUser)
				users.PUT("/:id", h.Auth.UpdateUser)
				users.DELETE("/:id", h.Auth.DeleteUser)
				users.PATCH("/:id/toggle", h.Auth.ToggleUserActive)
			}

			// Meetings — master, admin
			meetings := adminProtected.Group("/meetings")
			meetings.Use(middleware.RoleAuth(entity.RoleMaster, entity.RoleAdmin))
			{
				meetings.POST("", h.Meeting.Create)
				meetings.GET("", h.Meeting.GetAll)
				meetings.GET("/:id", h.Meeting.GetByID)
				meetings.PUT("/:id", h.Meeting.Update)
				meetings.DELETE("/:id", h.Meeting.Delete)

				// Sub-resources
				meetings.POST("/:id/participants", h.Meeting.AddParticipant)
				meetings.DELETE("/:id/participants/:participantId", h.Meeting.DeleteParticipant)
				meetings.POST("/:id/notes", h.Meeting.SaveNotes)
				meetings.POST("/:id/action-plans", h.Meeting.AddActionPlan)
				meetings.PUT("/:id/action-plans/:actionId", h.Meeting.UpdateActionPlan)
				meetings.DELETE("/:id/action-plans/:actionId", h.Meeting.DeleteActionPlan)
				meetings.POST("/:id/upload", h.Meeting.UploadFile)
				meetings.DELETE("/:id/files/:fileId", h.Meeting.DeleteFile)
			}

			// System Settings — master only
			settings := adminProtected.Group("/settings")
			settings.Use(middleware.RoleAuth(entity.RoleMaster, entity.RoleAdmin))
			{
				settings.GET("", h.Setting.GetAll)
				settings.GET("/:key", h.Setting.GetByKey)
				settings.PUT("", h.Setting.BulkUpdate)
			}

			// Invoices — master, admin, finance
			invoices := adminProtected.Group("/invoices")
			{
				invoices.POST("", middleware.RoleAuth(entity.RoleMaster, entity.RoleAdmin, entity.RoleFinance), h.Invoice.Create)
				invoices.GET("/partnership/:partnership_id", h.Invoice.GetByPartnership)
				invoices.GET("/:id", h.Invoice.GetByID)
				invoices.PUT("/:id/approve", middleware.RoleAuth(entity.RoleMaster, entity.RoleAdmin, entity.RoleFinance), h.Invoice.ManualApprove)
			}

			// Location Submissions
			locations := adminProtected.Group("/location-submissions")
			{
				locations.POST("", h.LocationSub.Create)
				locations.GET("", h.LocationSub.GetAll)
				locations.GET("/partnership/:partnershipId", h.LocationSub.GetByPartnership)
				locations.GET("/:id", h.LocationSub.GetByID)
				locations.PUT("/:id", h.LocationSub.Update)
				locations.PATCH("/:id/status", h.LocationSub.UpdateStatus)
				locations.POST("/:id/recalculate", h.LocationSub.RecalculateScore)
				locations.DELETE("/:id", middleware.RoleAuth(entity.RoleMaster, entity.RoleAdmin), h.LocationSub.Delete)
				locations.POST("/:id/surveys", h.LocationSub.CreateSurvey)
				locations.POST("/:id/files", h.LocationSub.AddFile)
				locations.DELETE("/:id/files/:fileId", h.LocationSub.DeleteFile)
				locations.POST("/:id/approve", middleware.RoleAuth(entity.RoleMaster, entity.RoleAdmin), h.LocationSub.Approve)
			}
		}
	}

	// ═══════════════════════════════════════════
	// LEGACY AUTH (kept for backward compat / mitra later)
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
