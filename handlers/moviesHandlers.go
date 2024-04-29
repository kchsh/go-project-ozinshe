package handlers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
	"ozinshe-final-project/models"
	"ozinshe-final-project/repositories"
	"path/filepath"
	"strconv"
	"time"
)

type MoviesHandler struct {
	moviesRepo *repositories.MoviesRepository
	genresRepo *repositories.GenresRepository
}

func NewMoviesHandler(moviesRepo *repositories.MoviesRepository, genresRepo *repositories.GenresRepository) *MoviesHandler {
	return &MoviesHandler{moviesRepo: moviesRepo, genresRepo: genresRepo}
}

func (h *MoviesHandler) HandleFindById(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, NewApiError("Invalid Movie Id"))
		return
	}

	movie, err := h.moviesRepo.FindById(c, id)
	if err != nil {
		c.JSON(http.StatusNotFound, NewApiError(err.Error()))
		return
	}
	c.JSON(http.StatusOK, movie)
}

func (h *MoviesHandler) HandleFindAll(c *gin.Context) {
	filters := models.MovieFilters{
		SearchTerm: c.Query("search"),
		IsWatched:  c.Query("iswatched"),
		GenreIds:   c.QueryArray("genreids"),
		Sort:       c.Query("sort"),
	}

	movies, err := h.moviesRepo.FindAll(c, filters)
	if err != nil {
		c.JSON(http.StatusInternalServerError, NewApiError(err.Error()))
		return
	}

	c.JSON(http.StatusOK, movies)
}

func (h *MoviesHandler) getGenresByIds(c *gin.Context, ids []int) ([]models.Genre, error) {
	genres, err := h.genresRepo.FindAll(c)
	if err != nil {
		return nil, err
	}

	selected := make([]models.Genre, 0)
	for _, genre := range genres {
		for _, id := range ids {
			if genre.Id == id {
				selected = append(selected, genre)
			}
		}
	}

	return selected, nil
}

func (h *MoviesHandler) getSavePathForImage(filename string) string {
	const imagesPath = "static/img/"
	return fmt.Sprintf("%s%s%s", imagesPath, uuid.New(), filepath.Ext(filename))
}

func (h *MoviesHandler) HandleCreate(c *gin.Context) {
	_, err := c.MultipartForm()
	if err != nil {
		c.JSON(http.StatusBadRequest, NewApiError("Invalid request payload"))
		return
	}

	title := c.PostForm("title")
	description := c.PostForm("description")
	date, _ := time.Parse("2006-01-02", c.PostForm("dateOfRelease"))
	director := c.PostForm("director")
	trailerUrl := c.PostForm("trailerUrl")

	genresArray := c.PostFormArray("genreIds")
	genreIds := make([]int, len(genresArray))
	for i, idStr := range genresArray {
		id, err := strconv.Atoi(idStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, NewApiError("Invalid genre id"))
			return
		}

		genreIds[i] = id
	}

	genres, err := h.getGenresByIds(c, genreIds)
	if err != nil {
		c.JSON(http.StatusInternalServerError, NewApiError(err.Error()))
		return
	}

	poster, _ := c.FormFile("poster")
	posterPath := h.getSavePathForImage(poster.Filename)
	err = c.SaveUploadedFile(poster, posterPath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, NewApiError(err.Error()))
	}

	movie := models.Movie{
		Title:         title,
		Description:   description,
		DateOfRelease: date,
		Director:      director,
		TrailerUrl:    trailerUrl,
		PosterUrl:     posterPath,
		Genres:        genres,
	}

	id, err := h.moviesRepo.Create(c, movie)
	if err != nil {
		c.JSON(http.StatusInternalServerError, NewApiError(err.Error()))
		return
	}
	c.JSON(http.StatusOK, gin.H{"id": id})
}

func (h *MoviesHandler) HandleUpdate(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, NewApiError(err.Error()))
		return
	}

	_, err = h.moviesRepo.FindById(c, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	title := c.PostForm("title")
	description := c.PostForm("description")
	date, _ := time.Parse("2006-01-02", c.PostForm("dateOfRelease"))
	director := c.PostForm("director")
	trailerUrl := c.PostForm("trailerUrl")

	genresArray := c.PostFormArray("genreIds")
	genreIds := make([]int, len(genresArray))
	for i, idStr := range genresArray {
		id, err := strconv.Atoi(idStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, NewApiError("Invalid genre id"))
			return
		}

		genreIds[i] = id
	}

	genres, err := h.getGenresByIds(c, genreIds)
	if err != nil {
		c.JSON(http.StatusInternalServerError, NewApiError(err.Error()))
		return
	}

	poster, _ := c.FormFile("poster")
	posterPath := h.getSavePathForImage(poster.Filename)
	err = c.SaveUploadedFile(poster, posterPath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, NewApiError(err.Error()))
	}

	movie := models.Movie{
		Id:            id,
		Title:         title,
		Description:   description,
		DateOfRelease: date,
		Director:      director,
		TrailerUrl:    trailerUrl,
		PosterUrl:     posterPath,
		Genres:        genres,
	}

	err = h.moviesRepo.Update(c, id, movie)
	if err != nil {
		c.JSON(http.StatusInternalServerError, NewApiError(err.Error()))
		return
	}
	c.Status(http.StatusOK)
}

func (h *MoviesHandler) HandleDelete(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, NewApiError(err.Error()))
		return
	}

	_, err = h.moviesRepo.FindById(c, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	err = h.moviesRepo.Delete(c, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, NewApiError(err.Error()))
		return
	}
	c.Status(http.StatusOK)
}

func (h *MoviesHandler) HandleSetRating(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, NewApiError("Invalid movie Id"))
		return
	}

	ratingStr := c.Query("rating")
	rating, err := strconv.Atoi(ratingStr)
	if err != nil || rating > 5 || rating < 1 {
		c.JSON(http.StatusBadRequest, NewApiError("Invalid rating value"))
		return
	}

	_, err = h.moviesRepo.FindById(c, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, NewApiError(err.Error()))
		return
	}

	err = h.moviesRepo.SetRating(c, id, rating)
	if err != nil {
		c.JSON(http.StatusInternalServerError, NewApiError(err.Error()))
		return
	}

	c.Status(http.StatusOK)
}

func (h *MoviesHandler) HandleSetWatched(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, NewApiError("Invalid movie Id"))
		return
	}

	isWatchedStr := c.Query("isWatched")
	isWatched, err := strconv.ParseBool(isWatchedStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, NewApiError("Invalid isWatched value"))
		return
	}

	_, err = h.moviesRepo.FindById(c, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, NewApiError(err.Error()))
		return
	}

	err = h.moviesRepo.SetWatched(c, id, isWatched)
	if err != nil {
		c.JSON(http.StatusInternalServerError, NewApiError(err.Error()))
		return
	}

	c.Status(http.StatusOK)
}
