package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"ozinshe-final-project/config"
	"ozinshe-final-project/repositories"
	"strconv"
	"time"
)

type AuthHandlers struct {
	usersRepo *repositories.UsersRepository
}

func NewAuthHandlers(usersRepo *repositories.UsersRepository) *AuthHandlers {
	return &AuthHandlers{usersRepo: usersRepo}
}

type signInRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (h *AuthHandlers) HandleSignIn(c *gin.Context) {
	var request signInRequest
	if err := c.BindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, NewApiError("Invalid request payload"))
		return
	}

	user, err := h.usersRepo.FindByEmail(c, request.Email)
	if err != nil {
		c.JSON(http.StatusUnauthorized, NewApiError("Invalid credentials"))
		return
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(request.Password))
	if err != nil {
		c.JSON(http.StatusUnauthorized, NewApiError("Invalid credentials"))
		return
	}

	claims := jwt.RegisteredClaims{
		Subject:   strconv.Itoa(user.Id),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(config.Config.JwtExpiresIn)),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(config.Config.JwtSecretKey))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": tokenString})
}

func (h *AuthHandlers) HandleGetUserInfo(c *gin.Context) {
	userId := c.GetInt("userId")
	user, err := h.usersRepo.FindById(c, userId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	r := MapUserToResponse(user)

	c.JSON(http.StatusOK, r)
}

func (h *AuthHandlers) HandleSignOut(c *gin.Context) {
	c.Status(http.StatusOK)
}
