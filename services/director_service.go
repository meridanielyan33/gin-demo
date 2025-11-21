package services

import (
	"gin-demo/model"
	"gin-demo/repository"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type DirectorService struct {
	repo *repository.DirectorRepository
}

func NewDirectorService(repo *repository.DirectorRepository) *DirectorService {
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
