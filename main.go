package main

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
	"ozinshe-final-project/handlers"
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

	r.GET("genres", genreHandlers.HandleFindAll)
	r.GET("genres/:id", genreHandlers.HandleFindById)
	r.POST("genres", genreHandlers.HandleCreate)
	r.PUT("genres/:id", genreHandlers.HandleUpdate)
	r.DELETE("genres/:id", genreHandlers.HandleDelete)

	r.GET("movies", moviesHandler.HandleFindAll)
	r.GET("movies/:id", moviesHandler.HandleFindById)
	r.POST("movies", moviesHandler.HandleCreate)
	r.PUT("movies/:id", moviesHandler.HandleUpdate)
	r.DELETE("movies/:id", moviesHandler.HandleDelete)
	r.PATCH("movies/:id/rate", moviesHandler.HandleSetRating)
	r.PATCH("movies/:id/setWatched", moviesHandler.HandleSetWatched)

	r.GET("watchlist", watchlistHandlers.HandleGetMovies)
	r.POST("watchlist/add", watchlistHandlers.HandleAddMovie)
	r.DELETE("watchlist/remove", watchlistHandlers.HandleRemoveMovie)

	r.GET("users", userHandlers.HandleFindAll)
	r.GET("users/:id", userHandlers.HandleFindById)
	r.POST("users", userHandlers.HandleCreate)
	r.PUT("users/:id", userHandlers.HandleUpdate)
	r.PUT("users/:id/changePassword", userHandlers.HandleChangePassword)
	r.DELETE("users/:id", userHandlers.HandleDelete)

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
