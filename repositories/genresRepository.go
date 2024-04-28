package repositories

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"ozinshe-final-project/models"
)

type GenresRepository struct {
	db *pgxpool.Pool
}

func NewGenresRepository(db *pgxpool.Pool) *GenresRepository {
	return &GenresRepository{db: db}
}

func (r *GenresRepository) FindAll(c context.Context) ([]models.Genre, error) {
	rows, err := r.db.Query(c, "select id, title from genres")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	genres := make([]models.Genre, 0)
	for rows.Next() {
		var genre models.Genre
		if err := rows.Scan(&genre.Id, &genre.Title); err != nil {
			return nil, err
		}
		genres = append(genres, genre)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return genres, nil
}

func (r *GenresRepository) FindById(c context.Context, id int) (models.Genre, error) {
	var genre models.Genre
	err := r.db.QueryRow(c, "select id, title from genres where id = $1", id).Scan(&genre.Id, &genre.Title)
	if err != nil {
		return models.Genre{}, err
	}

	return genre, nil
}

func (r *GenresRepository) Create(c context.Context, genre models.Genre) (int, error) {
	var id int
	err := r.db.QueryRow(c, "insert into genres(title) values($1) returning id", genre.Title).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (r *GenresRepository) Update(c context.Context, id int, genre models.Genre) error {
	_, err := r.db.Exec(c, "update genres set title = $1 where id = $2", genre.Title, id)
	return err
}

func (r *GenresRepository) Delete(c context.Context, id int) error {
	_, err := r.db.Exec(c, "delete from genres where id = $1", id)
	return err
}
