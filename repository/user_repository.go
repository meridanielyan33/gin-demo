package repository

import (
	"context"
	"errors"
	"fmt"
	"gin-demo/model"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

/*
//relational database usage
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
*/

type UserRepository struct {
	collection *mongo.Collection
}

func NewUserRepository(db *mongo.Database) *UserRepository {
	return &UserRepository{
		collection: db.Collection("users"),
	}
}

func (r *UserRepository) CreateUser(user *model.User) error {
	var existing model.User

	err := r.collection.FindOne(context.Background(), bson.M{"username": user.Username}).Decode(&existing)
	if err == nil {
		return fmt.Errorf("username already exists")
	}
	if !errors.Is(err, mongo.ErrNoDocuments) {
		return err
	}

	err = r.collection.FindOne(context.Background(), bson.M{"email": user.Email}).Decode(&existing)
	if err == nil {
		return fmt.Errorf("email already exists")
	}
	if !errors.Is(mongo.ErrNoDocuments, err) {
		return err
	}

	_, err = r.collection.InsertOne(context.Background(), user)
	return err
}

func (r *UserRepository) FindByEmail(email string) (*model.User, error) {
	var user model.User
	err := r.collection.FindOne(context.Background(), bson.M{"email": email}).Decode(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) FindById(id primitive.ObjectID) (*model.User, error) {
	var user model.User
	err := r.collection.FindOne(context.Background(), bson.M{"_id": id}).Decode(&user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *UserRepository) FindAll(email string) []model.User {
	filter := bson.M{"email": bson.M{"$ne": email}}

	cursor, err := r.collection.Find(context.Background(), filter)
	if err != nil {
		return nil
	}
	defer cursor.Close(context.Background())

	var users []model.User
	err = cursor.All(context.Background(), &users)
	if err != nil {
		return nil
	}

	return users
}
