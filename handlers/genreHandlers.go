package handlers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"ozinshe-final-project/models"
	"ozinshe-final-project/repositories"
	"strconv"
)

type GenreHandlers struct {
	repo *repositories.GenresRepository
}

func NewGenreHandlers(repo *repositories.GenresRepository) *GenreHandlers {
	return &GenreHandlers{repo: repo}
}

func (h *GenreHandlers) HandleFindAll(c *gin.Context) {
	genres, err := h.repo.FindAll(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, NewApiError(err.Error()))
		return
	}

	c.JSON(http.StatusOK, genres)
}

func (h *GenreHandlers) HandleFindById(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, NewApiError("Invalid genre ID"))
		return
	}

	genre, err := h.repo.FindById(c, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, NewApiError(err.Error()))
		return
	}

	c.JSON(http.StatusOK, genre)
}

type createGenreRequest struct {
	Title string `json:"title"`
}

func (h *GenreHandlers) HandleCreate(c *gin.Context) {
	var request createGenreRequest
	if err := c.BindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, NewApiError("Invalid request payload"))
		return
	}

	genre := models.Genre{Title: request.Title}
	id, err := h.repo.Create(c, genre)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"id": id})
}

type updateGenreRequest struct {
	Title string `json:"title"`
}

func (h *GenreHandlers) HandleUpdate(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, NewApiError("Invalid genre id"))
		return
	}

	var request updateGenreRequest
	if err := c.BindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, NewApiError("Invalid request payload"))
		return
	}

	_, err = h.repo.FindById(c, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	genre := models.Genre{Title: request.Title}
	err = h.repo.Update(c, id, genre)
	if err != nil {
		c.JSON(http.StatusInternalServerError, NewApiError(err.Error()))
		return
	}
	c.Status(http.StatusOK)
}

func (h *GenreHandlers) HandleDelete(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, NewApiError("Invalid genre id"))
		return
	}

	_, err = h.repo.FindById(c, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	err = h.repo.Delete(c, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusOK)
}
