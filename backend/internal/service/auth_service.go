package service

import (
	"context"
	"errors"
	"strings"
	"time"

	"backend/internal/dto"
	"backend/internal/model"
	"backend/internal/repository"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

// AuthService orchestrates authentication flows.
type AuthService struct {
	repo      repository.UserRepository
	jwtSecret string
}

// NewAuthService builds an AuthService.
func NewAuthService(repo repository.UserRepository, jwtSecret string) *AuthService {
	return &AuthService{repo: repo, jwtSecret: jwtSecret}
}

// Register creates a new user and returns a token.
func (s *AuthService) Register(ctx context.Context, req dto.RegisterRequest) (*dto.AuthResponse, error) {
	if strings.TrimSpace(req.Email) == "" || strings.TrimSpace(req.Nickname) == "" || req.Password == "" {
		return nil, errors.New("email, nickname, and password are required")
	}

	if _, err := s.repo.FindByEmail(ctx, req.Email); err == nil {
		return nil, errors.New("email already registered")
	} else if !errors.Is(err, repository.ErrNotFound) {
		return nil, err
	}

	if _, err := s.repo.FindByNickname(ctx, req.Nickname); err == nil {
		return nil, errors.New("nickname already registered")
	} else if !errors.Is(err, repository.ErrNotFound) {
		return nil, err
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &model.User{
		Email:        req.Email,
		Nickname:     req.Nickname,
		PasswordHash: string(hash),
		FirstName:    req.FirstName,
		LastName:     req.LastName,
		Grade:        req.Grade,
		Major:        req.Major,
		City:         req.City,
	}

	if err := s.repo.CreateUser(ctx, user); err != nil {
		return nil, err
	}

	return s.buildAuthResponse(user)
}

// Login authenticates credentials.
func (s *AuthService) Login(ctx context.Context, req dto.LoginRequest) (*dto.AuthResponse, error) {
	if strings.TrimSpace(req.EmailOrUsername) == "" || req.Password == "" {
		return nil, errors.New("email_or_username and password are required")
	}

	user, err := s.repo.FindByEmail(ctx, req.EmailOrUsername)
	if err != nil {
		if !errors.Is(err, repository.ErrNotFound) {
			return nil, err
		}
		user, err = s.repo.FindByNickname(ctx, req.EmailOrUsername)
		if err != nil {
			if errors.Is(err, repository.ErrNotFound) {
				return nil, errors.New("invalid credentials")
			}
			return nil, err
		}
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		return nil, errors.New("invalid credentials")
	}

	return s.buildAuthResponse(user)
}

func (s *AuthService) buildAuthResponse(user *model.User) (*dto.AuthResponse, error) {
	token, err := s.buildJWT(user)
	if err != nil {
		return nil, err
	}
	return &dto.AuthResponse{
		User:        mapToUserResponse(user),
		AccessToken: token,
	}, nil
}

func (s *AuthService) buildJWT(user *model.User) (string, error) {
	claims := dto.JwtCustomClaims{
		UserID:   user.ID,
		Nickname: user.Nickname,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Subject:   user.Nickname,
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.jwtSecret))
}
