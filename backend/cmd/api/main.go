package main

import (
	"flag"
	"log"
	"os"

	"github.com/franchise-system/backend/config"
	"github.com/franchise-system/backend/internal/handler"
	"github.com/franchise-system/backend/internal/migration"
	"github.com/franchise-system/backend/internal/repository/postgres"
	"github.com/franchise-system/backend/internal/router"
	"github.com/franchise-system/backend/internal/seeder"
	"github.com/franchise-system/backend/internal/service/chatbot"
	"github.com/franchise-system/backend/internal/service/email"
	"github.com/franchise-system/backend/internal/service/midtrans"
	"github.com/franchise-system/backend/internal/usecase"
	"github.com/gin-gonic/gin"
)

func main() {
	migrateUp := flag.Bool("migrate", false, "Run database migrations")
	migrateDown := flag.Bool("migrate-down", false, "Rollback database migrations")
	migrateFresh := flag.Bool("migrate-fresh", false, "Fresh migration (drop all + migrate)")
	seed := flag.Bool("seed", false, "Run database seeders")
	flag.Parse()

	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("❌ Failed to load config: %v", err)
	}

	db, err := config.NewDatabase(cfg.Database)
	if err != nil {
		log.Fatalf("❌ Failed to connect database: %v", err)
	}
	defer db.Close()

	migrator := migration.NewMigrator(db, "./migrations")

	if *migrateUp {
		if err := migrator.Up(); err != nil {
			log.Fatalf("❌ Migration failed: %v", err)
		}
		log.Println("✅ Migrations applied successfully")
		if !*seed {
			os.Exit(0)
		}
	}
	if *migrateDown {
		if err := migrator.Down(); err != nil {
			log.Fatalf("❌ Migration rollback failed: %v", err)
		}
		log.Println("✅ Migrations rolled back successfully")
		os.Exit(0)
	}
	if *migrateFresh {
		if err := migrator.Fresh(); err != nil {
			log.Fatalf("❌ Fresh migration failed: %v", err)
		}
		log.Println("✅ Fresh migration completed")
		if !*seed {
			os.Exit(0)
		}
	}
	if *seed {
		s := seeder.NewSeeder(db)
		if err := s.Run(); err != nil {
			log.Fatalf("❌ Seeding failed: %v", err)
		}
		os.Exit(0)
	}

	userRepo := postgres.NewUserRepo(db)
	outletRepo := postgres.NewOutletRepo(db)
	outletCategoryRepo := postgres.NewOutletCategoryRepo(db)
	outletPackageRepo := postgres.NewOutletPackageRepo(db)
	partnershipRepo := postgres.NewPartnershipRepo(db)
	paymentRepo := postgres.NewPaymentRepo(db)
	agreementRepo := postgres.NewAgreementRepo(db)
	revenueRepo := postgres.NewRevenueRepo(db)
	activityLogRepo := postgres.NewActivityLogRepo(db)
	dashboardRepo := postgres.NewDashboardRepo(db)
	otpRepo := postgres.NewOTPRepo(db)
	meetingRepo := postgres.NewMeetingRepo(db)
	settingRepo := postgres.NewSystemSettingRepo(db)
	invoiceRepo := postgres.NewInvoiceRepo(db)
	locationSubRepo := postgres.NewLocationSubmissionRepo(db)
	partnershipAppRepo := postgres.NewPartnershipApplicationRepo(db)
	ebookRepo := postgres.NewEbookRepo(db)
	ebookOrderRepo := postgres.NewEbookOrderRepo(db)
	ebookCategoryRepo := postgres.NewEbookCategoryRepo(db)
	affCommissionRepo := postgres.NewAffiliatorCommissionRepo(db)
	affWithdrawalRepo := postgres.NewAffiliatorWithdrawalRepo(db)

	emailService := email.NewEmailService(email.SMTPConfig{
		Host:     cfg.SMTP.Host,
		Port:     cfg.SMTP.Port,
		Username: cfg.SMTP.Username,
		Password: cfg.SMTP.Password,
		From:     cfg.SMTP.From,
		FromName: cfg.SMTP.FromName,
	})

	authUC := usecase.NewAuthUseCase(userRepo, cfg.JWT.Secret, cfg.JWT.ExpiryHours)
	adminAuthUC := usecase.NewAdminAuthUseCase(userRepo, otpRepo, emailService, cfg.JWT.Secret, cfg.JWT.ExpiryHours)
	outletUC := usecase.NewOutletUseCase(outletRepo, outletCategoryRepo)
	outletCategoryUC := usecase.NewOutletCategoryUseCase(outletCategoryRepo)
	outletPackageUC := usecase.NewOutletPackageUseCase(outletPackageRepo)
	partnershipUC := usecase.NewPartnershipUseCase(partnershipRepo, activityLogRepo)
	paymentUC := usecase.NewPaymentUseCase(paymentRepo, partnershipRepo, activityLogRepo)
	agreementUC := usecase.NewAgreementUseCase(agreementRepo, partnershipRepo, activityLogRepo)
	revenueUC := usecase.NewRevenueUseCase(revenueRepo, partnershipRepo)
	dashboardUC := usecase.NewDashboardUseCase(dashboardRepo, db)
	meetingUC := usecase.NewMeetingUseCase(meetingRepo)

	chatService := chatbot.NewService(db, cfg.OpenAI.APIKey)
	aiKnowledgeRepo := postgres.NewAIKnowledgeRepo(db)

	affiliatorUC := usecase.NewAffiliatorUseCase(userRepo, partnershipRepo, affCommissionRepo, affWithdrawalRepo)

	invoiceHandler := handler.NewInvoiceHandler(invoiceRepo, partnershipRepo, settingRepo, midtrans.NewService(settingRepo))

	handlers := router.Handlers{
		Auth:      handler.NewAuthHandler(authUC),
		AdminAuth: handler.NewAdminAuthHandler(adminAuthUC),
		Mitra:     handler.NewMitraHandler(partnershipAppRepo, outletRepo, outletPackageRepo, partnershipRepo, userRepo),

		Outlet:         handler.NewOutletHandler(outletUC),
		OutletCategory: handler.NewOutletCategoryHandler(outletCategoryUC),

		Partnership:    handler.NewPartnershipHandler(partnershipUC),
		Payment:        handler.NewPaymentHandler(paymentUC),
		Agreement:      handler.NewAgreementHandler(agreementUC),
		Revenue:        handler.NewRevenueHandler(revenueUC),
		Dashboard:      handler.NewDashboardHandler(dashboardUC),
		Upload:         handler.NewUploadHandler(cfg.Upload.Dir, cfg.Upload.MaxSize, cfg.App.URL),
		OutletPackage:  handler.NewOutletPackageHandler(outletPackageUC),
		Meeting:        handler.NewMeetingHandler(meetingUC, cfg.Upload.Dir, cfg.App.URL),
		Setting:        handler.NewSettingHandler(settingRepo),
		Invoice:        invoiceHandler,
		LocationSub:    handler.NewLocationSubmissionHandler(locationSubRepo),
		Ebook:          handler.NewEbookHandler(ebookRepo, ebookOrderRepo, settingRepo, midtrans.NewService(settingRepo), cfg.Upload.Dir),
		EbookCategory:  handler.NewEbookCategoryHandler(ebookCategoryRepo),
		Chat:           handler.NewChatHandler(chatService),
		AIAdmin:        handler.NewAIAdminHandler(aiKnowledgeRepo, chatService),
		Affiliator:     handler.NewAffiliatorHandler(affiliatorUC),
	}

	if cfg.App.Env == "production" {
		gin.SetMode(gin.ReleaseMode)
	}
	r := gin.Default()

	router.Setup(r, handlers, cfg.JWT.Secret, cfg.CORS.AllowedOrigins)

	os.MkdirAll(cfg.Upload.Dir, os.ModePerm)
	r.Static("/uploads", cfg.Upload.Dir)

	// Start background invoice expiry scheduler
	invoiceHandler.StartExpiryScheduler()

	log.Printf("🚀 Server starting on port %s", cfg.App.Port)
	if err := r.Run(":" + cfg.App.Port); err != nil {
		log.Fatalf("❌ Server failed: %v", err)
	}
}
