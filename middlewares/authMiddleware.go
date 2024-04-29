package middlewares

import (
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
	"ozinshe-final-project/handlers"
	"strconv"
	"strings"
)

func AuthMiddleware(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		c.JSON(http.StatusUnauthorized, handlers.NewApiError("Authorization header required"))
		c.Abort()
		return
	}

	tokenString := strings.Split(authHeader, "Bearer ")[1]
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return handlers.JWTSecretKey, nil
	})
	if err != nil || !token.Valid {
		c.JSON(http.StatusUnauthorized, handlers.NewApiError("Invalid token"))
		c.Abort()
		return
	}

	subject, err := token.Claims.GetSubject()
	if err != nil {
		c.JSON(http.StatusUnauthorized, handlers.NewApiError("Error while getting subject"))
		c.Abort()
		return
	}

	userId, _ := strconv.Atoi(subject)
	c.Set("userId", userId)
	c.Next()
}
