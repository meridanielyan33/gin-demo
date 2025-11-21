package repository

import (
	"context"
	"errors"
	"fmt"
	"gin-demo/model"
	"gin-demo/utils"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MovieRepository struct {
	movies    *mongo.Collection
	actors    *mongo.Collection
	directors *mongo.Collection
}

func NewMovieRepository(db *mongo.Database) *MovieRepository {
	return &MovieRepository{
		movies:    db.Collection("movie"),
		actors:    db.Collection("actors"),
		directors: db.Collection("directors"),
	}
}

func (r *MovieRepository) Create(movie *model.Movie) (primitive.ObjectID, error) {
	movie.ID = primitive.NewObjectID()
	_, err := r.movies.InsertOne(context.Background(), movie)
	return movie.ID, err
}

func (r *MovieRepository) GetByID(id primitive.ObjectID) (*model.Movie, error) {
	var movie model.Movie
	err := r.movies.FindOne(context.Background(), bson.M{"_id": id}).Decode(&movie)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return nil, fmt.Errorf("movie not found with given id")
	}
	return &movie, err
}

func (r *MovieRepository) GetAll() ([]model.Movie, error) {
	cursor, err := r.movies.Find(context.Background(), bson.M{})
	if err != nil {
		return nil, err
	}
	var movies []model.Movie
	if err = cursor.All(context.Background(), &movies); err != nil {
		return nil, err
	}
	return movies, nil
}

func (r *MovieRepository) Update(id primitive.ObjectID, update bson.M) error {
	_, err := r.movies.UpdateOne(context.Background(), bson.M{"_id": id}, bson.M{"$set": update})
	if errors.Is(err, mongo.ErrNoDocuments) {
		return fmt.Errorf("movie not found with given id")
	}
	return err
}

func (r *MovieRepository) Delete(id primitive.ObjectID) error {
	_, err := r.movies.DeleteOne(context.Background(), bson.M{"_id": id})
	if errors.Is(err, mongo.ErrNoDocuments) {
		return fmt.Errorf("movie not found with given id")
	}
	return err
}

func (r *MovieRepository) CountByDirectorID(id primitive.ObjectID) (int64, error) {
	exists, err := r.directors.CountDocuments(
		context.Background(),
		bson.M{"_id": id},
	)
	if err != nil {
		return 0, err
	}
	if exists == 0 {
		return 0, fmt.Errorf("director not found with id: %s", id.Hex())
	}

	filter := bson.M{"director_id": id}
	return r.movies.CountDocuments(context.Background(), filter)
}

func (r *MovieRepository) CountByActorID(id primitive.ObjectID) (int64, error) {
	exists, err := r.actors.CountDocuments(
		context.Background(),
		bson.M{"_id": id},
	)
	if err != nil {
		return 0, err
	}
	if exists == 0 {
		return 0, fmt.Errorf("actor not found with id: %s", id.Hex())
	}

	filter := bson.M{
		"actors": bson.M{"$in": []primitive.ObjectID{id}},
	}

	return r.movies.CountDocuments(context.Background(), filter)
}

func (r *MovieRepository) GetByActor(
	actorID primitive.ObjectID,
	pagination *utils.Pagination,
	projection bson.M,
) ([]bson.M, error) {

	opts := options.Find()
	opts.SetSkip(pagination.GetOffset())
	opts.SetLimit(pagination.Limit)

	if projection != nil {
		opts.SetProjection(projection)
	}

	filter := bson.M{
		"actors": bson.M{"$in": []primitive.ObjectID{actorID}},
	}

	cursor, err := r.movies.Find(context.Background(), filter, opts)
	if err != nil {
		return nil, err
	}

	var movies []bson.M
	if err := cursor.All(context.Background(), &movies); err != nil {
		return nil, err
	}

	return movies, nil
}

func (r *MovieRepository) GetByDirector(
	directorID primitive.ObjectID,
	pagination *utils.Pagination,
	projection bson.M,
) ([]bson.M, error) {

	opts := options.Find()
	opts.SetSkip(pagination.GetOffset())
	opts.SetLimit(pagination.Limit)

	if projection != nil {
		opts.SetProjection(projection)
	}

	cursor, err := r.movies.Find(context.Background(),
		bson.M{"director_id": directorID},
		opts,
	)
	if err != nil {
		return nil, err
	}

	var movies []bson.M
	if err := cursor.All(context.Background(), &movies); err != nil {
		return nil, err
	}

	return movies, nil
}
