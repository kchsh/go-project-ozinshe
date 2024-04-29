package main

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/spf13/viper"
	"log"
	"ozinshe-final-project/config"
	"ozinshe-final-project/handlers"
	"ozinshe-final-project/middlewares"
	"ozinshe-final-project/repositories"
)

func main() {
	r := gin.Default()

	err := loadConfigs()
	if err != nil {
		log.Fatalf("Error reading config file. %s", err)
	}

	conn, err := connectToDb()
	if err != nil {
		log.Fatal("Unable to connect to db")
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

	r.Run(fmt.Sprintf(":%s", config.Config.AppPort))
}

func connectToDb() (*pgxpool.Pool, error) {
	conn, err := pgxpool.New(context.Background(), config.Config.DbConnectionString)
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
