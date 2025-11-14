package routes

import (
	"gin-demo/handler"
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
	jwtStrategy := middleware.NewJWTStrategy(redis_utils.GetRedisClient())
	userRepo := repository.NewUserRepository(db)
	userService := services.NewUserService(userRepo, *jwtStrategy)
	userServiceFacade := services.NewUserServiceFacade(userService)
	userHandler := handler.NewHandler(*userServiceFacade)

	protected := router.Group("/api")
	protected.Use(middleware.AuthMiddleware(*jwtStrategy))

	router.POST("/api/login", userHandler.Login)
	protected.POST("/users", userHandler.Register)
	protected.GET("/users", userHandler.GetUsers)
	protected.GET("/users/me", userHandler.GetAuthenticatedUser)
	protected.GET("/users/:id", userHandler.GetUserById)
	protected.GET("/logout", userHandler.Logout)
	return router
}
