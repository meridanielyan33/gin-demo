package routes

import (
	"gin-demo/controller"
	"gin-demo/middleware"
	"gin-demo/redis_utils"
	"gin-demo/repository"
	services "gin-demo/service"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupRouter(db *gorm.DB) *gin.Engine {
	router := gin.Default()
	router.SetTrustedProxies([]string{"0.0.0.0/0"})
	userRepo := repository.NewUserRepository(db)
	userService := services.NewUserService(userRepo, redis_utils.GetRedisClient())
	userHandler := controller.NewHandler(userService)

	protected := router.Group("/api")
	protected.Use(middleware.AuthMiddleware(redis_utils.GetRedisClient()))

	router.POST("/users", userHandler.Register)
	router.POST("/login", userHandler.Login)
	protected.GET("/users", userHandler.GetUsers)
	protected.GET("/users/me", userHandler.GetAuthenticatedUser)
	protected.GET("/users/:id", userHandler.GetUserById)
	protected.GET("/logout", userHandler.Logout)
	return router
}
