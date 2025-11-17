package repository

import (
	"fmt"
	"gin-demo/model"

	"gorm.io/gorm"
)

type UserRepository interface {
	CreateUser(user *model.User) error
	FindByEmail(email string) (*model.User, error)
	FindById(id string) (*model.User, error)
	FindAll(email string) []model.User
}

type userRepo struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepo{db}
}

func (r *userRepo) CreateUser(user *model.User) error {
	var existingUser model.User
	if err := r.db.Where("username = ?", user.Username).First(&existingUser).Error; err == nil {
		return fmt.Errorf("username already exists")
	} else if err != gorm.ErrRecordNotFound {
		return err
	}

	if err := r.db.Where("email = ?", user.Email).First(&existingUser).Error; err == nil {
		return fmt.Errorf("email already exists")
	} else if err != gorm.ErrRecordNotFound {
		return err
	}

	return r.db.Create(user).Error
}

func (r *userRepo) FindByEmail(email string) (*model.User, error) {
	var user model.User
	err := r.db.Where("email = ?", email).First(&user).Error
	return &user, err
}

func (r *userRepo) FindById(id string) (*model.User, error) {
	var user model.User
	err := r.db.Where("id = ?", id).First(&user).Error
	return &user, err
}

func (r *userRepo) FindAll(email string) []model.User {
	var users []model.User
	if err := r.db.Model(&model.User{}).
		Select("*").
		Where("email <> ?", email).
		Find(&users).Error; err != nil {
		return nil
	}

	return users
}
