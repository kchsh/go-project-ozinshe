package main

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/spf13/viper"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"log"
	"ozinshe-final-project/config"
	"ozinshe-final-project/docs"
	"ozinshe-final-project/handlers"
	"ozinshe-final-project/middlewares"
	"ozinshe-final-project/repositories"
)

// @title           Ozinshe API
// @version         1.0
// @description     This is a sample server celler server.
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  support@swagger.io

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8080
// @BasePath  /

// @securityDefinitions.apikey Bearer
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.

// @externalDocs.description  OpenAPI
// @externalDocs.url          https://swagger.io/resources/open-api/
func main() {
	r := gin.Default()

	err := loadConfigs()
	if err != nil {
		log.Fatal("Error reading config file", err)
	}

	conn, err := connectToDb()
	if err != nil {
		log.Fatal("Unable to connect to db", err)
	}
	defer conn.Close()

	genresRepository := repositories.NewGenresRepository(conn)
	genreHandlers := handlers.NewGenreHandlers(genresRepository)
	moviesRepository := repositories.NewMoviesRepository(conn)
	moviesHandler := handlers.NewMoviesHandler(moviesRepository, genresRepository)
	watchlistRepository := repositories.NewWatchlistRepository(conn)
	watchlistHandlers := handlers.NewWatchlistHandler(moviesRepository, watchlistRepository)
	usersRepository := repositories.NewUsersRepository(conn)
	userHandlers := handlers.NewUserHandlers(usersRepository)
	authHandlers := handlers.NewAuthHandlers(usersRepository)
	imageHandlers := handlers.NewImageHandlers()

	authorized := r.Group("/")
	authorized.Use(middlewares.AuthMiddleware)

	authorized.GET("genres", genreHandlers.HandleFindAll)
	authorized.GET("genres/:id", genreHandlers.HandleFindById)
	authorized.POST("genres", genreHandlers.HandleCreate)
	authorized.PUT("genres/:id", genreHandlers.HandleUpdate)
	authorized.DELETE("genres/:id", genreHandlers.HandleDelete)

	authorized.GET("movies", moviesHandler.HandleFindAll)
	authorized.GET("movies/:id", moviesHandler.HandleFindById)
	authorized.POST("movies", moviesHandler.HandleCreate)
	authorized.PUT("movies/:id", moviesHandler.HandleUpdate)
	authorized.DELETE("movies/:id", moviesHandler.HandleDelete)
	authorized.PATCH("movies/:id/rate", moviesHandler.HandleSetRating)
	authorized.PATCH("movies/:id/setWatched", moviesHandler.HandleSetWatched)

	authorized.GET("watchlist", watchlistHandlers.HandleGetMovies)
	authorized.POST("watchlist/add", watchlistHandlers.HandleAddMovie)
	authorized.DELETE("watchlist/remove", watchlistHandlers.HandleRemoveMovie)

	authorized.GET("users", userHandlers.HandleFindAll)
	authorized.GET("users/:id", userHandlers.HandleFindById)
	authorized.POST("users", userHandlers.HandleCreate)
	authorized.PUT("users/:id", userHandlers.HandleUpdate)
	authorized.PUT("users/:id/changePassword", userHandlers.HandleChangePassword)
	authorized.DELETE("users/:id", userHandlers.HandleDelete)

	authorized.GET("auth/userInfo", authHandlers.HandleGetUserInfo)
	authorized.POST("auth/signOut", authHandlers.HandleSignOut)

	unauthorized := r.Group("")
	unauthorized.POST("auth/signIn", authHandlers.HandleSignIn)
	unauthorized.GET("images", imageHandlers.HandleGetImageById)

	docs.SwaggerInfo.BasePath = "/"
	unauthorized.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	r.Run(fmt.Sprintf(":%s", config.Config.AppPort))
}

func connectToDb() (*pgxpool.Pool, error) {
	conn, err := pgxpool.New(context.Background(), config.Config.DbConnectionString)
	if err != nil {
		return nil, err
	}
	err = conn.Ping(context.Background())
	if err != nil {
		return nil, err
	}

	return conn, nil
}

func loadConfigs() error {
	viper.SetConfigFile(".env")

	err := viper.ReadInConfig()
	if err != nil {
		return err
	}

	var mapConfig config.MapConfig
	err = viper.Unmarshal(&mapConfig)
	if err != nil {
		return err
	}

	config.Config = &mapConfig

	return nil
}
