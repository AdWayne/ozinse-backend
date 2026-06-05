package service

import (
	"context"
	"errors"
	"time"
	"ozinse-backend/internal/config"
	"ozinse-backend/internal/model"
	"ozinse-backend/internal/repository"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type AuthService interface {
	Register(ctx context.Context, input model.RegisterInput) error
	Login(ctx context.Context, email, password string) (string, string, error)
	RefreshToken(ctx context.Context, token string) (string, string, error)
	ResetPassword(ctx context.Context, email string) error
	GetProfile(ctx context.Context, userID int) (*model.User, error)
	UpdateProfile(ctx context.Context, userID int, input model.UpdateProfileInput) error
	UpdatePassword(ctx context.Context, userID int, input model.UpdatePasswordInput) error
}

type authService struct {
	repo repository.AuthRepository
	cfg  *config.Config
}

func NewAuthService(repo repository.AuthRepository, cfg *config.Config) AuthService {
	return &authService{repo: repo, cfg: cfg}
}

func (s *authService) Register(ctx context.Context, input model.RegisterInput) error {
	// 1. Проверяем, совпадают ли пароли
	if input.Password != input.RepeatPassword {
		return errors.New("PASSWORDS_DO_NOT_MATCH")
	}

	// 2. Проверяем, занят ли email
	existing, _ := s.repo.GetUserByEmail(ctx, input.Email)
	if existing != nil {
		return errors.New("USER_EXISTS")
	}

	// 3. Хэшируем пароль
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	return s.repo.CreateUser(ctx, input.Email, string(hashedBytes), "")
}

func (s *authService) Login(ctx context.Context, email, password string) (string, string, error) {
	user, err := s.repo.GetUserByEmail(ctx, email)
	if err != nil || user == nil {
		return "", "", errors.New("INVALID_CREDENTIALS")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return "", "", errors.New("INVALID_CREDENTIALS")
	}

	return s.generateTokenPair(ctx, user.ID, user.RoleName)
}

func (s *authService) RefreshToken(ctx context.Context, tokenStr string) (string, string, error) {
	userID, expiresAt, err := s.repo.GetRefreshToken(ctx, tokenStr)
	if err != nil || time.Now().After(expiresAt) {
		_ = s.repo.DeleteRefreshToken(ctx, tokenStr)
		return "", "", errors.New("INVALID_REFRESH_TOKEN")
	}

	user, err := s.repo.GetUserByID(ctx, userID)
	if err != nil {
		return "", "", err
	}

	return s.generateTokenPair(ctx, user.ID, user.RoleName)
}

func (s *authService) ResetPassword(ctx context.Context, email string) error {
	user, err := s.repo.GetUserByEmail(ctx, email)
	if err != nil || user == nil {
		return errors.New("USER_NOT_FOUND")
	}
	return nil
}

func (s *authService) GetProfile(ctx context.Context, userID int) (*model.User, error) {
	return s.repo.GetUserByID(ctx, userID)
}

func (s *authService) UpdateProfile(ctx context.Context, userID int, input model.UpdateProfileInput) error {
	return s.repo.UpdateProfile(ctx, userID, input.FullName, input.Phone, input.BirthDate)
}

func (s *authService) UpdatePassword(ctx context.Context, userID int, input model.UpdatePasswordInput) error {
	user, err := s.repo.GetUserByID(ctx, userID)
	if err != nil {
		return err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(input.OldPassword)); err != nil {
		return errors.New("INCORRECT_OLD_PASSWORD")
	}

	newHashed, err := bcrypt.GenerateFromPassword([]byte(input.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	return s.repo.UpdatePassword(ctx, userID, string(newHashed))
}

func (s *authService) generateTokenPair(ctx context.Context, userID int, role string) (string, string, error) {
	accessClaims := jwt.MapClaims{
		"user_id": userID,
		"role":    role,
		"exp":     time.Now().Add(time.Hour * 60).Unix(),
	}
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	accessStr, err := accessToken.SignedString([]byte(s.cfg.JWTSecret))
	if err != nil {
		return "", "", err
	}

	refreshExpiry := time.Now().Add(time.Hour * 24 * 7)
	refreshClaims := jwt.MapClaims{
		"user_id": userID,
		"exp":     refreshExpiry.Unix(),
	}
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	refreshStr, err := refreshToken.SignedString([]byte(s.cfg.RefreshSecret))
	if err != nil {
		return "", "", err
	}

	err = s.repo.SaveRefreshToken(ctx, userID, refreshStr, refreshExpiry)
	if err != nil {
		return "", "", err
	}

	return accessStr, refreshStr, nil
}