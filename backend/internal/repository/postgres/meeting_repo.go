package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/franchise-system/backend/internal/entity"
	"github.com/google/uuid"
)

type MeetingRepo struct {
	db *sql.DB
}

func NewMeetingRepo(db *sql.DB) *MeetingRepo {
	return &MeetingRepo{db: db}
}

// ═══ MEETING CRUD ═══

func (r *MeetingRepo) Create(ctx context.Context, m *entity.Meeting) error {
	query := `INSERT INTO meetings (id, title, meeting_type, meeting_date, duration, meeting_link, location, status, created_by, created_at, updated_at)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11)`
	_, err := r.db.ExecContext(ctx, query,
		m.ID, m.Title, m.MeetingType, m.MeetingDate, m.Duration,
		m.MeetingLink, m.Location, m.Status, m.CreatedBy, m.CreatedAt, m.UpdatedAt,
	)
	return err
}

func (r *MeetingRepo) FindByID(ctx context.Context, id uuid.UUID) (*entity.Meeting, error) {
	m := &entity.Meeting{}
	var meetingLink, location sql.NullString
	query := `SELECT id, title, meeting_type, meeting_date, duration, meeting_link, location, status, created_by, created_at, updated_at
		FROM meetings WHERE id = $1 AND deleted_at IS NULL`
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&m.ID, &m.Title, &m.MeetingType, &m.MeetingDate, &m.Duration,
		&meetingLink, &location, &m.Status, &m.CreatedBy, &m.CreatedAt, &m.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("meeting tidak ditemukan")
	}
	if err != nil {
		return nil, err
	}
	m.MeetingLink = meetingLink.String
	m.Location = location.String

	// Load relations
	m.Participants, _ = r.FindParticipants(ctx, id)
	m.Notes, _ = r.FindNotes(ctx, id)
	m.ActionPlans, _ = r.FindActionPlans(ctx, id)
	m.Files, _ = r.FindFiles(ctx, id)
	return m, nil
}

func (r *MeetingRepo) FindAll(ctx context.Context, status string, meetingType string, page, limit int) ([]*entity.Meeting, int, error) {
	var total int
	countQ := "SELECT COUNT(*) FROM meetings WHERE deleted_at IS NULL"
	args := []interface{}{}
	idx := 1

	if status != "" {
		countQ += fmt.Sprintf(" AND status = $%d", idx)
		args = append(args, status)
		idx++
	}
	if meetingType != "" {
		countQ += fmt.Sprintf(" AND meeting_type = $%d", idx)
		args = append(args, meetingType)
		idx++
	}

	if err := r.db.QueryRowContext(ctx, countQ, args...).Scan(&total); err != nil {
		return nil, 0, err
	}

	q := "SELECT id, title, meeting_type, meeting_date, duration, meeting_link, location, status, created_by, created_at, updated_at FROM meetings WHERE deleted_at IS NULL"
	qArgs := []interface{}{}
	qIdx := 1

	if status != "" {
		q += fmt.Sprintf(" AND status = $%d", qIdx)
		qArgs = append(qArgs, status)
		qIdx++
	}
	if meetingType != "" {
		q += fmt.Sprintf(" AND meeting_type = $%d", qIdx)
		qArgs = append(qArgs, meetingType)
		qIdx++
	}

	q += fmt.Sprintf(" ORDER BY meeting_date DESC LIMIT $%d OFFSET $%d", qIdx, qIdx+1)
	qArgs = append(qArgs, limit, (page-1)*limit)

	rows, err := r.db.QueryContext(ctx, q, qArgs...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var list []*entity.Meeting
	for rows.Next() {
		m := &entity.Meeting{}
		var meetingLink, location sql.NullString
		if err := rows.Scan(&m.ID, &m.Title, &m.MeetingType, &m.MeetingDate, &m.Duration,
			&meetingLink, &location, &m.Status, &m.CreatedBy, &m.CreatedAt, &m.UpdatedAt); err != nil {
			return nil, 0, err
		}
		m.MeetingLink = meetingLink.String
		m.Location = location.String
		// Load participant count only for list
		participants, _ := r.FindParticipants(ctx, m.ID)
		m.Participants = participants
		list = append(list, m)
	}
	return list, total, nil
}

func (r *MeetingRepo) Update(ctx context.Context, m *entity.Meeting) error {
	query := `UPDATE meetings SET title=$1, meeting_type=$2, meeting_date=$3, duration=$4,
		meeting_link=$5, location=$6, status=$7, updated_at=$8 WHERE id=$9 AND deleted_at IS NULL`
	_, err := r.db.ExecContext(ctx, query,
		m.Title, m.MeetingType, m.MeetingDate, m.Duration,
		m.MeetingLink, m.Location, m.Status, time.Now(), m.ID,
	)
	return err
}

func (r *MeetingRepo) Delete(ctx context.Context, id uuid.UUID) error {
	_, err := r.db.ExecContext(ctx, `UPDATE meetings SET deleted_at = $1 WHERE id = $2`, time.Now(), id)
	return err
}

// ═══ PARTICIPANTS ═══

func (r *MeetingRepo) AddParticipant(ctx context.Context, p *entity.MeetingParticipant) error {
	query := `INSERT INTO meeting_participants (id, meeting_id, user_id, external_name, external_email, external_phone, role, created_at)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8)`
	_, err := r.db.ExecContext(ctx, query,
		p.ID, p.MeetingID, p.UserID, p.ExternalName, p.ExternalEmail, p.ExternalPhone, p.Role, p.CreatedAt,
	)
	return err
}

func (r *MeetingRepo) FindParticipants(ctx context.Context, meetingID uuid.UUID) ([]*entity.MeetingParticipant, error) {
	query := `SELECT mp.id, mp.meeting_id, mp.user_id, mp.external_name, mp.external_email, mp.external_phone, mp.role, mp.created_at,
		u.name
		FROM meeting_participants mp
		LEFT JOIN users u ON mp.user_id = u.id
		WHERE mp.meeting_id = $1 ORDER BY mp.created_at ASC`
	rows, err := r.db.QueryContext(ctx, query, meetingID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []*entity.MeetingParticipant
	for rows.Next() {
		p := &entity.MeetingParticipant{}
		var userID sql.NullString
		var extName, extEmail, extPhone sql.NullString
		var userName sql.NullString
		if err := rows.Scan(&p.ID, &p.MeetingID, &userID, &extName, &extEmail, &extPhone, &p.Role, &p.CreatedAt, &userName); err != nil {
			return nil, err
		}
		if userID.Valid {
			uid, _ := uuid.Parse(userID.String)
			p.UserID = &uid
			p.User = &entity.User{Name: userName.String}
		}
		p.ExternalName = extName.String
		p.ExternalEmail = extEmail.String
		p.ExternalPhone = extPhone.String
		list = append(list, p)
	}
	return list, nil
}

func (r *MeetingRepo) DeleteParticipant(ctx context.Context, id uuid.UUID) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM meeting_participants WHERE id = $1`, id)
	return err
}

// ═══ NOTES ═══

func (r *MeetingRepo) UpsertNotes(ctx context.Context, n *entity.MeetingNotes) error {
	query := `INSERT INTO meeting_notes (id, meeting_id, purpose, summary, discussion_points, decision, next_step, created_at, updated_at)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9)
		ON CONFLICT (meeting_id) DO UPDATE SET
			purpose = EXCLUDED.purpose,
			summary = EXCLUDED.summary,
			discussion_points = EXCLUDED.discussion_points,
			decision = EXCLUDED.decision,
			next_step = EXCLUDED.next_step,
			updated_at = EXCLUDED.updated_at`
	_, err := r.db.ExecContext(ctx, query,
		n.ID, n.MeetingID, n.Purpose, n.Summary, n.DiscussionPoints, n.Decision, n.NextStep, n.CreatedAt, n.UpdatedAt,
	)
	return err
}

func (r *MeetingRepo) FindNotes(ctx context.Context, meetingID uuid.UUID) (*entity.MeetingNotes, error) {
	n := &entity.MeetingNotes{}
	var purpose, summary, discussion, decision, nextStep sql.NullString
	err := r.db.QueryRowContext(ctx,
		`SELECT id, meeting_id, purpose, summary, discussion_points, decision, next_step, created_at, updated_at FROM meeting_notes WHERE meeting_id = $1`, meetingID,
	).Scan(&n.ID, &n.MeetingID, &purpose, &summary, &discussion, &decision, &nextStep, &n.CreatedAt, &n.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	n.Purpose = purpose.String
	n.Summary = summary.String
	n.DiscussionPoints = discussion.String
	n.Decision = decision.String
	n.NextStep = nextStep.String
	return n, nil
}

// ═══ ACTION PLANS ═══

func (r *MeetingRepo) AddActionPlan(ctx context.Context, a *entity.MeetingActionPlan) error {
	query := `INSERT INTO meeting_action_plans (id, meeting_id, task_name, pic, deadline, status, created_at, updated_at)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8)`
	_, err := r.db.ExecContext(ctx, query,
		a.ID, a.MeetingID, a.TaskName, a.PIC, a.Deadline, a.Status, a.CreatedAt, a.UpdatedAt,
	)
	return err
}

func (r *MeetingRepo) FindActionPlans(ctx context.Context, meetingID uuid.UUID) ([]*entity.MeetingActionPlan, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT id, meeting_id, task_name, pic, deadline, status, created_at, updated_at FROM meeting_action_plans WHERE meeting_id = $1 ORDER BY created_at ASC`, meetingID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []*entity.MeetingActionPlan
	for rows.Next() {
		a := &entity.MeetingActionPlan{}
		if err := rows.Scan(&a.ID, &a.MeetingID, &a.TaskName, &a.PIC, &a.Deadline, &a.Status, &a.CreatedAt, &a.UpdatedAt); err != nil {
			return nil, err
		}
		list = append(list, a)
	}
	return list, nil
}

func (r *MeetingRepo) UpdateActionPlan(ctx context.Context, a *entity.MeetingActionPlan) error {
	query := `UPDATE meeting_action_plans SET task_name=$1, pic=$2, deadline=$3, status=$4, updated_at=$5 WHERE id=$6`
	_, err := r.db.ExecContext(ctx, query, a.TaskName, a.PIC, a.Deadline, a.Status, time.Now(), a.ID)
	return err
}

func (r *MeetingRepo) DeleteActionPlan(ctx context.Context, id uuid.UUID) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM meeting_action_plans WHERE id = $1`, id)
	return err
}

// ═══ FILES ═══

func (r *MeetingRepo) AddFile(ctx context.Context, f *entity.MeetingFile) error {
	query := `INSERT INTO meeting_files (id, meeting_id, file_name, file_path, file_type, uploaded_by, uploaded_at)
		VALUES ($1,$2,$3,$4,$5,$6,$7)`
	_, err := r.db.ExecContext(ctx, query,
		f.ID, f.MeetingID, f.FileName, f.FilePath, f.FileType, f.UploadedBy, f.UploadedAt,
	)
	return err
}

func (r *MeetingRepo) FindFiles(ctx context.Context, meetingID uuid.UUID) ([]*entity.MeetingFile, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT id, meeting_id, file_name, file_path, file_type, uploaded_by, uploaded_at FROM meeting_files WHERE meeting_id = $1 ORDER BY uploaded_at ASC`, meetingID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []*entity.MeetingFile
	for rows.Next() {
		f := &entity.MeetingFile{}
		var uploadedBy sql.NullString
		if err := rows.Scan(&f.ID, &f.MeetingID, &f.FileName, &f.FilePath, &f.FileType, &uploadedBy, &f.UploadedAt); err != nil {
			return nil, err
		}
		if uploadedBy.Valid {
			uid, _ := uuid.Parse(uploadedBy.String)
			f.UploadedBy = &uid
		}
		list = append(list, f)
	}
	return list, nil
}

func (r *MeetingRepo) DeleteFile(ctx context.Context, id uuid.UUID) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM meeting_files WHERE id = $1`, id)
	return err
}
