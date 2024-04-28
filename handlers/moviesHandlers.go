package handlers

import (
	"github.com/gin-gonic/gin"
	"ozinshe-final-project/repositories"
	"strconv"
)

type MoviesHandler struct {
	repo *repositories.MoviesRepository
}

func NewMoviesHandler(repo *repositories.MoviesRepository) *MoviesHandler {
	return &MoviesHandler{repo: repo}
}

func (h *MoviesHandler) HandleFindById(ctx *gin.Context) {
	idParam := ctx.Query("id")
	id, _ := strconv.Atoi(idParam)

	movie := h.repo.FindById(id)

	ctx.JSON(200, movie)
}
