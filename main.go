package main

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
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
	defer conn.Close(context.Background())

	genresRepository := repositories.NewGenresRepository(conn)
	genreHandlers := handlers.NewGenreHandlers(genresRepository)

	moviesRepository := repositories.NewMoviesRepository(conn)
	moviesHandler := handlers.NewMoviesHandler(moviesRepository)

	r.GET("/movies/:id", moviesHandler.HandleFindById)

	r.GET("/genres", genreHandlers.HandleFindAll)
	r.GET("/genres/:id", genreHandlers.HandleFindById)
	r.POST("/genres", genreHandlers.HandleCreate)
	r.PUT("/genres/:id", genreHandlers.HandleUpdate)
	r.DELETE("/genres/:id", genreHandlers.HandleDelete)

	r.Run(":8080")
}

func connectToDb() (*pgx.Conn, error) {
	dbUrl := "postgres://postgres:postgrespw@localhost:55000"
	conn, err := pgx.Connect(context.Background(), dbUrl)
	if err != nil {
		return nil, err
	}

	return conn, nil
}
