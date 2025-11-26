package handler

import (
	"encoding/json"
	errMsg "gin-demo/errors"
	"gin-demo/model"
	"gin-demo/services"
	"gin-demo/utils"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type DirectorHandler struct {
	service services.IDirectorService
}

func NewDirectorHandler(service services.IDirectorService) *DirectorHandler {
	return &DirectorHandler{service: service}
}

func (h *DirectorHandler) CreateDirector(c *gin.Context) {
	var raw map[string]interface{}
	if err := c.ShouldBindJSON(&raw); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if birthStr, ok := raw["birth_date"].(string); ok {
		t, err := utils.ParseDate(birthStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": errMsg.TimeFormatWrong})
			return
		}
		raw["birth_date"] = t
	}

	var director model.Director
	bsonBytes, _ := bson.Marshal(raw)
	bson.Unmarshal(bsonBytes, &director)

	_, err := h.service.Create(&director)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, "Successfully created a director")
}

func (h *DirectorHandler) GetDirector(c *gin.Context) {
	idHex := c.Param("id")
	id, err := primitive.ObjectIDFromHex(idHex)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": errMsg.InvalidID})
		return
	}

	director, err := h.service.GetByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, director)
}

func (h *DirectorHandler) GetAllDirectors(c *gin.Context) {
	directors, err := h.service.GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, directors)
}

func (h *DirectorHandler) UpdateDirector(c *gin.Context) {
	idHex := c.Param("id")
	id, err := primitive.ObjectIDFromHex(idHex)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": errMsg.InvalidID})
		return
	}

	var body map[string]interface{}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var update model.Director
	jsonBody, _ := json.Marshal(body)
	json.Unmarshal(jsonBody, &update)

	updateBson := bson.M{}

	if update.FirstName != "" {
		updateBson["first_name"] = update.FirstName
	}
	if update.LastName != "" {
		updateBson["last_name"] = update.LastName
	}
	if raw, exists := body["birthDate"]; exists {
		if str, ok := raw.(string); ok && str != "" {
			t, err := time.Parse("2006-01-02", str)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": errMsg.TimeFormatWrong})
				return
			}
			updateBson["birth_date"] = t
		}
	}

	if len(updateBson) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": errMsg.NoFieldsToUpdate})
		return
	}

	if err := h.service.Update(id, updateBson); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, "Successfully updated a director")
}

func (h *DirectorHandler) DeleteDirector(c *gin.Context) {
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
	c.JSON(http.StatusOK, "Successfully deleted a director")
}
