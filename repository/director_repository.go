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

type DirectorRepository struct {
	collection *mongo.Collection
}

type IDirectorRepository interface {
	Create(director *model.Director) (primitive.ObjectID, error)
	GetByID(id primitive.ObjectID) (*model.Director, error)
	GetAll() ([]model.Director, error)
	Update(id primitive.ObjectID, update bson.M) error
	Delete(id primitive.ObjectID) error
}

func NewDirectorRepository(db *mongo.Database) IDirectorRepository {
	return &DirectorRepository{collection: db.Collection("directors")}
}

func (r *DirectorRepository) Create(director *model.Director) (primitive.ObjectID, error) {
	director.ID = primitive.NewObjectID()
	_, err := r.collection.InsertOne(context.Background(), director)
	return director.ID, err
}

func (r *DirectorRepository) GetByID(id primitive.ObjectID) (*model.Director, error) {
	var director model.Director
	err := r.collection.FindOne(context.Background(), bson.M{"_id": id}).Decode(&director)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return nil, fmt.Errorf("director not found with given id")
	}
	return &director, err
}

func (r *DirectorRepository) GetAll() ([]model.Director, error) {
	cursor, err := r.collection.Find(context.Background(), bson.M{})
	if err != nil {
		return nil, err
	}
	var directors []model.Director
	if err = cursor.All(context.Background(), &directors); err != nil {
		return nil, err
	}
	return directors, nil
}

func (r *DirectorRepository) Update(id primitive.ObjectID, update bson.M) error {
	_, err := r.collection.UpdateOne(context.Background(), bson.M{"_id": id}, bson.M{"$set": update})
	if errors.Is(err, mongo.ErrNoDocuments) {
		return fmt.Errorf("director not found with given id")
	}
	return err
}

func (r *DirectorRepository) Delete(id primitive.ObjectID) error {
	_, err := r.collection.DeleteOne(context.Background(), bson.M{"_id": id})
	if errors.Is(err, mongo.ErrNoDocuments) {
		return fmt.Errorf("director not found with given id")
	}
	return err
}
