package main

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
	"ozinshe-final-project/handlers"
	"ozinshe-final-project/middlewares"
	"ozinshe-final-project/repositories"
)

func main() {
	r := gin.Default()

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

	r.Run(":8080")
}

func connectToDb() (*pgxpool.Pool, error) {
	dbUrl := "postgres://postgres:postgrespw@localhost:55000"
	conn, err := pgxpool.New(context.Background(), dbUrl)
	if err != nil {
		return nil, err
	}

	return conn, nil
}
