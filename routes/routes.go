package routes

import (
	"gin-demo/handler"
	"gin-demo/middleware"
	"gin-demo/redis_utils"
	"gin-demo/repository"
	"gin-demo/services"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

func SetupRouter(db *mongo.Database) *gin.Engine {
	router := gin.Default()
	router.SetTrustedProxies([]string{"0.0.0.0/0"})
	jwtStrategy := middleware.NewJWTStrategy(redis_utils.GetRedisClient())
	userRepo := repository.NewUserRepository(db)
	userService := services.NewUserService(userRepo, *jwtStrategy)
	userServiceFacade := services.NewUserServiceFacade(userService)
	userHandler := handler.NewHandler(*userServiceFacade)

	actorRepo := repository.NewActorRepository(db)
	actorService := services.NewActorService(actorRepo)
	actorHandler := handler.NewActorHandler(actorService)

	directorRepo := repository.NewDirectorRepository(db)
	directorService := services.NewDirectorService(directorRepo)
	directorHandler := handler.NewDirectorHandler(directorService)

	movieRepo := repository.NewMovieRepository(db)
	movieHydrator := repository.NewMovieHydrator(directorRepo, actorRepo)
	movieService := services.NewMovieService(movieRepo, movieHydrator)
	movieHandler := handler.NewMovieHandler(movieService)

	protected := router.Group("/api")
	protected.Use(middleware.AuthMiddleware(*jwtStrategy))

	router.POST("/api/login", userHandler.Login)
	protected.POST("/users", userHandler.Register)
	protected.GET("/users", userHandler.GetUsers)
	protected.GET("/users/me", userHandler.GetAuthenticatedUser)
	protected.GET("/users/:id", userHandler.GetUserById)
	protected.GET("/logout", userHandler.Logout)

	protected.POST("/actors", actorHandler.CreateActor)
	protected.PUT("/actors/:id", actorHandler.UpdateActor)
	protected.DELETE("/actor/:id", actorHandler.DeleteActor)
	protected.GET("/all-actors", actorHandler.GetAllActors)
	protected.GET("/actor/:id", actorHandler.GetActor)

	protected.POST("/directors", directorHandler.CreateDirector)
	protected.PUT("/directors/:id", directorHandler.UpdateDirector)
	protected.GET("/all-directors", directorHandler.GetAllDirectors)
	protected.GET("/director/:id", directorHandler.GetDirector)
	protected.DELETE("/director/:id", directorHandler.DeleteDirector)

	protected.POST("/movies", movieHandler.CreateMovie)
	protected.GET("/movie/:id", movieHandler.GetMovie)
	protected.GET("/all-movies", movieHandler.GetAllMovies)
	protected.PUT("/movie/:id", movieHandler.UpdateMovies)
	protected.DELETE("/movie/:id", movieHandler.DeleteMovies)

	//for including the field the url has to look a like this way -->
	// api/actor-movies/692035ff46a473472ef22f5b?field=title,release_year
	//for excluding the field
	// api/actor-movies/692035ff46a473472ef22f5b?exclude=title,release_year
	protected.GET("/director-movies/:directorId", movieHandler.GetMoviesByDirector)
	protected.GET("/actor-movies/:actorId", movieHandler.GetMoviesByActor)

	return router
}
