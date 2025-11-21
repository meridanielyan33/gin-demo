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

type ActorHandler struct {
	service services.ActorService
}

func NewActorHandler(service services.ActorService) *ActorHandler {
	return &ActorHandler{service: service}
}

func (h *ActorHandler) CreateActor(c *gin.Context) {
	var raw map[string]interface{}
	if err := c.ShouldBindJSON(&raw); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if birthStr, ok := raw["birthDate"].(string); ok {
		t, err := utils.ParseDate(birthStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": errMsg.TimeFormatWrong})
			return
		}
		raw["birthDate"] = t
	}

	var actor model.Actor
	bsonBytes, _ := bson.Marshal(raw)
	bson.Unmarshal(bsonBytes, &actor)

	_, err := h.service.Create(&actor)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, "Successfully created an actor")
}

func (h *ActorHandler) GetActor(c *gin.Context) {
	idStr := c.Param("id")
	id, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": errMsg.InvalidID})
		return
	}
	actor, err := h.service.GetByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, actor)
}

func (h *ActorHandler) GetAllActors(c *gin.Context) {
	actors, err := h.service.GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, actors)
}

func (h *ActorHandler) UpdateActor(c *gin.Context) {
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

	var updatedActor model.Actor
	jsonBody, _ := json.Marshal(body)
	json.Unmarshal(jsonBody, &updatedActor)

	updateBson := bson.M{}

	if updatedActor.FirstName != "" {
		updateBson["first_name"] = updatedActor.FirstName
	}
	if updatedActor.LastName != "" {
		updateBson["last_name"] = updatedActor.LastName
	}
	if raw, exists := body["birth_date"]; exists {
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

	c.JSON(http.StatusOK, "Successfully updated an actor")
}

func (h *ActorHandler) DeleteActor(c *gin.Context) {
	idStr := c.Param("id")
	id, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": errMsg.InvalidID})
		return
	}

	err = h.service.Delete(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, "Successfully deleted an actor")
}
