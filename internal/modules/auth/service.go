package auth

import (
	"errors"

	"github.com/rakibulbanna/go-fiber-postgres/dtos"
	"github.com/rakibulbanna/go-fiber-postgres/models"
	"github.com/rakibulbanna/go-fiber-postgres/utils"
	"gorm.io/gorm"
)

type Service struct {
	db        *gorm.DB
	jwtSecret string
}

func NewService(db *gorm.DB, jwtSecret string) *Service {
	return &Service{
		db:        db,
		jwtSecret: jwtSecret,
	}
}

func (s *Service) SignUp(req *dtos.SignUpRequest) (*dtos.AuthResponse, error) {
	// Check if user already exists
	var existingUser models.User
	if err := s.db.Where("email = ?", req.Email).First(&existingUser).Error; err == nil {
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

	if err := s.db.Create(user).Error; err != nil {
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

func (s *Service) Login(req *dtos.LoginRequest) (*dtos.AuthResponse, error) {
	// Find user by email
	var user models.User
	if err := s.db.Where("email = ?", req.Email).First(&user).Error; err != nil {
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

func (s *Service) FindUserByID(id uint) (*models.User, error) {
	var user models.User
	if err := s.db.First(&user, id).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

