package models

import "time"

type Movie struct {
	Id          int
	Title       string
	Description string
	ReleaseDate time.Time
	Director    string
	Rating      float32
	TrailerUrl  string
	PosterUtl   string
}
