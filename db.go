package main

import (
	"fmt"
	"gin-demo/config"
	"gin-demo/model"
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var db *gorm.DB

func NewDatabase(config *config.Config) (*gorm.DB, error) {
	host := config.Database.Host
	port := config.Database.Port
	user := config.Database.User
	dbname := config.Database.Name
	if host == "" || port == "" || user == "" || dbname == "" {
		log.Fatalf("Missing required database configuration in config.json")
		return nil, fmt.Errorf("missing required database configuration")
	}

	dsn := fmt.Sprintf("%s@tcp(%s:%s)/%s?parseTime=true",
		user, host, port, dbname)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Printf("Error opening database: %v", err)
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		log.Printf("Error getting raw SQL DB instance: %v", err)
		return nil, err
	}
	if err := sqlDB.Ping(); err != nil {
		log.Printf("Error pinging database: %v", err)
		return nil, err
	}

	if err := db.AutoMigrate(&model.User{}); err != nil {
		log.Printf("Error auto-migrating database: %v", err)
		return nil, err
	}

	return db, nil
}

func GetDB() *gorm.DB {
	if db == nil {
		log.Fatal("Database not initialized")
	}
	return db
}
