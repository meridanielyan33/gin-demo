package handler

import (
	errMsg "gin-demo/errors"
	"gin-demo/model"
	"gin-demo/services"
	"gin-demo/utils"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type MovieHandler struct {
	service *services.MovieService
}

func NewMovieHandler(service *services.MovieService) *MovieHandler {
	return &MovieHandler{service: service}
}

func (h *MovieHandler) CreateMovie(c *gin.Context) {
	var movie model.Movie
	if err := c.ShouldBindJSON(&movie); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	_, err := h.service.Create(&movie)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, "Successfully created a movie")
}

func (h *MovieHandler) GetMovie(c *gin.Context) {
	idHex := c.Param("id")
	id, err := primitive.ObjectIDFromHex(idHex)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": errMsg.InvalidID})
		return
	}

	movie, err := h.service.GetByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, movie)
}

func (h *MovieHandler) GetAllMovies(c *gin.Context) {
	movies, err := h.service.GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, movies)
}

func (h *MovieHandler) UpdateMovies(c *gin.Context) {
	idHex := c.Param("id")
	id, err := primitive.ObjectIDFromHex(idHex)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": errMsg.InvalidID})
		return
	}

	var update model.Movie
	if err := c.ShouldBindJSON(&update); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updateBson := bson.M{}
	if update.Title != "" {
		updateBson["title"] = update.Title
	}
	if update.ReleaseYear != 0 {
		updateBson["release_year"] = update.ReleaseYear
	}
	if update.DirectorID != primitive.NilObjectID {
		updateBson["director_id"] = update.DirectorID
	}
	if len(update.Actors) > 0 {
		updateBson["actors"] = update.Actors
	}

	if len(updateBson) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": errMsg.NoFieldsToUpdate})
		return
	}

	if err := h.service.Update(id, updateBson); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, "Successfully updated a movie")
}

func (h *MovieHandler) DeleteMovies(c *gin.Context) {
	idHex := c.Param("id")
	id, err := primitive.ObjectIDFromHex(idHex)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": errMsg.InvalidID})
		return
	}

	if err := h.service.Delete(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, "Successfully deleted a movie")
}

func (h *MovieHandler) GetMoviesByDirector(c *gin.Context) {
	directorHex := c.Param("directorId")
	directorID, err := primitive.ObjectIDFromHex(directorHex)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": errMsg.InvalidDirectorID})
		return
	}

	page, _ := strconv.ParseInt(c.DefaultQuery("page", "1"), 10, 64)
	limit, _ := strconv.ParseInt(c.DefaultQuery("limit", "10"), 10, 64)
	pagination := utils.NewPagination(page, limit)

	fieldsToInclude := c.DefaultQuery("fields", "")
	fieldsToExclude := c.DefaultQuery("exclude", "")
	projection := utils.BuildProjection(fieldsToInclude, fieldsToExclude)

	movies, totalRows, err := h.service.GetByDirector(directorID, pagination, projection)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	pagination.SetTotal(totalRows)

	c.JSON(http.StatusOK, gin.H{
		"movies": movies,
		"total":  pagination.TotalRows,
		"page":   pagination.Page,
		"limit":  pagination.Limit,
	})
}

func (h *MovieHandler) GetMoviesByActor(c *gin.Context) {
	actorHex := c.Param("actorId")
	actorID, err := primitive.ObjectIDFromHex(actorHex)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": errMsg.InvalidActorID})
		return
	}

	skip, _ := strconv.Atoi(c.DefaultQuery("skip", "0"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	pagination := utils.NewPagination(int64(skip), int64(limit))
	fieldsToInclude := c.DefaultQuery("fields", "")
	fieldsToExclude := c.DefaultQuery("exclude", "")
	projection := utils.BuildProjection(fieldsToInclude, fieldsToExclude)

	movies, totalRows, err := h.service.GetByActor(actorID, pagination, projection)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	pagination.SetTotal(totalRows)

	c.JSON(http.StatusOK, gin.H{
		"movies": movies,
		"total":  pagination.TotalRows,
		"page":   pagination.Page,
		"limit":  pagination.Limit,
	})
}
