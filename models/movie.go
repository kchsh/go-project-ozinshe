package models

import "time"

type MovieFilters struct {
	SearchTerm string
	GenreIds   []string
	IsWatched  string
	Sort       string
}

type Movie struct {
	Id            int
	Title         string
	Description   string
	DateOfRelease time.Time
	Director      string
	Rating        int
	TrailerUrl    string
	PosterUrl     string
	IsWatched     bool
	Genres        []Genre
}
