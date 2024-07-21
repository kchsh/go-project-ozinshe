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

// HandleFindAll godoc
// @Tags genres
// @Summary      Get genres list
// @Accept       json
// @Produce      json
// @Success      200  {array} models.Genre "OK"
// @Failure   	 400  {object} models.ApiError "Validation error"
// @Failure   	 500  {object} models.ApiError
// @Router       /genres [get]
// @Security Bearer
func (h *GenreHandlers) HandleFindAll(c *gin.Context) {
	genres, err := h.repo.FindAll(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.NewApiError(err.Error()))
		return
	}

	c.JSON(http.StatusOK, genres)
}

// HandleFindById godoc
// @Summary      Find genre by id
// @Tags genres
// @Accept       json
// @Produce      json
// @Param id path int true "Genre ID"
// @Success      200  {object} models.Genre "OK"
// @Failure   	 400  {object} models.ApiError "Validation error"
// @Failure   	 500  {object} models.ApiError
// @Router       /genres/{id} [get]
// @Security Bearer
func (h *GenreHandlers) HandleFindById(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.NewApiError("Invalid genre ID"))
		return
	}

	genre, err := h.repo.FindById(c, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.NewApiError(err.Error()))
		return
	}

	c.JSON(http.StatusOK, genre)
}

type createGenreRequest struct {
	Title string `json:"title"`
}

// HandleCreate godoc
// @Summary      Create genre
// @Tags genres
// @Accept       json
// @Produce      json
// @Param request body handlers.createGenreRequest true "Genre model"
// @Success      200  {object} object{id=int}  "OK"
// @Failure   	 400  {object} models.ApiError "Validation error"
// @Failure   	 500  {object} models.ApiError
// @Router       /genres [post]
// @Security Bearer
func (h *GenreHandlers) HandleCreate(c *gin.Context) {
	var request createGenreRequest
	if err := c.BindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, models.NewApiError("Invalid request payload"))
		return
	}

	genre := models.Genre{Title: request.Title}
	id, err := h.repo.Create(c, genre)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"id": id})
}

type updateGenreRequest struct {
	Title string `json:"title"`
}

// HandleUpdate godoc
// @Summary      Update genre
// @Tags genres
// @Accept       json
// @Produce      json
// @Param id path int true "Genre id"
// @Param request body handlers.updateGenreRequest true "Genre model"
// @Success      200
// @Failure   	 400  {object} models.ApiError "Validation error"
// @Failure   	 500  {object} models.ApiError
// @Router       /genres/{id} [put]
// @Security Bearer
func (h *GenreHandlers) HandleUpdate(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.NewApiError("Invalid genre id"))
		return
	}

	var request updateGenreRequest
	if err := c.BindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, models.NewApiError("Invalid request payload"))
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
		c.JSON(http.StatusInternalServerError, models.NewApiError(err.Error()))
		return
	}
	c.Status(http.StatusOK)
}

// HandleDelete godoc
// @Summary      Delete genre
// @Tags genres
// @Accept       json
// @Produce      json
// @Param id path int true "Genre id"
// @Success      200
// @Failure   	 400  {object} models.ApiError "Validation error"
// @Failure   	 500  {object} models.ApiError
// @Router       /genres/{id} [delete]
// @Security Bearer
func (h *GenreHandlers) HandleDelete(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.NewApiError("Invalid genre id"))
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
