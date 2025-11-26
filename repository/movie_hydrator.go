package repository

import (
	"gin-demo/model"
	"gin-demo/utils"

	"go.mongodb.org/mongo-driver/bson"
)

type HydrationOptions struct {
	ForceActors   bool
	ForceDirector bool
}

type MovieHydrator struct {
	DirectorsRepo IDirectorRepository
	ActorsRepo    IActorRepository
}

func NewMovieHydrator(dRepo IDirectorRepository, aRepo IActorRepository) *MovieHydrator {
	return &MovieHydrator{
		DirectorsRepo: dRepo,
		ActorsRepo:    aRepo,
	}
}

func (h *MovieHydrator) Hydrate(
	m *model.Movie,
	projection bson.M,
	opts HydrationOptions,
) error {

	hydrateDirector :=
		opts.ForceDirector ||
			projection == nil ||
			utils.FieldIncluded(projection, "director") ||
			utils.FieldIncluded(projection, "director_id")

	if hydrateDirector && !m.DirectorID.IsZero() {
		director, err := h.DirectorsRepo.GetByID(m.DirectorID)
		if err != nil {
			return err
		}
		m.Director = director
	}

	// --- ACTORS ---
	hydrateActors :=
		opts.ForceActors ||
			projection == nil ||
			utils.FieldIncluded(projection, "actors") ||
			utils.FieldIncluded(projection, "actors_details")

	if hydrateActors && len(m.Actors) > 0 {
		actors, err := h.ActorsRepo.GetByIDs(m.Actors)
		if err != nil {
			return err
		}
		m.ActorsDetails = actors
	}

	return nil
}
