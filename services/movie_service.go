package services

import (
	"gin-demo/model"
	"gin-demo/repository"
	"gin-demo/utils"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type MovieService struct {
	repo     *repository.MovieRepository
	hydrator *repository.MovieHydrator
}

func NewMovieService(repo *repository.MovieRepository, hydrator *repository.MovieHydrator) *MovieService {
	return &MovieService{repo: repo, hydrator: hydrator}
}

func (s *MovieService) Create(movie *model.Movie) (primitive.ObjectID, error) {
	return s.repo.Create(movie)
}

func (s *MovieService) GetByID(id primitive.ObjectID) (*model.MovieResponse, error) {
	movie, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}
	hydratorOption := repository.HydrationOptions{
		ForceDirector: true,
		ForceActors:   true,
	}
	err = s.hydrator.Hydrate(movie, nil, hydratorOption)
	if err != nil {
		return nil, err
	}
	movieResponse := &model.MovieResponse{
		Title:         movie.Title,
		ReleaseYear:   movie.ReleaseYear,
		Director:      movie.Director,
		ActorsDetails: movie.ActorsDetails,
	}
	return movieResponse, nil
}

func (s *MovieService) GetAll() ([]model.MovieResponse, error) {
	movies, err := s.repo.GetAll()
	if err != nil {
		return nil, err
	}
	hydratorOption := repository.HydrationOptions{
		ForceDirector: true,
		ForceActors:   true,
	}
	for i := range movies {
		if err := s.hydrator.Hydrate(&movies[i], nil, hydratorOption); err != nil {
			return nil, err
		}
	}

	movieResponses := make([]model.MovieResponse, 0, len(movies))

	for _, movie := range movies {
		movieResponses = append(movieResponses, model.MovieResponse{
			Title:         movie.Title,
			ReleaseYear:   movie.ReleaseYear,
			Director:      movie.Director,
			ActorsDetails: movie.ActorsDetails,
		})
	}

	return movieResponses, nil
}

func (s *MovieService) Update(id primitive.ObjectID, update bson.M) error {
	return s.repo.Update(id, update)
}

func (s *MovieService) Delete(id primitive.ObjectID) error {
	return s.repo.Delete(id)
}

func (s *MovieService) GetByActor(
	actorID primitive.ObjectID,
	pagination *utils.Pagination,
	projection bson.M,
) ([]bson.M, int64, error) {

	totalRows, err := s.repo.CountByActorID(actorID)
	if err != nil {
		return nil, 0, err
	}
	rawMovies, err := s.repo.GetByActor(actorID, pagination, projection)
	if err != nil {
		return nil, 0, err
	}

	opts := repository.HydrationOptions{
		ForceActors: true,
	}

	hydrated, err := s.hydrateRawMovies(rawMovies, projection, opts)
	return hydrated, totalRows, err
}

func (s *MovieService) GetByDirector(
	directorID primitive.ObjectID,
	pagination *utils.Pagination,
	projection bson.M,
) ([]bson.M, int64, error) {

	totalRows, err := s.repo.CountByDirectorID(directorID)
	if err != nil {
		return nil, 0, err
	}
	rawMovies, err := s.repo.GetByDirector(directorID, pagination, projection)
	if err != nil {
		return nil, 0, err
	}

	opts := repository.HydrationOptions{
		ForceDirector: true,
	}

	hydrated, err := s.hydrateRawMovies(rawMovies, projection, opts)
	return hydrated, totalRows, err
}

func (s *MovieService) hydrateRawMovies(
	rawMovies []bson.M,
	projection bson.M,
	opts repository.HydrationOptions,
) ([]bson.M, error) {

	hydratedMovies := make([]bson.M, 0, len(rawMovies))

	for _, raw := range rawMovies {

		var m model.Movie
		bsonBytes, _ := bson.Marshal(raw)
		bson.Unmarshal(bsonBytes, &m)

		if err := s.hydrator.Hydrate(&m, projection, opts); err != nil {
			return nil, err
		}

		if utils.FieldIncluded(projection, "director") || utils.FieldIncluded(projection, "director_id") {
			raw["director"] = m.Director
		}

		if utils.FieldIncluded(projection, "actors") || utils.FieldIncluded(projection, "actors_details") {
			raw["actors_details"] = m.ActorsDetails
		}

		delete(raw, "director_id")
		delete(raw, "actors")

		hydratedMovies = append(hydratedMovies, raw)
	}

	return hydratedMovies, nil
}
