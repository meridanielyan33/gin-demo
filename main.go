package main

import (
	"flag"
	"gin-demo/config"
	"gin-demo/redis_utils"
	"gin-demo/routes"
	"log"
)

func main() {
	env := config.Env

	configFile := flag.String("config", "./config/config.json", "Path to the config file")
	flag.Parse()

	if err := config.LoadConfig(*configFile); err != nil {
		log.Fatalf("Error loading config: %v", err)
	}
	err := NewDatabase(config.GetConfig())
	if err != nil {
		log.Fatalf("could not initialize database connection: %s", err)
	}
	db := GetDatabase(config.GetConfig())
	redis_utils.InitRedis(env)
	router := routes.SetupRouter(db)

	router.Run(":8080")
}
