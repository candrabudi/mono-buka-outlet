package router

import (
	"github.com/franchise-system/backend/internal/handler"
	"github.com/franchise-system/backend/internal/middleware"
	"github.com/gin-gonic/gin"
)

type Handlers struct {
	Auth      *handler.AuthHandler
	AdminAuth *handler.AdminAuthHandler
	Mitra     *handler.MitraHandler

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

	// Register all route groups
	RegisterPublicRoutes(r, h, jwtSecret)
	RegisterPanelRoutes(r, h, jwtSecret)
	RegisterMitraRoutes(r, h, jwtSecret)
}
