package repositories

import (
	"github.com/jackc/pgx/v5"
	"ozinshe-final-project/models"
)

type MoviesRepository struct {
	db *pgx.Conn
}

func NewMoviesRepository(db *pgx.Conn) *MoviesRepository {
	return &MoviesRepository{db: db}
}

func (r *MoviesRepository) FindById(id int) *models.Movie {
	return nil
}
