package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/franchise-system/backend/internal/entity"
	"github.com/franchise-system/backend/internal/repository"
	"github.com/google/uuid"
)

type MeetingUseCase struct {
	repo repository.MeetingRepository
}

func NewMeetingUseCase(repo repository.MeetingRepository) *MeetingUseCase {
	return &MeetingUseCase{repo: repo}
}

// ─── Request Types ───

type CreateMeetingRequest struct {
	Title       string `json:"title" binding:"required"`
	MeetingType string `json:"meeting_type" binding:"required"`
	MeetingDate string `json:"meeting_date" binding:"required"`
	Duration    int    `json:"duration"`
	MeetingLink string `json:"meeting_link"`
	Location    string `json:"location"`
}

type UpdateMeetingRequest struct {
	Title       string `json:"title"`
	MeetingType string `json:"meeting_type"`
	MeetingDate string `json:"meeting_date"`
	Duration    int    `json:"duration"`
	MeetingLink string `json:"meeting_link"`
	Location    string `json:"location"`
	Status      string `json:"status"`
}

type AddParticipantRequest struct {
	UserID        string `json:"user_id"`
	ExternalName  string `json:"external_name"`
	ExternalEmail string `json:"external_email"`
	ExternalPhone string `json:"external_phone"`
	Role          string `json:"role" binding:"required"`
}

type SaveNotesRequest struct {
	Purpose          string `json:"purpose"`
	Summary          string `json:"summary"`
	DiscussionPoints string `json:"discussion_points"`
	Decision         string `json:"decision"`
	NextStep         string `json:"next_step"`
}

type AddActionPlanRequest struct {
	TaskName string `json:"task_name" binding:"required"`
	PIC      string `json:"pic" binding:"required"`
	Deadline string `json:"deadline"`
	Status   string `json:"status"`
}

type UpdateActionPlanRequest struct {
	TaskName string `json:"task_name"`
	PIC      string `json:"pic"`
	Deadline string `json:"deadline"`
	Status   string `json:"status"`
}

// ─── Meeting CRUD ───

func (uc *MeetingUseCase) Create(ctx context.Context, req CreateMeetingRequest, createdBy uuid.UUID) (*entity.Meeting, error) {
	meetingDate, err := time.Parse(time.RFC3339, req.MeetingDate)
	if err != nil {
		meetingDate, err = time.Parse("2006-01-02T15:04", req.MeetingDate)
		if err != nil {
			return nil, fmt.Errorf("format tanggal tidak valid")
		}
	}
	duration := req.Duration
	if duration <= 0 {
		duration = 60
	}

	m := &entity.Meeting{
		ID:          uuid.New(),
		Title:       req.Title,
		MeetingType: req.MeetingType,
		MeetingDate: meetingDate,
		Duration:    duration,
		MeetingLink: req.MeetingLink,
		Location:    req.Location,
		Status:      entity.MeetingStatusScheduled,
		CreatedBy:   createdBy,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	if err := uc.repo.Create(ctx, m); err != nil {
		return nil, fmt.Errorf("gagal membuat meeting: %w", err)
	}
	return m, nil
}

func (uc *MeetingUseCase) GetByID(ctx context.Context, id uuid.UUID) (*entity.Meeting, error) {
	return uc.repo.FindByID(ctx, id)
}

func (uc *MeetingUseCase) GetAll(ctx context.Context, status, meetingType string, page, limit int) ([]*entity.Meeting, int, error) {
	return uc.repo.FindAll(ctx, status, meetingType, page, limit)
}

func (uc *MeetingUseCase) Update(ctx context.Context, id uuid.UUID, req UpdateMeetingRequest) (*entity.Meeting, error) {
	m, err := uc.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if req.Title != "" {
		m.Title = req.Title
	}
	if req.MeetingType != "" {
		m.MeetingType = req.MeetingType
	}
	if req.MeetingDate != "" {
		t, err := time.Parse(time.RFC3339, req.MeetingDate)
		if err != nil {
			t, err = time.Parse("2006-01-02T15:04", req.MeetingDate)
			if err != nil {
				return nil, fmt.Errorf("format tanggal tidak valid")
			}
		}
		m.MeetingDate = t
	}
	if req.Duration > 0 {
		m.Duration = req.Duration
	}
	if req.MeetingLink != "" {
		m.MeetingLink = req.MeetingLink
	}
	if req.Location != "" {
		m.Location = req.Location
	}
	if req.Status != "" {
		m.Status = req.Status
	}
	if err := uc.repo.Update(ctx, m); err != nil {
		return nil, fmt.Errorf("gagal update meeting: %w", err)
	}
	return uc.repo.FindByID(ctx, id)
}

func (uc *MeetingUseCase) Delete(ctx context.Context, id uuid.UUID) error {
	return uc.repo.Delete(ctx, id)
}

// ─── Participants ───

func (uc *MeetingUseCase) AddParticipant(ctx context.Context, meetingID uuid.UUID, req AddParticipantRequest) (*entity.MeetingParticipant, error) {
	p := &entity.MeetingParticipant{
		ID:            uuid.New(),
		MeetingID:     meetingID,
		ExternalName:  req.ExternalName,
		ExternalEmail: req.ExternalEmail,
		ExternalPhone: req.ExternalPhone,
		Role:          req.Role,
		CreatedAt:     time.Now(),
	}
	if req.UserID != "" {
		uid, err := uuid.Parse(req.UserID)
		if err == nil {
			p.UserID = &uid
		}
	}
	if err := uc.repo.AddParticipant(ctx, p); err != nil {
		return nil, fmt.Errorf("gagal menambah peserta: %w", err)
	}
	return p, nil
}

func (uc *MeetingUseCase) DeleteParticipant(ctx context.Context, id uuid.UUID) error {
	return uc.repo.DeleteParticipant(ctx, id)
}

// ─── Notes ───

func (uc *MeetingUseCase) SaveNotes(ctx context.Context, meetingID uuid.UUID, req SaveNotesRequest) (*entity.MeetingNotes, error) {
	n := &entity.MeetingNotes{
		ID:               uuid.New(),
		MeetingID:        meetingID,
		Purpose:          req.Purpose,
		Summary:          req.Summary,
		DiscussionPoints: req.DiscussionPoints,
		Decision:         req.Decision,
		NextStep:         req.NextStep,
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
	}
	if err := uc.repo.UpsertNotes(ctx, n); err != nil {
		return nil, fmt.Errorf("gagal menyimpan notulensi: %w", err)
	}
	return uc.repo.FindNotes(ctx, meetingID)
}

// ─── Action Plans ───

func (uc *MeetingUseCase) AddActionPlan(ctx context.Context, meetingID uuid.UUID, req AddActionPlanRequest) (*entity.MeetingActionPlan, error) {
	a := &entity.MeetingActionPlan{
		ID:        uuid.New(),
		MeetingID: meetingID,
		TaskName:  req.TaskName,
		PIC:       req.PIC,
		Status:    "pending",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	if req.Deadline != "" {
		t, err := time.Parse("2006-01-02", req.Deadline)
		if err == nil {
			a.Deadline = &t
		}
	}
	if req.Status != "" {
		a.Status = req.Status
	}
	if err := uc.repo.AddActionPlan(ctx, a); err != nil {
		return nil, fmt.Errorf("gagal menambah action plan: %w", err)
	}
	return a, nil
}

func (uc *MeetingUseCase) UpdateActionPlan(ctx context.Context, id uuid.UUID, req UpdateActionPlanRequest) error {
	a := &entity.MeetingActionPlan{ID: id, UpdatedAt: time.Now()}
	if req.TaskName != "" {
		a.TaskName = req.TaskName
	}
	if req.PIC != "" {
		a.PIC = req.PIC
	}
	if req.Status != "" {
		a.Status = req.Status
	}
	if req.Deadline != "" {
		t, err := time.Parse("2006-01-02", req.Deadline)
		if err == nil {
			a.Deadline = &t
		}
	}
	return uc.repo.UpdateActionPlan(ctx, a)
}

func (uc *MeetingUseCase) DeleteActionPlan(ctx context.Context, id uuid.UUID) error {
	return uc.repo.DeleteActionPlan(ctx, id)
}

// ─── Files ───

func (uc *MeetingUseCase) AddFile(ctx context.Context, meetingID uuid.UUID, f *entity.MeetingFile) error {
	f.ID = uuid.New()
	f.MeetingID = meetingID
	f.UploadedAt = time.Now()
	return uc.repo.AddFile(ctx, f)
}

func (uc *MeetingUseCase) DeleteFile(ctx context.Context, id uuid.UUID) error {
	return uc.repo.DeleteFile(ctx, id)
}
