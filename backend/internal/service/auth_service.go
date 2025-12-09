package service

import (
	"errors"
	"time"

	"backend/internal/dto"
	"backend/internal/mapper"
	"backend/internal/model"
	"backend/internal/repository"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	repo      repository.UserRepository
	jwtSecret string
}

func NewAuthService(repo repository.UserRepository, jwtSecret string) *AuthService {
	return &AuthService{repo: repo, jwtSecret: jwtSecret}
}

func (s *AuthService) Register(req dto.RegisterRequest) (*dto.AuthResponse, error) {
	if req.Email == "" || req.Nickname == "" || req.Password == "" {
		return nil, errors.New("email, nickname and password required")
	}

	_, err := s.repo.GetByEmail(req.Email)
	if err == nil {
		return nil, errors.New("email already registered")
	}
	if !errors.Is(err, repository.ErrNotFound) {
		return nil, err
	}

	_, err = s.repo.GetByNickname(req.Nickname)
	if err == nil {
		return nil, errors.New("nickname already registered")
	}
	if !errors.Is(err, repository.ErrNotFound) {
		return nil, err
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := mapper.MapRegisterRequestToUser(req, string(hash))

	if err := s.repo.Create(user); err != nil {
		return nil, err
	}

	// build response with counts (all zero on fresh user)
	return s.buildAuthResponse(user)
}

func (s *AuthService) Login(req dto.LoginRequest) (*dto.AuthResponse, error) {
	if req.EmailOrUsername == "" || req.Password == "" {
		return nil, errors.New("email_or_username and password required")
	}

	user, err := s.repo.GetByEmail(req.EmailOrUsername)
	if errors.Is(err, repository.ErrNotFound) {
		user, err = s.repo.GetByNickname(req.EmailOrUsername)
	}
	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	if bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)) != nil {
		return nil, errors.New("invalid credentials")
	}

	return s.buildAuthResponse(user)
}

func (s *AuthService) buildAuthResponse(user *model.User) (*dto.AuthResponse, error) {
	// get counts
	postsCnt, err := s.repo.GetPostsCount(user.ID)
	if err != nil {
		return nil, err
	}
	followersCnt, err := s.repo.GetFollowersCount(user.ID)
	if err != nil {
		return nil, err
	}
	followingCnt, err := s.repo.GetFollowingCount(user.ID)
	if err != nil {
		return nil, err
	}

	token, err := s.buildJWT(user)
	if err != nil {
		return nil, err
	}

	return &dto.AuthResponse{
		User:        mapper.MapUserToResponseWithCounts(user, postsCnt, followersCnt, followingCnt),
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
