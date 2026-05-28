package service

import (
	"context"
	"database/sql"
	"errors"

	"assignment-platform/internal/dto"
	"assignment-platform/internal/model"
	"assignment-platform/internal/repository"
)

// ReviewService berisi business logic untuk review dan request revision.
type ReviewService struct {
	db         *sql.DB
	reviewRepo *repository.ReviewRepository
	classRepo  *repository.ClassRepository
}

func NewReviewService(db *sql.DB, reviewRepo *repository.ReviewRepository, classRepo *repository.ClassRepository) *ReviewService {
	return &ReviewService{db: db, reviewRepo: reviewRepo, classRepo: classRepo}
}

// ReviewSubmission dipakai trainer/admin untuk memberi nilai/feedback atau meminta revisi.
func (s *ReviewService) ReviewSubmission(ctx context.Context, submissionID uint64, trainerID uint64, role string, req dto.ReviewSubmissionRequest) error {
	if req.Status != "reviewed" && req.Status != "revision_required" {
		return errors.New("status review harus reviewed atau revision_required")
	}
	if req.Score < 0 || req.Score > 100 {
		return errors.New("score harus 0 sampai 100")
	}

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// SELECT ... FOR UPDATE mencegah dua trainer review submission yang sama bersamaan.
	submission, classID, err := s.reviewRepo.GetSubmissionForUpdate(ctx, tx, submissionID)
	if err != nil {
		return err
	}

	if role == "trainer" {
		assigned, err := s.classRepo.IsTrainerAssigned(ctx, classID, trainerID)
		if err != nil {
			return err
		}
		if !assigned {
			return errors.New("trainer tidak terdaftar pada class submission ini")
		}
	}

	if submission.Status != "submitted" && submission.Status != "late" {
		return errors.New("submission tidak dalam status yang bisa direview")
	}

	review := model.SubmissionReview{
		SubmissionID: submissionID,
		TrainerID:    trainerID,
		Score:        req.Score,
		Feedback:     req.Feedback,
		ReviewStatus: req.Status,
	}
	if err := s.reviewRepo.InsertReview(ctx, tx, review); err != nil {
		return err
	}

	if err := s.reviewRepo.UpdateSubmissionStatus(ctx, tx, submissionID, req.Status); err != nil {
		return err
	}

	history := model.SubmissionHistory{
		SubmissionID: submissionID,
		OldStatus:    submission.Status,
		NewStatus:    req.Status,
		Note:         req.Feedback,
		CreatedBy:    trainerID,
	}
	if err := s.reviewRepo.InsertHistory(ctx, tx, history); err != nil {
		return err
	}

	return tx.Commit()
}
