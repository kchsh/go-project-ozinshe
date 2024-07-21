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

// HandleFindAll godoc
// @Tags users
// @Summary      Get users list
// @Accept       json
// @Produce      json
// @Success      200  {array} handlers.UserResponse "OK"
// @Failure   	 500  {object} models.ApiError
// @Router       /users [get]
// @Security Bearer
func (h *UserHandlers) HandleFindAll(c *gin.Context) {
	users, err := h.repo.FindAll(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.NewApiError(err.Error()))
		return
	}

	r := MapUsersToResponse(users)

	c.JSON(http.StatusOK, r)
}

// HandleFindById godoc
// @Tags users
// @Summary      Find users by id
// @Accept       json
// @Produce      json
// @Param id path int true "User id"
// @Success      200  {array} handlers.UserResponse "OK"
// @Failure   	 400  {object} models.ApiError "Invalid user id"
// @Failure   	 404  {object} models.ApiError "User not found"
// @Failure   	 500  {object} models.ApiError
// @Router       /users/{id} [get]
// @Security Bearer
func (h *UserHandlers) HandleFindById(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.NewApiError("Invalid user Id"))
		return
	}

	user, err := h.repo.FindById(c, id)
	if err != nil {
		c.JSON(http.StatusNotFound, models.NewApiError("User not found"))
		return
	}

	r := MapUserToResponse(user)

	c.JSON(http.StatusOK, r)
}

// HandleCreate godoc
// @Tags users
// @Summary      Create user
// @Accept       json
// @Produce      json
// @Param request body handlers.createUserRequest true "User data"
// @Success      200  {object} object{id=int} "OK"
// @Failure   	 400  {object} models.ApiError "Invalid data"
// @Failure   	 500  {object} models.ApiError
// @Router       /users [post]
// @Security Bearer
func (h *UserHandlers) HandleCreate(c *gin.Context) {
	var request createUserRequest
	if err := c.BindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, models.NewApiError("Invalid request payload"))
		return
	}

	if request.Password != request.ConfirmPassword {
		c.JSON(http.StatusBadRequest, models.NewApiError("Passwords miss match"))
		return
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.NewApiError("Failed to hash password"))
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

// HandleUpdate godoc
// @Tags users
// @Summary      Update user
// @Accept       json
// @Produce      json
// @Param id path int true "User id"
// @Param request body handlers.updateUserRequest true "User data"
// @Success      200  {object} object{id=int} "OK"
// @Failure   	 400  {object} models.ApiError "Invalid data"
// @Failure   	 404  {object} models.ApiError "User not found"
// @Failure   	 500  {object} models.ApiError
// @Router       /users/{id} [put]
// @Security Bearer
func (h *UserHandlers) HandleUpdate(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.NewApiError("Invalid user Id"))
		return
	}

	var request updateUserRequest
	if err := c.BindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, models.NewApiError("Invalid request payload"))
		return
	}

	user, err := h.repo.FindById(c, id)
	if err != nil {
		c.JSON(http.StatusNotFound, models.NewApiError("User not found"))
		return
	}

	user.Name = request.Name
	user.Email = request.Email

	err = h.repo.Update(c, id, user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.NewApiError(err.Error()))
		return
	}

	c.Status(http.StatusOK)
}

// HandleDelete godoc
// @Tags users
// @Summary      Delete user
// @Accept       json
// @Produce      json
// @Param id path int true "User id"
// @Success      200  "OK"
// @Failure   	 400  {object} models.ApiError "Invalid data"
// @Failure   	 404  {object} models.ApiError "User not found"
// @Failure   	 500  {object} models.ApiError
// @Router       /users/{id} [delete]
// @Security Bearer
func (h *UserHandlers) HandleDelete(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.NewApiError("Invalid user Id"))
		return
	}

	_, err = h.repo.FindById(c, id)
	if err != nil {
		c.JSON(http.StatusNotFound, models.NewApiError("User not found"))
		return
	}

	err = h.repo.Delete(c, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.NewApiError(err.Error()))
		return
	}

	c.Status(http.StatusOK)
}

// HandleChangePassword godoc
// @Tags users
// @Summary      Change user password
// @Accept       json
// @Produce      json
// @Param id path int true "User id"
// @Param request body handlers.changePasswordRequest true "Password data"
// @Success      200  "OK"
// @Failure   	 400  {object} models.ApiError "Invalid data"
// @Failure   	 404  {object} models.ApiError "User not found"
// @Failure   	 500  {object} models.ApiError
// @Router       /users/{id}/changePassword [put]
// @Security Bearer
func (h *UserHandlers) HandleChangePassword(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.NewApiError("Invalid user Id"))
		return
	}

	var request changePasswordRequest
	if err := c.BindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, models.NewApiError("Invalid request payload"))
		return
	}

	if request.Password != request.ConfirmPassword {
		c.JSON(http.StatusBadRequest, models.NewApiError("Passwords miss match"))
		return
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.NewApiError("Failed to hash password"))
		return
	}

	user, err := h.repo.FindById(c, id)
	if err != nil {
		c.JSON(http.StatusNotFound, models.NewApiError("User not found"))
		return
	}

	user.PasswordHash = string(passwordHash)

	err = h.repo.Update(c, id, user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.NewApiError(err.Error()))
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
