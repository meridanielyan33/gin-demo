package services

import (
	"errors"
	"fmt"
	"gin-demo/middleware"
	"gin-demo/model"
	"gin-demo/repository"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

type UserData struct {
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
}

type IUserService interface {
	Register(user *model.User) error
	Login(loginRequest *model.UserLoginRequest) (*model.UserLoginResponse, error)
	Logout(logoutRequest *model.UserLogoutRequest) (*model.UserLogoutResponse, error)
	GetUsers(email string) []model.User
	GetUserById(id primitive.ObjectID) (*model.User, error)
	GetUserByEmail(email string) (*model.User, error)
}

type UserService struct {
	repo          repository.IUserRepository
	tokenStrategy middleware.JWTStrategy
}

func NewUserService(repo repository.IUserRepository, tokenStrategy middleware.JWTStrategy) IUserService {
	return &UserService{
		repo:          repo,
		tokenStrategy: tokenStrategy,
	}
}

func (s *UserService) Register(user *model.User) error {
	hashed, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.Password = string(hashed)
	return s.repo.CreateUser(user)
}

func (s *UserService) Login(loginRequest *model.UserLoginRequest) (*model.UserLoginResponse, error) {
	userAuth, err := s.repo.FindByEmail(loginRequest.Email)
	if err != nil {
		return nil, errors.New("no such user with specified email")
	}

	err = bcrypt.CompareHashAndPassword([]byte(userAuth.Password), []byte(loginRequest.Password))
	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	accessToken, err := s.tokenStrategy.GenerateAccessToken(loginRequest.Email)
	if err != nil {
		return nil, fmt.Errorf("failed to generate access token: %w", err)
	}

	response := model.UserLoginResponse{
		JWTToken: accessToken,
		Message:  "Welcome to our page dear " + userAuth.Username,
	}
	return &response, nil
}

func (s *UserService) Logout(logoutRequest *model.UserLogoutRequest) (*model.UserLogoutResponse, error) {
	user, err := s.repo.FindByEmail(logoutRequest.Email)
	if err != nil {
		return nil, fmt.Errorf("failed to find user: %w", err)
	}
	if user == nil {
		return nil, fmt.Errorf("user %s not found", logoutRequest.Email)
	}

	err = s.tokenStrategy.InvalidateToken(logoutRequest.Email)
	if err != nil {
		return nil, fmt.Errorf("logout failed: %w", err)
	}

	return &model.UserLogoutResponse{
		Message: fmt.Sprintf("User %s logged out successfully", logoutRequest.Email),
	}, nil
}

func (s *UserService) GetUsers(email string) []model.User {
	return s.repo.FindAll(email)
}

func (s *UserService) GetUserById(id primitive.ObjectID) (*model.User, error) {
	return s.repo.FindById(id)
}

func (s *UserService) GetUserByEmail(email string) (*model.User, error) {
	user, err := s.repo.FindByEmail(email)
	if err != nil {
		return nil, fmt.Errorf("failed to find user: %w", err)
	}
	return user, nil
}
