package services

import (
	"gin-demo/model"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Implemented facade structural design pattern
// if the user service code logic changes,
// or our user service starts to need using more service
// this way the code will be more cleaner and flexible
type UserServiceFacade struct {
	userService IUserService
}

func NewUserServiceFacade(userService IUserService) *UserServiceFacade {
	return &UserServiceFacade{
		userService: userService,
	}
}

func (f *UserServiceFacade) Register(user *model.User) error {
	return f.userService.Register(user)
}

func (f *UserServiceFacade) Login(req *model.UserLoginRequest) (*model.UserLoginResponse, error) {
	return f.userService.Login(req)
}

func (f *UserServiceFacade) Logout(req *model.UserLogoutRequest) (*model.UserLogoutResponse, error) {
	return f.userService.Logout(req)
}

func (f *UserServiceFacade) GetUsers(email string) []model.User {
	return f.userService.GetUsers(email)
}

func (f *UserServiceFacade) GetUserById(id primitive.ObjectID) (*model.User, error) {
	return f.userService.GetUserById(id)
}

func (f *UserServiceFacade) GetUserByEmail(email string) (*model.User, error) {
	return f.userService.GetUserByEmail(email)
}
