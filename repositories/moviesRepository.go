package repositories

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"ozinshe-final-project/models"
	"strconv"
)

type MoviesRepository struct {
	db *pgxpool.Pool
}

func NewMoviesRepository(db *pgxpool.Pool) *MoviesRepository {
	return &MoviesRepository{db: db}
}

func (r *MoviesRepository) FindAll(c context.Context, filters models.MovieFilters) ([]models.Movie, error) {
	sql :=
		`
select m.id, 
       m.title, 
       m.description, 
       m.date_of_release, 
       m.director, 
       m.rating, 
       m.trailer_url, 
       m.poster_url,
       g.id,
       g.title
from movies m 
join movie_genres mg on mg.movie_id = m.id
join genres g on g.id = mg.genre_id
where 1 = 1`

	params := pgx.NamedArgs{}

	if filters.SearchTerm != "" {
		sql = fmt.Sprintf("%s and m.title ilike @s", sql)
		params["s"] = fmt.Sprintf("%%%s%%", filters.SearchTerm)
	}
	if filters.IsWatched != "" {
		isWatched, _ := strconv.ParseBool(filters.IsWatched)

		sql = fmt.Sprintf("%s and m.is_watched = @isWatched", sql)
		params["isWatched"] = isWatched
	}
	if len(filters.GenreIds) > 0 {
		sql = fmt.Sprintf("%s and g.id = any(@genreIds)")
		params["genreIds"] = filters.GenreIds
	}

	if filters.Sort != "" {
		o := "asc"

		// If reverse order
		if string(filters.Sort[0]) == "-" {
			o = "desc"
			filters.Sort = filters.Sort[1:]
		}

		identifier := pgx.Identifier{filters.Sort}
		q := fmt.Sprintf("order by %s %s", identifier.Sanitize(), o)
		sql = fmt.Sprintf("%s %s", sql, q)
	}

	rows, err := r.db.Query(c, sql, params)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	movies := make(map[int]models.Movie)
	for rows.Next() {
		var movie models.Movie
		var genre models.Genre
		err := rows.Scan(&movie.Id, &movie.Title, &movie.Description, &movie.DateOfRelease, &movie.Director,
			&movie.Rating, &movie.TrailerUrl, &movie.PosterUrl, &genre.Id, &genre.Title)
		if err != nil {
			return nil, err
		}

		_, exists := movies[movie.Id]
		if exists {
			movie = movies[movie.Id]
		}

		movie.Genres = append(movie.Genres, genre)
		movies[movie.Id] = movie
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	slice := make([]models.Movie, 0, len(movies))
	for _, m := range movies {
		slice = append(slice, m)
	}

	return slice, nil
}

func (r *MoviesRepository) FindById(c context.Context, id int) (models.Movie, error) {
	sql :=
		`
select m.id, 
       m.title, 
       m.description, 
       m.date_of_release, 
       m.director, 
       m.rating, 
       m.trailer_url, 
       m.poster_url,
       g.id,
       g.title
from movies m 
join movie_genres mg on mg.movie_id = m.id
join genres g on g.id = mg.genre_id
where m.id = $1
`

	rows, err := r.db.Query(c, sql, id)
	if err != nil {
		return models.Movie{}, err
	}
	defer rows.Close()

	movies := make(map[int]models.Movie)
	for rows.Next() {
		var movie models.Movie
		var genre models.Genre
		err := rows.Scan(&movie.Id, &movie.Title, &movie.Description, &movie.DateOfRelease, &movie.Director,
			&movie.Rating, &movie.TrailerUrl, &movie.PosterUrl, &genre.Id, &genre.Title)
		if err != nil {
			return models.Movie{}, err
		}

		_, exists := movies[movie.Id]
		if exists {
			movie = movies[movie.Id]
		}

		movie.Genres = append(movie.Genres, genre)
		movies[movie.Id] = movie
	}
	if err := rows.Err(); err != nil {
		return models.Movie{}, err
	}

	return movies[id], nil
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
