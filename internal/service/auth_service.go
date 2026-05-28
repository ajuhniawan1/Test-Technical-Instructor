package service

import (
	"context"
	"database/sql"
	"errors"

	"assignment-platform/internal/dto"
	"assignment-platform/internal/model"
	"assignment-platform/internal/repository"
	"assignment-platform/internal/utils"
)

// AuthService berisi business logic login dan profile.
type AuthService struct {
	userRepo  *repository.UserRepository
	jwtSecret string
}

func NewAuthService(userRepo *repository.UserRepository, jwtSecret string) *AuthService {
	return &AuthService{userRepo: userRepo, jwtSecret: jwtSecret}
}

// Login memvalidasi email/password dan membuat JWT token.
func (s *AuthService) Login(ctx context.Context, req dto.LoginRequest) (*dto.LoginResponse, error) {
	user, err := s.userRepo.FindByEmail(ctx, req.Email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("email atau password salah")
		}
		return nil, err
	}

	// Demo ini memakai SHA-256 agar mudah dijalankan.
	// Untuk production, gunakan bcrypt.
	if !utils.CheckPasswordDemo(req.Password, user.PasswordHash) {
		return nil, errors.New("email atau password salah")
	}

	token, err := utils.GenerateToken(user.ID, user.Email, user.Role, s.jwtSecret)
	if err != nil {
		return nil, err
	}

	return &dto.LoginResponse{
		Token: token,
		Role:  user.Role,
		Name:  user.Name,
		Email: user.Email,
	}, nil
}

// Me mengambil profile user yang sedang login.
func (s *AuthService) Me(ctx context.Context, userID uint64) (*model.User, error) {
	return s.userRepo.FindByID(ctx, userID)
}
