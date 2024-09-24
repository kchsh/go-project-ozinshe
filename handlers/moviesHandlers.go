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
)

type MoviesHandler struct {
	moviesRepo *repositories.MoviesRepository
	genresRepo *repositories.GenresRepository
}

func NewMoviesHandler(moviesRepo *repositories.MoviesRepository, genresRepo *repositories.GenresRepository) *MoviesHandler {
	return &MoviesHandler{moviesRepo: moviesRepo, genresRepo: genresRepo}
}

// HandleFindById godoc
// @Summary      Find by id
// @Tags movies
// @Accept       json
// @Produce      json
// @Param id path int true "Movie id"
// @Success      200  {object} models.Movie "OK"
// @Failure   	 400  {object} models.ApiError "Invalid movie id"
// @Failure   	 404  {object} models.ApiError "Movie not found"
// @Failure   	 500  {object} models.ApiError
// @Router       /movies/{id} [get]
// @Security Bearer
func (h *MoviesHandler) HandleFindById(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.NewApiError("Invalid Movie Id"))
		return
	}

	movie, err := h.moviesRepo.FindById(c, id)
	if err != nil {
		c.JSON(http.StatusNotFound, models.NewApiError(err.Error()))
		return
	}
	c.JSON(http.StatusOK, movie)
}

// HandleFindAll godoc
// @Summary      Get all movies
// @Tags movies
// @Accept       json
// @Produce      json
// @Param filters query models.MovieFilters true "Movie filters"
// @Success      200  {object} models.Movie "OK"
// @Failure   	 500  {object} models.ApiError
// @Router       /movies [get]
// @Security Bearer
func (h *MoviesHandler) HandleFindAll(c *gin.Context) {
	filters := models.MovieFilters{
		SearchTerm: c.Query("search"),
		IsWatched:  c.Query("iswatched"),
		GenreIds:   c.QueryArray("genreids"),
		Sort:       c.Query("sort"),
	}

	movies, err := h.moviesRepo.FindAll(c, filters)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.NewApiError(err.Error()))
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

func (h *MoviesHandler) transformFilename(poster string) string {
	return fmt.Sprintf("%s%s", uuid.New(), filepath.Ext(poster))

}

func (h *MoviesHandler) getFilePath(filename string) string {
	return fmt.Sprintf("images/%s", filename)
}

// HandleCreate godoc
// @Summary      Create movie
// @Tags movies
// @Accept       multipart/form-data
// @Produce      json
// @Param title formData string true "Title"
// @Param description formData string true "Description"
// @Param releaseYear formData int true "Year of release"
// @Param director formData string true "Director"
// @Param trailerUrl formData string true "Trailer URL"
// @Param genreIds formData []int true "Genre ids"
// @Param poster formData file true "Poster image"
// @Success      200  {object} object{id=int} "OK"
// @Failure   	 400  {object} models.ApiError "Invalid data"
// @Failure   	 500  {object} models.ApiError
// @Router       /movies [post]
// @Security Bearer
func (h *MoviesHandler) HandleCreate(c *gin.Context) {
	_, err := c.MultipartForm()
	if err != nil {
		c.JSON(http.StatusBadRequest, models.NewApiError("Invalid request payload"))
		return
	}

	title := c.PostForm("title")
	description := c.PostForm("description")
	releaseYearStr := c.PostForm("releaseYear")
	releaseYear, err := strconv.Atoi(releaseYearStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.NewApiError(err.Error()))
	}
	director := c.PostForm("director")
	trailerUrl := c.PostForm("trailerUrl")

	genresArray := c.PostFormArray("genreIds")
	genreIds := make([]int, len(genresArray))
	for i, idStr := range genresArray {
		id, err := strconv.Atoi(idStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, models.NewApiError("Invalid genre id"))
			return
		}

		genreIds[i] = id
	}

	genres, err := h.getGenresByIds(c, genreIds)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.NewApiError(err.Error()))
		return
	}

	poster, _ := c.FormFile("poster")
	filename := h.transformFilename(poster.Filename)
	filePath := h.getFilePath(filename)
	err = c.SaveUploadedFile(poster, filePath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.NewApiError(err.Error()))
	}

	movie := models.Movie{
		Title:       title,
		Description: description,
		ReleaseYear: releaseYear,
		Director:    director,
		TrailerUrl:  trailerUrl,
		PosterUrl:   filename,
		Genres:      genres,
	}

	id, err := h.moviesRepo.Create(c, movie)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.NewApiError(err.Error()))
		return
	}
	c.JSON(http.StatusOK, gin.H{"id": id})
}

// HandleUpdate godoc
// @Summary      Update movie
// @Tags movies
// @Accept       multipart/form-data
// @Produce      json
// @Param id path int true "Movie id"
// @Param title formData string true "Title"
// @Param description formData string true "Description"
// @Param releaseYear formData int true "Year of release"
// @Param director formData string true "Director"
// @Param trailerUrl formData string true "Trailer URL"
// @Param genreIds formData []int true "Genre ids"
// @Param poster formData file true "Poster image"
// @Success      200  {object} object{id=int} "OK"
// @Failure   	 400  {object} models.ApiError "Invalid data"
// @Failure   	 500  {object} models.ApiError
// @Router       /movies/{id} [put]
// @Security Bearer
func (h *MoviesHandler) HandleUpdate(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.NewApiError(err.Error()))
		return
	}

	_, err = h.moviesRepo.FindById(c, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	title := c.PostForm("title")
	description := c.PostForm("description")
	releaseYearStr := c.PostForm("releaseYear")
	releaseYear, err := strconv.Atoi(releaseYearStr)
	director := c.PostForm("director")
	trailerUrl := c.PostForm("trailerUrl")

	genresArray := c.PostFormArray("genreIds")
	genreIds := make([]int, len(genresArray))
	for i, idStr := range genresArray {
		id, err := strconv.Atoi(idStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, models.NewApiError("Invalid genre id"))
			return
		}

		genreIds[i] = id
	}

	genres, err := h.getGenresByIds(c, genreIds)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.NewApiError(err.Error()))
		return
	}

	poster, _ := c.FormFile("poster")
	filename := h.transformFilename(poster.Filename)
	filePath := h.getFilePath(filename)
	err = c.SaveUploadedFile(poster, filePath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.NewApiError(err.Error()))
	}

	movie := models.Movie{
		Id:          id,
		Title:       title,
		Description: description,
		ReleaseYear: releaseYear,
		Director:    director,
		TrailerUrl:  trailerUrl,
		PosterUrl:   filename,
		Genres:      genres,
	}

	err = h.moviesRepo.Update(c, id, movie)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.NewApiError(err.Error()))
		return
	}
	c.Status(http.StatusOK)
}

// HandleDelete godoc
// @Summary      Delete movie
// @Tags movies
// @Accept       json
// @Produce      json
// @Param id path int true "Movie id"
// @Success      200  "OK"
// @Failure   	 400  {object} models.ApiError "Invalid data"
// @Failure   	 500  {object} models.ApiError
// @Router       /movies/{id} [delete]
// @Security Bearer
func (h *MoviesHandler) HandleDelete(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.NewApiError(err.Error()))
		return
	}

	_, err = h.moviesRepo.FindById(c, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	err = h.moviesRepo.Delete(c, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.NewApiError(err.Error()))
		return
	}
	c.Status(http.StatusOK)
}

// HandleSetRating godoc
// @Summary      Set movie rating
// @Tags movies
// @Accept       json
// @Produce      json
// @Param id path int true "Movie id"
// @Param rating query int true "Movie rating"
// @Success      200  "OK"
// @Failure   	 400  {object} models.ApiError "Invalid data"
// @Failure   	 500  {object} models.ApiError
// @Router       /movies/{id}/rate [patch]
// @Security Bearer
func (h *MoviesHandler) HandleSetRating(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.NewApiError("Invalid movie Id"))
		return
	}

	ratingStr := c.Query("rating")
	rating, err := strconv.Atoi(ratingStr)
	if err != nil || rating > 5 || rating < 1 {
		c.JSON(http.StatusBadRequest, models.NewApiError("Invalid rating value"))
		return
	}

	_, err = h.moviesRepo.FindById(c, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.NewApiError(err.Error()))
		return
	}

	err = h.moviesRepo.SetRating(c, id, rating)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.NewApiError(err.Error()))
		return
	}

	c.Status(http.StatusOK)
}

// HandleSetWatched godoc
// @Summary      Mark movie as watched
// @Tags movies
// @Accept       json
// @Produce      json
// @Param id path int true "Movie id"
// @Param isWatched query bool true "Flag value"
// @Success      200  "OK"
// @Failure   	 400  {object} models.ApiError "Invalid data"
// @Failure   	 500  {object} models.ApiError
// @Router       /movies/{id}/setWatched [patch]
// @Security Bearer
func (h *MoviesHandler) HandleSetWatched(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.NewApiError("Invalid movie Id"))
		return
	}

	isWatchedStr := c.Query("isWatched")
	isWatched, err := strconv.ParseBool(isWatchedStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.NewApiError("Invalid isWatched value"))
		return
	}

	_, err = h.moviesRepo.FindById(c, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.NewApiError(err.Error()))
		return
	}

	err = h.moviesRepo.SetWatched(c, id, isWatched)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.NewApiError(err.Error()))
		return
	}

	c.Status(http.StatusOK)
}
