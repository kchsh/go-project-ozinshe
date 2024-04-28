package repositories

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"ozinshe-final-project/models"
)

type MoviesRepository struct {
	db *pgxpool.Pool
}

func NewMoviesRepository(db *pgxpool.Pool) *MoviesRepository {
	return &MoviesRepository{db: db}
}

func (r *MoviesRepository) FindAll(c context.Context) ([]models.Movie, error) {
	rows, err := r.db.Query(c, "select id, title, description, date_of_release, director, rating, trailer_url, poster_url from movies")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	movies := make([]models.Movie, 0)
	for rows.Next() {
		var movie models.Movie
		if err := rows.Scan(&movie.Id, &movie.Title, &movie.Description, &movie.DateOfRelease, &movie.Director, &movie.Rating, &movie.TrailerUrl, &movie.PosterUrl); err != nil {
			return nil, err
		}

		genres, err := r.getGenresForMovie(c, movie.Id)
		if err != nil {
			return nil, err
		}
		movie.Genres = genres
		movies = append(movies, movie)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return movies, nil
}

func (r *MoviesRepository) FindById(c context.Context, id int) (models.Movie, error) {
	var movie models.Movie
	err := r.db.QueryRow(
		c,
		`
select id, 
       title, 
       description, 
       date_of_release, 
       director, 
       rating, 
       trailer_url, 
       poster_url 
from movies 
where id = $1
`,
		id).Scan(
		&movie.Id,
		&movie.Title,
		&movie.Description,
		&movie.DateOfRelease,
		&movie.Director,
		&movie.Rating,
		&movie.TrailerUrl,
		&movie.PosterUrl)
	if err != nil {
		return models.Movie{}, err
	}

	genres, err := r.getGenresForMovie(c, movie.Id)
	if err != nil {
		return models.Movie{}, err
	}
	movie.Genres = genres
	return movie, nil
}

func (r *MoviesRepository) Create(c context.Context, movie models.Movie) (int, error) {
	var id int
	err := r.db.QueryRow(
		c,
		`
insert into movies(title, description, date_of_release, director, trailer_url, poster_url) 
values($1, $2, $3, $4, $5, $6) 
returning id`,
		movie.Title,
		movie.Description,
		movie.DateOfRelease,
		movie.Director,
		movie.TrailerUrl,
		movie.PosterUrl,
	).Scan(&id)
	if err != nil {
		return 0, err
	}

	for _, genre := range movie.Genres {
		_, err := r.db.Exec(c, "insert into movie_genres(movie_id, genre_id) values($1, $2)", id, genre.Id)
		if err != nil {
			return 0, err
		}
	}

	return id, nil
}

func (r *MoviesRepository) Update(c context.Context, id int, movie models.Movie) error {
	_, err := r.db.Exec(
		c,
		`
update movies 
set title = $1, 
    description = $2, 
    date_of_release = $3, 
    director = $4, 
    trailer_url = $5, 
    poster_url = $6 
where id = $7
`,
		movie.Title,
		movie.Description,
		movie.DateOfRelease,
		movie.Director,
		movie.TrailerUrl,
		movie.PosterUrl,
		id)
	if err != nil {
		return err
	}

	_, err = r.db.Exec(c, "delete from movie_genres where movie_id = $1", id)
	if err != nil {
		return err
	}
	for _, genre := range movie.Genres {
		_, err := r.db.Exec(c, "insert into movie_genres(movie_id, genre_id) values($1, $2)", id, genre.Id)
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *MoviesRepository) Delete(c context.Context, id int) error {
	_, err := r.db.Exec(c, "delete from movie_genres where movie_id = $1", id)
	if err != nil {
		return err
	}

	_, err = r.db.Exec(c, "delete from movies where id = $1", id)
	if err != nil {
		return err
	}

	return err
}

func (r *MoviesRepository) SetRating(c context.Context, movieId int, rating int) error {
	_, err := r.db.Exec(c, "update movies set rating = $1 where id = $2", rating, movieId)
	if err != nil {
		return err
	}

	return nil
}

func (r *MoviesRepository) SetWatched(c context.Context, movieId int, isWatched bool) error {
	_, err := r.db.Exec(c, "update movies set is_watched = $1 where id = $2", isWatched, movieId)
	if err != nil {
		return err
	}

	return nil
}

func (r *MoviesRepository) getGenresForMovie(c context.Context, movieID int) ([]models.Genre, error) {
	rows, err := r.db.Query(
		c,
		`
select g.id, g.title 
from genres g
join movie_genres mg on mg.genre_id = g.id
join movies m on m.id = $1
`,
		movieID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var genres []models.Genre
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
