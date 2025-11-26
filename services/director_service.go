package services

import (
	"gin-demo/model"
	"gin-demo/repository"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type IDirectorService interface {
	Create(director *model.Director) (primitive.ObjectID, error)
	GetByID(id primitive.ObjectID) (*model.Director, error)
	GetAll() ([]model.Director, error)
	Update(id primitive.ObjectID, update bson.M) error
	Delete(id primitive.ObjectID) error
}

type DirectorService struct {
	repo repository.IDirectorRepository
}

func NewDirectorService(repo repository.IDirectorRepository) IDirectorService {
	return &DirectorService{repo: repo}
}

func (d *DirectorService) Create(director *model.Director) (primitive.ObjectID, error) {
	return d.repo.Create(director)
}

func (d *DirectorService) GetByID(id primitive.ObjectID) (*model.Director, error) {
	return d.repo.GetByID(id)
}

func (d *DirectorService) GetAll() ([]model.Director, error) {
	return d.repo.GetAll()
}

func (d *DirectorService) Update(id primitive.ObjectID, update bson.M) error {
	return d.repo.Update(id, update)
}

func (d *DirectorService) Delete(id primitive.ObjectID) error {
	return d.repo.Delete(id)
}
