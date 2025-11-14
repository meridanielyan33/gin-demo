package services

import (
	"context"
	"errors"
	"fmt"
	"gin-demo/middleware"
	"gin-demo/model"
	"gin-demo/repository"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type UserService interface {
	Register(user *model.User) error
	Login(loginRequest *model.UserLoginRequest) (*model.UserLoginResponse, error)
	Logout(logoutRequest *model.UserLogoutRequest) (*model.UserLogoutResponse, error)
	GetUsers(email string) []model.UserData
	GetUserById(id string) (*model.UserData, error)
	GetUserByEmail(email string) (*model.UserData, error)
}

type UserData struct {
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
}

type userService struct {
	repo          repository.UserRepository
	tokenStrategy middleware.JWTStrategy
}

func NewUserService(repo repository.UserRepository, tokenStrategy middleware.JWTStrategy) UserService {
	return &userService{
		repo:          repo,
		tokenStrategy: tokenStrategy,
	}
}

func (s *userService) Register(user *model.User) error {
	hashed, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.Password = string(hashed)
	return s.repo.CreateUser(user)
}

func (s *userService) Login(loginRequest *model.UserLoginRequest) (*model.UserLoginResponse, error) {
	userAuth, err := s.repo.FindByEmail(loginRequest.Email)
	if err != nil {
		return nil, errors.New("no such user with specified email")
	}

	err = bcrypt.CompareHashAndPassword([]byte(userAuth.Password), []byte(loginRequest.Password))
	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	accessToken, err := s.tokenStrategy.GenerateAccessToken(context.Background(), loginRequest.Email)
	if err != nil {
		return nil, fmt.Errorf("failed to generate access token: %w", err)
	}

	response := model.UserLoginResponse{
		JWTToken: accessToken,
		Message:  "Welcome to our page dear " + userAuth.Username,
	}
	return &response, nil
}

func (s *userService) Logout(logoutRequest *model.UserLogoutRequest) (*model.UserLogoutResponse, error) {
	user, err := s.repo.FindByEmail(logoutRequest.Email)
	if err != nil {
		return nil, fmt.Errorf("failed to find user: %w", err)
	}
	if user == nil {
		return nil, fmt.Errorf("user %s not found", logoutRequest.Email)
	}

	err = s.tokenStrategy.InvalidateToken(context.Background(), logoutRequest.Email)
	if err != nil {
		return nil, fmt.Errorf("logout failed: %w", err)
	}

	return &model.UserLogoutResponse{
		Message: fmt.Sprintf("User %s logged out successfully", logoutRequest.Email),
	}, nil
}

func (s *userService) GetUsers(email string) []model.UserData {
	return s.repo.FindAll(email)
}

func (s *userService) GetUserById(id string) (*model.UserData, error) {
	return s.repo.FindById(id)
}

func (s *userService) GetUserByEmail(email string) (*model.UserData, error) {
	user, err := s.repo.FindByEmail(email)
	if err != nil {
		return nil, fmt.Errorf("failed to find user: %w", err)
	}
	userData := &model.UserData{
		Username:  user.Username,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
		FirstName: user.FirstName,
		LastName:  user.LastName,
	}
	return userData, nil
}
