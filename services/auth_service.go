package services

import (
	"errors"

	"github.com/rakibulbanna/go-fiber-postgres/dtos"
	"github.com/rakibulbanna/go-fiber-postgres/models"
	"github.com/rakibulbanna/go-fiber-postgres/repositories"
	"github.com/rakibulbanna/go-fiber-postgres/utils"
)

type SignUpRequest struct {
	Email    string
	Password string
	Name     string
}

type LoginRequest struct {
	Email    string
	Password string
}

type AuthService struct {
	userRepo  *repositories.UserRepository
	jwtSecret string
}

func NewAuthService(userRepo *repositories.UserRepository, jwtSecret string) *AuthService {
	return &AuthService{
		userRepo:  userRepo,
		jwtSecret: jwtSecret,
	}
}

func (s *AuthService) SignUp(req *SignUpRequest) (*dtos.AuthResponse, error) {
	// Check if user already exists
	existingUser, _ := s.userRepo.FindByEmail(req.Email)
	if existingUser != nil {
		return nil, errors.New("user with this email already exists")
	}

	// Hash password
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, errors.New("failed to hash password")
	}

	// Create user
	user := &models.User{
		Email:    req.Email,
		Password: hashedPassword,
		Name:     req.Name,
	}

	if err := s.userRepo.Create(user); err != nil {
		return nil, errors.New("failed to create user")
	}

	// Generate token
	token, err := utils.GenerateToken(user.ID, user.Email, s.jwtSecret)
	if err != nil {
		return nil, errors.New("failed to generate token")
	}

	return &dtos.AuthResponse{
		Token: token,
		User: dtos.UserResponse{
			ID:    user.ID,
			Email: user.Email,
			Name:  user.Name,
		},
	}, nil
}

func (s *AuthService) Login(req *LoginRequest) (*dtos.AuthResponse, error) {
	// Find user by email
	user, err := s.userRepo.FindByEmail(req.Email)
	if err != nil {
		return nil, errors.New("invalid email or password")
	}

	// Check password
	if !utils.CheckPasswordHash(req.Password, user.Password) {
		return nil, errors.New("invalid email or password")
	}

	// Generate token
	token, err := utils.GenerateToken(user.ID, user.Email, s.jwtSecret)
	if err != nil {
		return nil, errors.New("failed to generate token")
	}

	return &dtos.AuthResponse{
		Token: token,
		User: dtos.UserResponse{
			ID:    user.ID,
			Email: user.Email,
			Name:  user.Name,
		},
	}, nil
}
