package service

import (
	"errors"

	"projectwebcurhat/config/pkg/errs"
	"projectwebcurhat/config/pkg/token"
	"projectwebcurhat/contract"
	"projectwebcurhat/database"
	"projectwebcurhat/dto"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type authService struct {
	repo *contract.Repository
}

func NewAuthService(repo *contract.Repository) contract.AuthService {
	return &authService{repo: repo}
}

func (s *authService) Register(payload *dto.RegisterRequest) (*dto.AuthResponse, error) {
	// Check if email already exists
	_, err := s.repo.User.GetUserByEmail(payload.Email)
	if err == nil {
		return nil, errs.BadRequest("Email already registered")
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errs.InternalServerError("Failed to check email")
	}

	// Check if username already exists
	_, err = s.repo.User.GetUserByUsername(payload.Username)
	if err == nil {
		return nil, errs.BadRequest("Username already taken")
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errs.InternalServerError("Failed to check username")
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(payload.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, errs.InternalServerError("Failed to hash password")
	}

	// Create user
	user := &database.User{
		Username: payload.Username,
		Email:    payload.Email,
		Password: string(hashedPassword),
	}

	createdUser, err := s.repo.User.CreateUser(user)
	if err != nil {
		return nil, errs.InternalServerError("Failed to create user")
	}

	// Generate JWT token
	tokenString, err := token.GenerateToken(createdUser.ID, createdUser.Username, createdUser.Email)
	if err != nil {
		return nil, errs.InternalServerError("Failed to generate token")
	}

	return &dto.AuthResponse{
		Token: tokenString,
		User: dto.UserProfile{
			ID:       createdUser.ID,
			Username: createdUser.Username,
			Email:    createdUser.Email,
			IsOnline: createdUser.IsOnline,
		},
	}, nil
}

func (s *authService) Login(payload *dto.LoginRequest) (*dto.AuthResponse, error) {
	// Find user by email
	user, err := s.repo.User.GetUserByEmail(payload.Email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errs.Unauthorized("Invalid email or password")
		}
		return nil, errs.InternalServerError("Failed to find user")
	}

	// Compare password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(payload.Password)); err != nil {
		return nil, errs.Unauthorized("Invalid email or password")
	}

	// Generate JWT token
	tokenString, err := token.GenerateToken(user.ID, user.Username, user.Email)
	if err != nil {
		return nil, errs.InternalServerError("Failed to generate token")
	}

	return &dto.AuthResponse{
		Token: tokenString,
		User: dto.UserProfile{
			ID:       user.ID,
			Username: user.Username,
			Email:    user.Email,
			IsOnline: user.IsOnline,
		},
	}, nil
}

func (s *authService) GetProfile(userID int) (*dto.UserProfile, error) {
	user, err := s.repo.User.GetUserByID(userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errs.NotFound("User not found")
		}
		return nil, errs.InternalServerError("Failed to get user")
	}

	return &dto.UserProfile{
		ID:       user.ID,
		Username: user.Username,
		Email:    user.Email,
		IsOnline: user.IsOnline,
	}, nil
}
