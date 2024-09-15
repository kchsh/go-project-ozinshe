package handlers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"ozinshe-final-project/models"
	"ozinshe-final-project/repositories"
	"strconv"
)

type WatchlistHandler struct {
	moviesRepo    *repositories.MoviesRepository
	watchlistRepo *repositories.WatchlistRepository
}

func NewWatchlistHandler(moviesRepo *repositories.MoviesRepository, watchlistRepo *repositories.WatchlistRepository) *WatchlistHandler {
	return &WatchlistHandler{moviesRepo: moviesRepo, watchlistRepo: watchlistRepo}
}

// HandleGetMovies godoc
// @Summary      Get movies watchlist
// @Tags watchlist
// @Accept       json
// @Produce      json
// @Success      200 {array} models.Movie  "OK"
// @Failure   	 500  {object} models.ApiError
// @Router       /watchlist [get]
// @Security Bearer
func (h *WatchlistHandler) HandleGetMovies(c *gin.Context) {
	movies, err := h.watchlistRepo.GetMoviesFromWatchlist(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.NewApiError(err.Error()))
		return
	}

	c.JSON(http.StatusOK, movies)
}

// HandleAddMovie godoc
// @Summary      Add movie to watchlist
// @Tags watchlist
// @Accept       json
// @Produce      json
// @Param movieId path int true "Movie id"
// @Success      200 "OK"
// @Failure   	 400  {object} models.ApiError "Invalid data"
// @Failure   	 500  {object} models.ApiError
// @Router       /watchlist/:movieId [post]
// @Security Bearer
func (h *WatchlistHandler) HandleAddMovie(c *gin.Context) {
	idStr := c.Param("movieId")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.NewApiError("Invalid movie id"))
		return
	}
	_, err = h.moviesRepo.FindById(c, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.NewApiError(err.Error()))
		return
	}

	err = h.watchlistRepo.AddToWatchlist(c, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.NewApiError(err.Error()))
		return
	}

	c.Status(http.StatusOK)
}

// HandleRemoveMovie godoc
// @Summary      Remove movie from watchlist
// @Tags watchlist
// @Accept       json
// @Produce      json
// @Param movieId path int true "Movie id"
// @Success      200 "OK"
// @Failure   	 400  {object} models.ApiError "Invalid data"
// @Failure   	 500  {object} models.ApiError
// @Router       /watchlist/:movieId [delete]
// @Security Bearer
func (h *WatchlistHandler) HandleRemoveMovie(c *gin.Context) {
	idStr := c.Param("movieId")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.NewApiError("Invalid movie id"))
		return
	}

	_, err = h.moviesRepo.FindById(c, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.NewApiError(err.Error()))
		return
	}

	err = h.watchlistRepo.RemoveFromWatchlist(c, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.NewApiError(err.Error()))
		return
	}

	c.Status(http.StatusOK)
}
