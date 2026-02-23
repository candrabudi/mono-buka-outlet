package repository

import (
	"context"

	"github.com/franchise-system/backend/internal/entity"
	"github.com/google/uuid"
)

type UserRepository interface {
	Create(ctx context.Context, user *entity.User) error
	FindByID(ctx context.Context, id uuid.UUID) (*entity.User, error)
	FindByEmail(ctx context.Context, email string) (*entity.User, error)
	FindAll(ctx context.Context, role string, page, limit int) ([]*entity.User, int, error)
	Update(ctx context.Context, user *entity.User) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type OutletCategoryRepository interface {
	Create(ctx context.Context, cat *entity.OutletCategory) error
	FindByID(ctx context.Context, id uuid.UUID) (*entity.OutletCategory, error)
	FindBySlug(ctx context.Context, slug string) (*entity.OutletCategory, error)
	FindAll(ctx context.Context, activeOnly bool) ([]*entity.OutletCategory, error)
	Update(ctx context.Context, cat *entity.OutletCategory) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type OutletRepository interface {
	Create(ctx context.Context, outlet *entity.Outlet) error
	FindByID(ctx context.Context, id uuid.UUID) (*entity.Outlet, error)
	FindBySlug(ctx context.Context, slug string) (*entity.Outlet, error)
	FindAll(ctx context.Context, filter OutletFilter) ([]*entity.Outlet, int, error)
	Update(ctx context.Context, outlet *entity.Outlet) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type OutletFilter struct {
	Search     string
	CategoryID *uuid.UUID
	City       string
	Province   string
	Active     *bool
	Featured   *bool
	Page       int
	Limit      int
}

type LeadRepository interface {
	Create(ctx context.Context, lead *entity.Lead) error
	FindByID(ctx context.Context, id uuid.UUID) (*entity.Lead, error)
	FindAll(ctx context.Context, brandID *uuid.UUID, status string, page, limit int) ([]*entity.Lead, int, error)
	FindByBrandGrouped(ctx context.Context, brandID *uuid.UUID) (map[string][]*entity.Lead, error)
	Update(ctx context.Context, lead *entity.Lead) error
	UpdateStatus(ctx context.Context, id uuid.UUID, status string, progress int) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type LocationRepository interface {
	Create(ctx context.Context, location *entity.Location) error
	FindByID(ctx context.Context, id uuid.UUID) (*entity.Location, error)
	FindByLeadID(ctx context.Context, leadID uuid.UUID) ([]*entity.Location, error)
	Update(ctx context.Context, location *entity.Location) error
	UpdateApproval(ctx context.Context, id uuid.UUID, status, notes string) error
}

type PartnershipRepository interface {
	Create(ctx context.Context, partnership *entity.Partnership) error
	FindByID(ctx context.Context, id uuid.UUID) (*entity.Partnership, error)
	FindAll(ctx context.Context, brandID *uuid.UUID, mitraID *uuid.UUID, page, limit int) ([]*entity.Partnership, int, error)
	FindByMitraID(ctx context.Context, mitraID uuid.UUID) ([]*entity.Partnership, error)
	Update(ctx context.Context, partnership *entity.Partnership) error
	UpdateProgress(ctx context.Context, id uuid.UUID, progress int, status string) error
}

type PaymentRepository interface {
	Create(ctx context.Context, payment *entity.Payment) error
	FindByID(ctx context.Context, id uuid.UUID) (*entity.Payment, error)
	FindByPartnershipID(ctx context.Context, partnershipID uuid.UUID) ([]*entity.Payment, error)
	Update(ctx context.Context, payment *entity.Payment) error
	Verify(ctx context.Context, id uuid.UUID, status string, verifiedBy uuid.UUID) error
}

type AgreementRepository interface {
	Create(ctx context.Context, agreement *entity.Agreement) error
	FindByID(ctx context.Context, id uuid.UUID) (*entity.Agreement, error)
	FindByPartnershipID(ctx context.Context, partnershipID uuid.UUID) ([]*entity.Agreement, error)
	Update(ctx context.Context, agreement *entity.Agreement) error
	Sign(ctx context.Context, id uuid.UUID) error
}

type RevenueRepository interface {
	Create(ctx context.Context, revenue *entity.Revenue) error
	FindByID(ctx context.Context, id uuid.UUID) (*entity.Revenue, error)
	FindByPartnershipID(ctx context.Context, partnershipID uuid.UUID) ([]*entity.Revenue, error)
	FindByBrandID(ctx context.Context, brandID uuid.UUID, month string) ([]*entity.Revenue, error)
	Update(ctx context.Context, revenue *entity.Revenue) error
}

type ActivityLogRepository interface {
	Create(ctx context.Context, log *entity.ActivityLog) error
	FindByEntity(ctx context.Context, entityType string, entityID uuid.UUID) ([]*entity.ActivityLog, error)
}

type NotificationRepository interface {
	Create(ctx context.Context, notification *entity.Notification) error
	FindByUserID(ctx context.Context, userID uuid.UUID, unreadOnly bool) ([]*entity.Notification, error)
	MarkAsRead(ctx context.Context, id uuid.UUID) error
	MarkAllAsRead(ctx context.Context, userID uuid.UUID) error
}

type DashboardRepository interface {
	GetTotalLeads(ctx context.Context, brandID *uuid.UUID) (int, error)
	GetActiveMitra(ctx context.Context, brandID *uuid.UUID) (int, error)
	GetTotalInvestment(ctx context.Context, brandID *uuid.UUID) (float64, error)
	GetMonthlyRevenue(ctx context.Context, brandID *uuid.UUID, month string) (float64, error)
	GetLeadsByStatus(ctx context.Context, brandID *uuid.UUID) (map[string]int, error)
	GetRevenueChart(ctx context.Context, brandID *uuid.UUID, months int) ([]map[string]interface{}, error)
}

type OutletPackageRepository interface {
	Create(ctx context.Context, pkg *entity.OutletPackage) error
	FindByID(ctx context.Context, id uuid.UUID) (*entity.OutletPackage, error)
	FindByOutletID(ctx context.Context, outletID uuid.UUID) ([]*entity.OutletPackage, error)
	Update(ctx context.Context, pkg *entity.OutletPackage) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type MeetingRepository interface {
	Create(ctx context.Context, m *entity.Meeting) error
	FindByID(ctx context.Context, id uuid.UUID) (*entity.Meeting, error)
	FindAll(ctx context.Context, status string, meetingType string, page, limit int) ([]*entity.Meeting, int, error)
	Update(ctx context.Context, m *entity.Meeting) error
	Delete(ctx context.Context, id uuid.UUID) error

	// Participants
	AddParticipant(ctx context.Context, p *entity.MeetingParticipant) error
	FindParticipants(ctx context.Context, meetingID uuid.UUID) ([]*entity.MeetingParticipant, error)
	DeleteParticipant(ctx context.Context, id uuid.UUID) error

	// Notes
	UpsertNotes(ctx context.Context, n *entity.MeetingNotes) error
	FindNotes(ctx context.Context, meetingID uuid.UUID) (*entity.MeetingNotes, error)

	// Action Plans
	AddActionPlan(ctx context.Context, a *entity.MeetingActionPlan) error
	FindActionPlans(ctx context.Context, meetingID uuid.UUID) ([]*entity.MeetingActionPlan, error)
	UpdateActionPlan(ctx context.Context, a *entity.MeetingActionPlan) error
	DeleteActionPlan(ctx context.Context, id uuid.UUID) error

	// Files
	AddFile(ctx context.Context, f *entity.MeetingFile) error
	FindFiles(ctx context.Context, meetingID uuid.UUID) ([]*entity.MeetingFile, error)
	DeleteFile(ctx context.Context, id uuid.UUID) error
}

type SystemSettingRepository interface {
	FindAll(ctx context.Context) ([]*entity.SystemSetting, error)
	FindByGroup(ctx context.Context, group string) ([]*entity.SystemSetting, error)
	FindByKey(ctx context.Context, key string) (*entity.SystemSetting, error)
	Upsert(ctx context.Context, key, value string) error
	BulkUpsert(ctx context.Context, settings map[string]string) error
}

type InvoiceRepository interface {
	Create(ctx context.Context, inv *entity.Invoice) error
	FindByID(ctx context.Context, id uuid.UUID) (*entity.Invoice, error)
	FindByPartnershipID(ctx context.Context, partnershipID uuid.UUID) ([]*entity.Invoice, error)
	FindByOrderID(ctx context.Context, orderID string) (*entity.Invoice, error)
	UpdateMidtransStatus(ctx context.Context, orderID string, txnStatus, paymentType, txnID string) error
	ManualApprove(ctx context.Context, id uuid.UUID, proofURL string) error
	GenerateInvoiceNumber(ctx context.Context) (string, error)
}

type PartnershipApplicationRepository interface {
	Create(ctx context.Context, app *entity.PartnershipApplication) error
	FindByID(ctx context.Context, id uuid.UUID) (*entity.PartnershipApplication, error)
	FindByMitraID(ctx context.Context, mitraID uuid.UUID) ([]*entity.PartnershipApplication, error)
	FindAll(ctx context.Context, status string, page, limit int) ([]*entity.PartnershipApplication, int, error)
	UpdateStatus(ctx context.Context, id uuid.UUID, status, adminNotes string, reviewedBy uuid.UUID) error
	HasActiveApplication(ctx context.Context, mitraID, outletID, packageID uuid.UUID) (bool, error)
	CancelByMitra(ctx context.Context, id, mitraID uuid.UUID) error
}

type EbookRepository interface {
	Create(ctx context.Context, ebook *entity.Ebook) error
	FindByID(ctx context.Context, id uuid.UUID) (*entity.Ebook, error)
	FindBySlug(ctx context.Context, slug string) (*entity.Ebook, error)
	FindAll(ctx context.Context, activeOnly bool, search string, page, limit int) ([]*entity.Ebook, int, error)
	Update(ctx context.Context, ebook *entity.Ebook) error
	Delete(ctx context.Context, id uuid.UUID) error
	IncrementSold(ctx context.Context, id uuid.UUID) error
	SyncCategories(ctx context.Context, ebookID uuid.UUID, categoryIDs []uuid.UUID) error
}

type EbookCategoryRepository interface {
	Create(ctx context.Context, cat *entity.EbookCategory) error
	FindByID(ctx context.Context, id uuid.UUID) (*entity.EbookCategory, error)
	FindBySlug(ctx context.Context, slug string) (*entity.EbookCategory, error)
	FindAll(ctx context.Context, activeOnly bool) ([]*entity.EbookCategory, error)
	Update(ctx context.Context, cat *entity.EbookCategory) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type EbookOrderRepository interface {
	Create(ctx context.Context, order *entity.EbookOrder) error
	FindByID(ctx context.Context, id uuid.UUID) (*entity.EbookOrder, error)
	FindByUserID(ctx context.Context, userID uuid.UUID) ([]*entity.EbookOrder, error)
	FindByMidtransOrderID(ctx context.Context, orderID string) (*entity.EbookOrder, error)
	HasUserPurchased(ctx context.Context, userID, ebookID uuid.UUID) (bool, error)
	UpdateMidtransStatus(ctx context.Context, orderID string, txnStatus, paymentType, txnID string) error
	GenerateOrderNumber(ctx context.Context) (string, error)

	// Download approval workflow
	RequestDownload(ctx context.Context, id uuid.UUID) error
	ApproveDownload(ctx context.Context, id uuid.UUID, note string) error
	RejectDownload(ctx context.Context, id uuid.UUID, note string) error
	FindPendingDownloads(ctx context.Context) ([]*entity.EbookOrder, error)
	FindAllOrders(ctx context.Context, status, downloadStatus string, page, limit int) ([]*entity.EbookOrder, int, error)
	CancelOrder(ctx context.Context, id uuid.UUID) error
	UploadPaymentProof(ctx context.Context, id uuid.UUID, proofURL string) error
	ApprovePayment(ctx context.Context, id uuid.UUID) error
	RejectPayment(ctx context.Context, id uuid.UUID, note string) error
}
