package models

type MovieFilters struct {
	SearchTerm string
	GenreIds   []string
	IsWatched  string
	Sort       string
}

type Movie struct {
	Id          int
	Title       string
	Description string
	ReleaseYear int
	Director    string
	Rating      int
	TrailerUrl  string
	PosterUrl   string
	IsWatched   bool
	Genres      []Genre
}
