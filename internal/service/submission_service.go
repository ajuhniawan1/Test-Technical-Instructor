package service

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"assignment-platform/internal/dto"
	"assignment-platform/internal/model"
	"assignment-platform/internal/repository"
)

// SubmissionService berisi business logic untuk submit, resubmit, dan progress.
type SubmissionService struct {
	db             *sql.DB
	submissionRepo *repository.SubmissionRepository
	assignmentRepo *repository.AssignmentRepository
	classRepo      *repository.ClassRepository
}

func NewSubmissionService(db *sql.DB, submissionRepo *repository.SubmissionRepository, assignmentRepo *repository.AssignmentRepository, classRepo *repository.ClassRepository) *SubmissionService {
	return &SubmissionService{db: db, submissionRepo: submissionRepo, assignmentRepo: assignmentRepo, classRepo: classRepo}
}

// SubmitAssignment dipakai talent untuk submit assignment pertama kali.
func (s *SubmissionService) SubmitAssignment(ctx context.Context, assignmentID uint64, talentID uint64, req dto.SubmitAssignmentRequest) (uint64, error) {
	if req.GithubURL == "" && req.DeploymentURL == "" {
		return 0, errors.New("github_url atau deployment_url wajib diisi salah satu")
	}

	assignment, err := s.assignmentRepo.FindByID(ctx, assignmentID)
	if err != nil {
		return 0, err
	}

	if assignment.Status != "active" {
		return 0, errors.New("assignment sudah tidak aktif")
	}

	assigned, err := s.classRepo.IsTalentAssigned(ctx, assignment.ClassID, talentID)
	if err != nil {
		return 0, err
	}
	if !assigned {
		return 0, errors.New("talent tidak terdaftar pada class assignment ini")
	}

	status := "submitted"
	if time.Now().After(assignment.Deadline) {
		status = "late"
	}

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return 0, err
	}
	defer tx.Rollback()

	existing, err := s.submissionRepo.FindByAssignmentAndTalent(ctx, tx, assignmentID, talentID)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return 0, err
	}
	if existing != nil {
		return 0, errors.New("submission sudah ada; gunakan PUT /api/v1/submissions/:id untuk resubmit jika status revision_required")
	}

	newSubmission := model.Submission{
		AssignmentID:  assignmentID,
		TalentID:      talentID,
		GithubURL:     req.GithubURL,
		DeploymentURL: req.DeploymentURL,
		Notes:         req.Notes,
		Status:        status,
	}

	submissionID, err := s.submissionRepo.CreateSubmission(ctx, tx, newSubmission)
	if err != nil {
		return 0, err
	}

	history := model.SubmissionHistory{
		SubmissionID: submissionID,
		OldStatus:    "not_submitted",
		NewStatus:    status,
		Note:         "talent submit assignment",
		CreatedBy:    talentID,
	}
	if err := s.submissionRepo.InsertHistory(ctx, tx, history); err != nil {
		return 0, err
	}

	if err := tx.Commit(); err != nil {
		return 0, err
	}
	return submissionID, nil
}

// ResubmitSubmission dipakai talent untuk submit ulang hanya saat status revision_required.
func (s *SubmissionService) ResubmitSubmission(ctx context.Context, submissionID uint64, talentID uint64, req dto.SubmitAssignmentRequest) error {
	if req.GithubURL == "" && req.DeploymentURL == "" {
		return errors.New("github_url atau deployment_url wajib diisi salah satu")
	}

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	submission, _, err := s.submissionRepo.GetSubmissionForUpdate(ctx, tx, submissionID)
	if err != nil {
		return err
	}

	if submission.TalentID != talentID {
		return errors.New("talent tidak boleh resubmit submission milik orang lain")
	}
	if submission.Status != "revision_required" {
		return errors.New("resubmit hanya boleh saat status submission revision_required")
	}

	assignment, err := s.assignmentRepo.FindByID(ctx, submission.AssignmentID)
	if err != nil {
		return err
	}
	if assignment.Status != "active" {
		return errors.New("assignment sudah tidak aktif")
	}

	status := "submitted"
	if time.Now().After(assignment.Deadline) {
		status = "late"
	}

	oldStatus := submission.Status
	if err := s.submissionRepo.UpdateSubmissionLink(ctx, tx, submission.ID, req.GithubURL, req.DeploymentURL, req.Notes, status); err != nil {
		return err
	}

	history := model.SubmissionHistory{
		SubmissionID: submission.ID,
		OldStatus:    oldStatus,
		NewStatus:    status,
		Note:         "talent resubmit setelah revisi",
		CreatedBy:    talentID,
	}
	if err := s.submissionRepo.InsertHistory(ctx, tx, history); err != nil {
		return err
	}

	return tx.Commit()
}

func (s *SubmissionService) ListMySubmissions(ctx context.Context, talentID uint64) ([]map[string]any, error) {
	return s.submissionRepo.ListByTalent(ctx, talentID)
}

func (s *SubmissionService) ListSubmissionsByClass(ctx context.Context, classID uint64, userID uint64, role string) ([]map[string]any, error) {
	if role == "trainer" {
		assigned, err := s.classRepo.IsTrainerAssigned(ctx, classID, userID)
		if err != nil {
			return nil, err
		}
		if !assigned {
			return nil, errors.New("trainer tidak terdaftar pada class ini")
		}
	}
	return s.submissionRepo.ListByClass(ctx, classID)
}

// GetClassProgress menghitung progress langsung dari MySQL.
// Belum menggunakan Redis agar data tetap fresh dan sederhana.
func (s *SubmissionService) GetClassProgress(ctx context.Context, classID uint64, userID uint64, role string) (map[string]any, error) {
	if role == "trainer" {
		assigned, err := s.classRepo.IsTrainerAssigned(ctx, classID, userID)
		if err != nil {
			return nil, err
		}
		if !assigned {
			return nil, errors.New("trainer tidak terdaftar pada class ini")
		}
	}

	talentCount, err := s.classRepo.CountTalentsInClass(ctx, classID)
	if err != nil {
		return nil, err
	}

	assignmentCount, err := s.assignmentRepo.CountByClass(ctx, classID)
	if err != nil {
		return nil, err
	}

	statusCounts, submittedRecords, err := s.submissionRepo.CountStatusByClass(ctx, classID)
	if err != nil {
		return nil, err
	}

	totalExpected := talentCount * assignmentCount
	notSubmitted := totalExpected - submittedRecords
	if notSubmitted < 0 {
		notSubmitted = 0
	}

	return map[string]any{
		"class_id":             classID,
		"total_talents":        talentCount,
		"total_assignments":    assignmentCount,
		"expected_submissions": totalExpected,
		"not_submitted":        notSubmitted,
		"submitted":            statusCounts["submitted"],
		"reviewed":             statusCounts["reviewed"],
		"revision_required":    statusCounts["revision_required"],
		"late":                 statusCounts["late"],
	}, nil
}

// GetTalentProgress menghitung progress satu talent pada semua class yang dia ikuti.
func (s *SubmissionService) GetTalentProgress(ctx context.Context, talentID uint64, requesterID uint64, role string) (map[string]any, error) {
	if role == "talent" && requesterID != talentID {
		return nil, errors.New("talent hanya boleh melihat progress miliknya sendiri")
	}

	classCount, err := s.classRepo.CountClassesForTalent(ctx, talentID)
	if err != nil {
		return nil, err
	}

	assignmentCount, err := s.assignmentRepo.CountForTalent(ctx, talentID)
	if err != nil {
		return nil, err
	}

	statusCounts, submittedRecords, err := s.submissionRepo.CountStatusByTalent(ctx, talentID)
	if err != nil {
		return nil, err
	}

	notSubmitted := assignmentCount - submittedRecords
	if notSubmitted < 0 {
		notSubmitted = 0
	}

	return map[string]any{
		"talent_id":            talentID,
		"total_classes":        classCount,
		"total_assignments":    assignmentCount,
		"expected_submissions": assignmentCount,
		"not_submitted":        notSubmitted,
		"submitted":            statusCounts["submitted"],
		"reviewed":             statusCounts["reviewed"],
		"revision_required":    statusCounts["revision_required"],
		"late":                 statusCounts["late"],
	}, nil
}
