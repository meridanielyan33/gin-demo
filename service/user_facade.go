package services

import (
	"context"
	"gin-demo/model"
)

// Implemented facade structural design pattern
// if the user service code logic changes,
// or our user service starts to need using more service
// this way the code will be more cleaner and flexible
type UserServiceFacade struct {
	userService UserService
}

func NewUserServiceFacade(userService UserService) *UserServiceFacade {
	return &UserServiceFacade{
		userService: userService,
	}
}

func (f *UserServiceFacade) Register(user *model.User) error {
	return f.userService.Register(user)
}

func (f *UserServiceFacade) Login(ctx context.Context, req *model.UserLoginRequest) (*model.UserLoginResponse, error) {
	return f.userService.Login(req)
}

func (f *UserServiceFacade) Logout(ctx context.Context, req *model.UserLogoutRequest) (*model.UserLogoutResponse, error) {
	return f.userService.Logout(req)
}

func (f *UserServiceFacade) GetUsers(email string) []model.UserData {
	return f.userService.GetUsers(email)
}

func (f *UserServiceFacade) GetUserById(id string) (*model.UserData, error) {
	return f.userService.GetUserById(id)
}

func (f *UserServiceFacade) GetUserByEmail(email string) (*model.UserData, error) {
	return f.userService.GetUserByEmail(email)
}
