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

type ActorRepository struct {
	collection *mongo.Collection
}

type IActorRepository interface {
	Create(actor *model.Actor) (primitive.ObjectID, error)
	GetByID(id primitive.ObjectID) (*model.Actor, error)
	GetByIDs(ids []primitive.ObjectID) ([]model.Actor, error)
	GetAll() ([]model.Actor, error)
	Update(id primitive.ObjectID, update bson.M) error
	Delete(id primitive.ObjectID) error
}

func NewActorRepository(db *mongo.Database) IActorRepository {
	return &ActorRepository{collection: db.Collection("actors")}
}

func (r *ActorRepository) Create(actor *model.Actor) (primitive.ObjectID, error) {
	actor.ID = primitive.NewObjectID()
	_, err := r.collection.InsertOne(context.Background(), actor)
	return actor.ID, err
}

func (r *ActorRepository) GetByID(id primitive.ObjectID) (*model.Actor, error) {
	var actor model.Actor
	err := r.collection.FindOne(context.Background(), bson.M{"_id": id}).Decode(&actor)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return nil, fmt.Errorf("actor not found with given id")
	}
	return &actor, err
}

func (r *ActorRepository) GetByIDs(ids []primitive.ObjectID) ([]model.Actor, error) {
	if len(ids) == 0 {
		return []model.Actor{}, nil
	}
	cursor, err := r.collection.Find(context.Background(), bson.M{
		"_id": bson.M{"$in": ids},
	})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())

	var actors []model.Actor
	if err := cursor.All(context.Background(), &actors); err != nil {
		return nil, err
	}

	return actors, nil
}

func (r *ActorRepository) GetAll() ([]model.Actor, error) {
	cursor, err := r.collection.Find(context.Background(), bson.M{})
	if err != nil {
		return nil, err
	}
	var actors []model.Actor
	if err = cursor.All(context.Background(), &actors); err != nil {
		return nil, err
	}
	return actors, nil
}

func (r *ActorRepository) Update(id primitive.ObjectID, update bson.M) error {
	_, err := r.collection.UpdateOne(context.Background(), bson.M{"_id": id}, bson.M{"$set": update})
	if errors.Is(err, mongo.ErrNoDocuments) {
		return fmt.Errorf("actor not found with given id")
	}
	return err
}

func (r *ActorRepository) Delete(id primitive.ObjectID) error {
	_, err := r.collection.DeleteOne(context.Background(), bson.M{"_id": id})
	if errors.Is(err, mongo.ErrNoDocuments) {
		return fmt.Errorf("actor not found with given id")
	}
	return err
}
