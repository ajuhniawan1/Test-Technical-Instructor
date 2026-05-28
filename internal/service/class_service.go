package service

import (
	"context"
	"errors"

	"assignment-platform/internal/dto"
	"assignment-platform/internal/model"
	"assignment-platform/internal/repository"
)

// ClassService berisi business logic untuk class/batch.
type ClassService struct {
	classRepo *repository.ClassRepository
}

func NewClassService(classRepo *repository.ClassRepository) *ClassService {
	return &ClassService{classRepo: classRepo}
}

// CreateClass membuat class baru.
func (s *ClassService) CreateClass(ctx context.Context, req dto.CreateClassRequest) (uint64, error) {
	startDate, err := repository.ParseDate(req.StartDate)
	if err != nil {
		return 0, errors.New("start_date harus format YYYY-MM-DD")
	}
	endDate, err := repository.ParseDate(req.EndDate)
	if err != nil {
		return 0, errors.New("end_date harus format YYYY-MM-DD")
	}
	if endDate.Before(startDate) {
		return 0, errors.New("end_date tidak boleh sebelum start_date")
	}

	class := model.Class{
		Name:        req.Name,
		Description: req.Description,
		StartDate:   startDate,
		EndDate:     endDate,
	}

	return s.classRepo.Create(ctx, class)
}

// ListClasses mengambil class dengan pagination.
func (s *ClassService) ListClasses(ctx context.Context, page int, limit int) (dto.PaginatedResponse, error) {
	if page <= 0 {
		page = 1
	}
	if limit <= 0 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}

	offset := (page - 1) * limit
	items, total, err := s.classRepo.ListPaginated(ctx, limit, offset)
	if err != nil {
		return dto.PaginatedResponse{}, err
	}

	totalPages := 0
	if total > 0 {
		totalPages = (total + limit - 1) / limit
	}

	return dto.PaginatedResponse{
		Items: items,
		Meta: dto.PaginationMeta{
			Page:       page,
			Limit:      limit,
			Total:      total,
			TotalPages: totalPages,
		},
	}, nil
}

// GetClassByID mengambil detail class berdasarkan id.
func (s *ClassService) GetClassByID(ctx context.Context, classID uint64) (*model.Class, error) {
	return s.classRepo.FindByID(ctx, classID)
}

func (s *ClassService) AssignTrainer(ctx context.Context, classID uint64, trainerID uint64) error {
	// Pastikan class ada sebelum assign.
	if _, err := s.classRepo.FindByID(ctx, classID); err != nil {
		return err
	}
	return s.classRepo.AssignTrainer(ctx, classID, trainerID)
}

func (s *ClassService) AssignTalent(ctx context.Context, classID uint64, talentID uint64) error {
	// Pastikan class ada sebelum assign.
	if _, err := s.classRepo.FindByID(ctx, classID); err != nil {
		return err
	}
	return s.classRepo.AssignTalent(ctx, classID, talentID)
}
