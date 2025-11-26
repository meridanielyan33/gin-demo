package handler_test

import (
	"context"
	"gin-demo/handler"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func setupActorRouter(handler *handler.ActorHandler) *gin.Engine {
	r := gin.Default()
	r.POST("/actors", handler.CreateActor)
	r.PUT("/actors/:id", handler.UpdateActor)
	r.DELETE("/actor/:id", handler.DeleteActor)
	r.GET("/all-actors", handler.GetAllActors)
	r.GET("/actor/:id", handler.GetActor)
	return r
}

func GetDB() *mongo.Client {
	url := "mongodb://localhost:27017/film_actor_director"
	clientOptions := options.Client().ApplyURI(url)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil
	}

	if err := client.Ping(ctx, nil); err != nil {
		return nil
	}
	return client
}

// func TestCreateActor_Success(t *testing.T) {
// 	db := GetDB()
// 	actorRepo := repository.NewActorRepository(db.Database("actor"))
// 	actorSvc := services.NewActorService(actorRepo)
// 	actorHandler := handler.NewActorHandler(*actorSvc)

// }
