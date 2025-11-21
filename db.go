package main

import (
	"context"
	"fmt"
	"gin-demo/config"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var instance *mongo.Client

// Used creational pattern - singleton in the database creation
func NewDatabase(cfg *config.Config) error {
	if instance != nil {
		return nil
	}

	uri := fmt.Sprintf("mongodb://%s:%s", cfg.Database.Host, cfg.Database.Port)

	clientOptions := options.Client().ApplyURI(uri)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return err
	}

	if err := client.Ping(ctx, nil); err != nil {
		return err
	}

	instance = client
	return nil
}

func GetDB() *mongo.Client {
	if instance == nil {
		log.Fatalf("GetDB() called before NewDatabase()")
	}
	return instance
}

func GetDatabase(cfg *config.Config) *mongo.Database {
	return GetDB().Database(cfg.Database.DbName)
}
