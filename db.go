package main

import (
	"fmt"
	"gin-demo/config"
	"gin-demo/model"
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var instance *gorm.DB

// Used creational pattern - singleton in the database creation
func NewDatabase(cfg *config.Config) error {
	if instance != nil {
		return nil
	}

	host := cfg.Database.Host
	port := cfg.Database.Port
	user := cfg.Database.User
	dbname := cfg.Database.Name

	dsn := fmt.Sprintf("%s@tcp(%s:%s)/%s?parseTime=true",
		user, host, port, dbname)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return err
	}

	if err := db.AutoMigrate(&model.User{}); err != nil {
		return err
	}

	instance = db
	return nil
}

func GetDB() *gorm.DB {
	if instance == nil {
		log.Fatalf("GetDB() called before InitDB()")
	}
	return instance
}
