package service

import (
	"fmt"
	"time"

	"backend/internal/dto"
	"backend/internal/mapper"
	"backend/internal/model"
	"backend/internal/repository"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

// AuthService defines auth behaviour
type AuthService interface {
	Register(req dto.RegisterRequest) (*dto.AuthResponse, error)
	Login(req dto.LoginRequest) (*dto.AuthResponse, error)
}

type authService struct {
	repo      repository.UserRepository
	jwtSecret string
}

func NewAuthService(repo repository.UserRepository, secret string) AuthService {
	return &authService{repo: repo, jwtSecret: secret}
}

func (s *authService) Register(req dto.RegisterRequest) (*dto.AuthResponse, error) {
	if req.Email == "" || req.Nickname == "" || req.Password == "" {
		return nil, fmt.Errorf("email, nickname and password required")
	}

	if _, err := s.repo.GetByEmail(req.Email); err == nil {
		return nil, fmt.Errorf("email already registered")
	} else if !isRepoNotFound(err) {
		return nil, fmt.Errorf("checking email: %w", err)
	}

	if _, err := s.repo.GetByNickname(req.Nickname); err == nil {
		return nil, fmt.Errorf("nickname already registered")
	} else if !isRepoNotFound(err) {
		return nil, fmt.Errorf("checking nickname: %w", err)
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("hashing password: %w", err)
	}

	user := mapper.MapRegisterRequestToUser(req, string(hash))

	if err := s.repo.Create(user); err != nil {
		return nil, fmt.Errorf("create user: %w", err)
	}

	return s.buildAuthResponse(user)
}

func (s *authService) Login(req dto.LoginRequest) (*dto.AuthResponse, error) {
	if req.EmailOrUsername == "" || req.Password == "" {
		return nil, fmt.Errorf("email_or_username and password required")
	}

	user, err := s.repo.GetByEmail(req.EmailOrUsername)
	if err != nil && isRepoNotFound(err) {
		// try nickname
		user, err = s.repo.GetByNickname(req.EmailOrUsername)
	}
	if err != nil {
		// do not expose whether email/nickname exists
		return nil, fmt.Errorf("invalid credentials")
	}

	if bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)) != nil {
		return nil, fmt.Errorf("invalid credentials")
	}

	return s.buildAuthResponse(user)
}

func (s *authService) buildAuthResponse(user *model.User) (*dto.AuthResponse, error) {
	postsCnt, err := s.repo.GetPostsCount(user.ID)
	if err != nil {
		return nil, fmt.Errorf("count posts: %w", err)
	}
	followersCnt, err := s.repo.GetFollowersCount(user.ID)
	if err != nil {
		return nil, fmt.Errorf("count followers: %w", err)
	}
	followingCnt, err := s.repo.GetFollowingCount(user.ID)
	if err != nil {
		return nil, fmt.Errorf("count following: %w", err)
	}

	token, err := s.buildJWT(user)
	if err != nil {
		return nil, fmt.Errorf("build token: %w", err)
	}

	return &dto.AuthResponse{
		User:        mapper.MapUserToResponseWithCounts(user, postsCnt, followersCnt, followingCnt),
		AccessToken: token,
	}, nil
}

func (s *authService) buildJWT(user *model.User) (string, error) {
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

func isRepoNotFound(err error) bool {
	return err == repository.ErrNotFound || (err != nil && err.Error() == repository.ErrNotFound.Error())
}
