package models

import "time"

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
