package handlers

import (
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"ozinshe-final-project/models"
	"ozinshe-final-project/repositories"
	"strconv"
)

type UserHandlers struct {
	repo *repositories.UsersRepository
}

func NewUserHandlers(repo *repositories.UsersRepository) *UserHandlers {
	return &UserHandlers{repo: repo}
}

type createUserRequest struct {
	Name            string `json:"name"`
	Email           string `json:"email"`
	Password        string `json:"password"`
	ConfirmPassword string `json:"confirmPassword"`
}

type updateUserRequest struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

type changePasswordRequest struct {
	Password        string `json:"password"`
	ConfirmPassword string `json:"confirmPassword"`
}

type UserResponse struct {
	Id    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

func (h *UserHandlers) HandleFindAll(c *gin.Context) {
	users, err := h.repo.FindAll(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, NewApiError(err.Error()))
		return
	}

	r := MapUsersToResponse(users)

	c.JSON(http.StatusOK, r)
}

func (h *UserHandlers) HandleFindById(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, NewApiError("Invalid user Id"))
		return
	}

	user, err := h.repo.FindById(c, id)
	if err != nil {
		c.JSON(http.StatusNotFound, NewApiError("User not found"))
		return
	}

	r := MapUserToResponse(user)

	c.JSON(http.StatusOK, r)
}

func (h *UserHandlers) HandleCreate(c *gin.Context) {
	var request createUserRequest
	if err := c.BindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, NewApiError("Invalid request payload"))
		return
	}

	if request.Password != request.ConfirmPassword {
		c.JSON(http.StatusBadRequest, NewApiError("Passwords miss match"))
		return
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, NewApiError("Failed to hash password"))
		return
	}

	user := models.User{
		Name:         request.Name,
		Email:        request.Email,
		PasswordHash: string(passwordHash),
	}

	id, err := h.repo.Create(c, user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"id": id})
}

func (h *UserHandlers) HandleUpdate(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, NewApiError("Invalid user Id"))
		return
	}

	var request updateUserRequest
	if err := c.BindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, NewApiError("Invalid request payload"))
		return
	}

	user, err := h.repo.FindById(c, id)
	if err != nil {
		c.JSON(http.StatusNotFound, NewApiError("User not found"))
		return
	}

	user.Name = request.Name
	user.Email = request.Email

	err = h.repo.Update(c, id, user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, NewApiError(err.Error()))
		return
	}

	c.Status(http.StatusOK)
}

func (h *UserHandlers) HandleDelete(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, NewApiError("Invalid user Id"))
		return
	}

	_, err = h.repo.FindById(c, id)
	if err != nil {
		c.JSON(http.StatusNotFound, NewApiError("User not found"))
		return
	}

	err = h.repo.Delete(c, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, NewApiError(err.Error()))
		return
	}

	c.Status(http.StatusOK)
}

func (h *UserHandlers) HandleChangePassword(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, NewApiError("Invalid user Id"))
		return
	}

	var request changePasswordRequest
	if err := c.BindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, NewApiError("Invalid request payload"))
		return
	}

	if request.Password != request.ConfirmPassword {
		c.JSON(http.StatusBadRequest, NewApiError("Passwords miss match"))
		return
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, NewApiError("Failed to hash password"))
		return
	}

	user, err := h.repo.FindById(c, id)
	if err != nil {
		c.JSON(http.StatusNotFound, NewApiError("User not found"))
		return
	}

	user.PasswordHash = string(passwordHash)

	err = h.repo.Update(c, id, user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, NewApiError(err.Error()))
		return
	}

	c.Status(http.StatusOK)
}

func MapUsersToResponse(users []models.User) []UserResponse {
	usersResponse := make([]UserResponse, 0, len(users))

	for _, user := range users {
		r := UserResponse{
			Id:    user.Id,
			Name:  user.Name,
			Email: user.Email,
		}

		usersResponse = append(usersResponse, r)
	}

	return usersResponse
}

func MapUserToResponse(user models.User) UserResponse {
	return UserResponse{
		Id:    user.Id,
		Name:  user.Name,
		Email: user.Email,
	}
}
