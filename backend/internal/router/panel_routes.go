package router

import (
	"github.com/franchise-system/backend/internal/entity"
	"github.com/franchise-system/backend/internal/middleware"
	"github.com/gin-gonic/gin"
)

// RegisterPanelRoutes sets up all /api/v1/admin/* routes for the admin panel.
func RegisterPanelRoutes(r *gin.Engine, h Handlers, jwtSecret string) {
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

			// Partnerships — master, admin (read), finance (can view)
			partnerships := adminProtected.Group("/partnerships")
			{
				partnerships.POST("", middleware.RoleAuth(entity.RoleMaster, entity.RoleAdmin), h.Partnership.Create)
				partnerships.GET("", middleware.RoleAuth(
					entity.RoleMaster, entity.RoleAdmin, entity.RoleFinance,
				), h.Partnership.GetAll)
				partnerships.GET("/:id", h.Partnership.GetByID)
				partnerships.PATCH("/:id/status", middleware.RoleAuth(entity.RoleMaster, entity.RoleAdmin), h.Partnership.UpdateStatus)
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
				invoices.GET("/:id/check-status", h.Invoice.CheckStatus)
				invoices.PUT("/:id/approve", middleware.RoleAuth(entity.RoleMaster, entity.RoleAdmin, entity.RoleFinance), h.Invoice.ManualApprove)
				invoices.POST("/sync-pending", middleware.RoleAuth(entity.RoleMaster, entity.RoleAdmin), h.Invoice.SyncAllPending)
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

			// Partnership Applications — master, admin
			apps := adminProtected.Group("/partnership-applications")
			apps.Use(middleware.RoleAuth(entity.RoleMaster, entity.RoleAdmin))
			{
				apps.GET("", h.Mitra.ListAllApplications)
				apps.GET("/:id", h.Mitra.GetApplication)
				apps.PATCH("/:id/review", h.Mitra.ReviewApplication)
			}

			// Ebooks — master, admin
			ebooks := adminProtected.Group("/ebooks")
			ebooks.Use(middleware.RoleAuth(entity.RoleMaster, entity.RoleAdmin))
			{
				ebooks.POST("", h.Ebook.Create)
				ebooks.GET("", h.Ebook.GetAll)
				ebooks.GET("/:id", h.Ebook.GetByID)
				ebooks.PUT("/:id", h.Ebook.Update)
				ebooks.DELETE("/:id", h.Ebook.Delete)
				ebooks.PATCH("/:id/toggle", h.Ebook.ToggleActive)
			}

			// Ebook Categories — master, admin
			ebookCats := adminProtected.Group("/ebook-categories")
			ebookCats.Use(middleware.RoleAuth(entity.RoleMaster, entity.RoleAdmin))
			{
				ebookCats.POST("", h.EbookCategory.Create)
				ebookCats.GET("", h.EbookCategory.GetAll)
				ebookCats.GET("/:id", h.EbookCategory.GetByID)
				ebookCats.PUT("/:id", h.EbookCategory.Update)
				ebookCats.DELETE("/:id", h.EbookCategory.Delete)
				ebookCats.PATCH("/:id/toggle", h.EbookCategory.ToggleActive)
			}

			// Ebook Orders — master, admin (download approval)
			ebookOrders := adminProtected.Group("/ebook-orders")
			ebookOrders.Use(middleware.RoleAuth(entity.RoleMaster, entity.RoleAdmin))
			{
				ebookOrders.GET("", h.Ebook.ListAllOrders)
				ebookOrders.GET("/download-requests", h.Ebook.ListDownloadRequests)
				ebookOrders.PATCH("/:id/approve-download", h.Ebook.ApproveDownload)
				ebookOrders.PATCH("/:id/reject-download", h.Ebook.RejectDownload)
				ebookOrders.PATCH("/:id/approve-payment", h.Ebook.ApprovePayment)
				ebookOrders.PATCH("/:id/reject-payment", h.Ebook.RejectPayment)
			}

			// AI Konsultan Settings — master, admin
			ai := adminProtected.Group("/ai")
			ai.Use(middleware.RoleAuth(entity.RoleMaster, entity.RoleAdmin))
			{
				// Knowledge Base
				ai.GET("/knowledge", h.AIAdmin.ListKnowledge)
				ai.POST("/knowledge", h.AIAdmin.CreateKnowledge)
				ai.PUT("/knowledge/:id", h.AIAdmin.UpdateKnowledge)
				ai.DELETE("/knowledge/:id", h.AIAdmin.DeleteKnowledge)

				// Categories
				ai.GET("/categories", h.AIAdmin.ListCategories)
				ai.POST("/categories", h.AIAdmin.CreateCategory)

				// System Prompts
				ai.GET("/prompts", h.AIAdmin.ListPrompts)
				ai.POST("/prompts", h.AIAdmin.CreatePrompt)
				ai.PUT("/prompts/:id", h.AIAdmin.UpdatePrompt)

				// Config
				ai.GET("/config", h.AIAdmin.GetConfig)
				ai.PUT("/config", h.AIAdmin.UpdateConfig)

				// Cache
				ai.POST("/cache/invalidate", h.AIAdmin.InvalidateCache)
			}

			// Affiliator — affiliator only (self-service)
			affiliator := adminProtected.Group("/affiliator")
			affiliator.Use(middleware.RoleAuth(entity.RoleAffiliator))
			{
				affiliator.GET("/dashboard", h.Affiliator.GetDashboard)
				affiliator.GET("/partnerships", h.Affiliator.GetMyPartnerships)
				affiliator.GET("/referral-code", h.Affiliator.GetReferralCode)
			}
		}
	}
}
