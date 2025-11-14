package main

import (
	"fmt"
	"gin-demo/config"
	"gin-demo/model"
	"log"
	"sync"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var (
	instance *gorm.DB
	once     sync.Once
	initErr  error
)

// Used creational pattern - singleton in the database creation
func GetDB(config *config.Config) (*gorm.DB, error) {
	once.Do(func() {
		host := config.Database.Host
		port := config.Database.Port
		user := config.Database.User
		dbname := config.Database.Name
		if host == "" || port == "" || user == "" || dbname == "" {
			log.Fatalf("Missing required database configuration in config.json")
			initErr = fmt.Errorf("Missing required database configuration in config.json")
			return
		}

		dsn := fmt.Sprintf("%s@tcp(%s:%s)/%s?parseTime=true",
			user, host, port, dbname)

		db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
		if err != nil {
			log.Printf("Error opening database: %v", err)
			initErr = err
			return
		}

		sqlDB, err := db.DB()
		if err != nil {
			log.Printf("Error getting raw SQL DB instance: %v", err)
			initErr = err
			return
		}
		if err := sqlDB.Ping(); err != nil {
			log.Printf("Error pinging database: %v", err)
			initErr = err
			return
		}

		if err := db.AutoMigrate(&model.User{}); err != nil {
			log.Printf("Error auto-migrating database: %v", err)
			initErr = err
			return
		}
		instance = db
	})
	return instance, initErr
}
