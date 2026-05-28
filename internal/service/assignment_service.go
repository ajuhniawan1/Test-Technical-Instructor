package service

import (
	"context"
	"errors"
	"time"

	"assignment-platform/internal/dto"
	"assignment-platform/internal/model"
	"assignment-platform/internal/repository"
)

// AssignmentService berisi logic untuk membuat, mengambil, update, dan menutup assignment.
type AssignmentService struct {
	assignmentRepo *repository.AssignmentRepository
	classRepo      *repository.ClassRepository
}

func NewAssignmentService(assignmentRepo *repository.AssignmentRepository, classRepo *repository.ClassRepository) *AssignmentService {
	return &AssignmentService{assignmentRepo: assignmentRepo, classRepo: classRepo}
}

// CreateAssignment membuat assignment untuk suatu class.
func (s *AssignmentService) CreateAssignment(ctx context.Context, classID uint64, userID uint64, role string, req dto.CreateAssignmentRequest) (uint64, error) {
	if err := s.ensureCanManageClass(ctx, classID, userID, role); err != nil {
		return 0, err
	}

	deadline, err := time.Parse(time.RFC3339, req.Deadline)
	if err != nil {
		return 0, errors.New("deadline harus format RFC3339, contoh: 2026-05-30T23:59:00+07:00")
	}

	assignment := model.Assignment{
		ClassID:     classID,
		Title:       req.Title,
		Description: req.Description,
		Deadline:    deadline,
		Status:      "active",
		CreatedBy:   userID,
	}

	return s.assignmentRepo.Create(ctx, assignment)
}

// ListByClass mengambil daftar assignment berdasarkan class.
func (s *AssignmentService) ListByClass(ctx context.Context, classID uint64, userID uint64, role string) ([]model.Assignment, error) {
	if err := s.ensureCanViewClass(ctx, classID, userID, role); err != nil {
		return nil, err
	}
	return s.assignmentRepo.ListByClass(ctx, classID)
}

// GetAssignmentByID mengambil detail assignment.
func (s *AssignmentService) GetAssignmentByID(ctx context.Context, assignmentID uint64, userID uint64, role string) (*model.Assignment, error) {
	assignment, err := s.assignmentRepo.FindByID(ctx, assignmentID)
	if err != nil {
		return nil, err
	}
	if err := s.ensureCanViewClass(ctx, assignment.ClassID, userID, role); err != nil {
		return nil, err
	}
	return assignment, nil
}

// UpdateAssignment mengubah assignment. Admin boleh semua, trainer hanya class yang dia handle.
func (s *AssignmentService) UpdateAssignment(ctx context.Context, assignmentID uint64, userID uint64, role string, req dto.UpdateAssignmentRequest) error {
	assignment, err := s.assignmentRepo.FindByID(ctx, assignmentID)
	if err != nil {
		return err
	}
	if err := s.ensureCanManageClass(ctx, assignment.ClassID, userID, role); err != nil {
		return err
	}

	deadline, err := time.Parse(time.RFC3339, req.Deadline)
	if err != nil {
		return errors.New("deadline harus format RFC3339, contoh: 2026-05-30T23:59:00+07:00")
	}

	status := req.Status
	if status == "" {
		status = assignment.Status
	}
	if status != "active" && status != "closed" && status != "archived" {
		return errors.New("status assignment harus active, closed, atau archived")
	}

	assignment.Title = req.Title
	assignment.Description = req.Description
	assignment.Deadline = deadline
	assignment.Status = status

	return s.assignmentRepo.Update(ctx, *assignment)
}

// CloseAssignment menutup assignment agar tidak bisa disubmit lagi.
func (s *AssignmentService) CloseAssignment(ctx context.Context, assignmentID uint64, userID uint64, role string) error {
	assignment, err := s.assignmentRepo.FindByID(ctx, assignmentID)
	if err != nil {
		return err
	}
	if err := s.ensureCanManageClass(ctx, assignment.ClassID, userID, role); err != nil {
		return err
	}
	return s.assignmentRepo.Close(ctx, assignmentID)
}

func (s *AssignmentService) ensureCanViewClass(ctx context.Context, classID uint64, userID uint64, role string) error {
	switch role {
	case "admin":
		return nil
	case "trainer":
		assigned, err := s.classRepo.IsTrainerAssigned(ctx, classID, userID)
		if err != nil {
			return err
		}
		if !assigned {
			return errors.New("trainer tidak terdaftar pada class ini")
		}
		return nil
	case "talent":
		assigned, err := s.classRepo.IsTalentAssigned(ctx, classID, userID)
		if err != nil {
			return err
		}
		if !assigned {
			return errors.New("talent tidak terdaftar pada class ini")
		}
		return nil
	default:
		return errors.New("role tidak dikenali")
	}
}

func (s *AssignmentService) ensureCanManageClass(ctx context.Context, classID uint64, userID uint64, role string) error {
	if role == "admin" {
		return nil
	}
	if role != "trainer" {
		return errors.New("hanya admin atau trainer yang boleh mengelola assignment")
	}
	assigned, err := s.classRepo.IsTrainerAssigned(ctx, classID, userID)
	if err != nil {
		return err
	}
	if !assigned {
		return errors.New("trainer tidak terdaftar pada class ini")
	}
	return nil
}
