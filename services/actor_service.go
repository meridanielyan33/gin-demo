package services

import (
	"gin-demo/model"
	"gin-demo/repository"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ActorService struct {
	repo *repository.ActorRepository
}

func NewActorService(repo *repository.ActorRepository) *ActorService {
	return &ActorService{repo: repo}
}

func (s *ActorService) Create(actor *model.Actor) (primitive.ObjectID, error) {
	return s.repo.Create(actor)
}

func (s *ActorService) GetByID(id primitive.ObjectID) (*model.Actor, error) {
	return s.repo.GetByID(id)
}

func (s *ActorService) GetAll() ([]model.Actor, error) {
	return s.repo.GetAll()
}

func (s *ActorService) Update(id primitive.ObjectID, update bson.M) error {
	return s.repo.Update(id, update)
}

func (s *ActorService) Delete(id primitive.ObjectID) error {
	return s.repo.Delete(id)
}
