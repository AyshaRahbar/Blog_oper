package handlers

import (
	"go-blog/models"
	"go-blog/service"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type AuthInput struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type UserHandler struct {
	service *service.UserService
}

func NewUserHandler(service *service.UserService) *UserHandler {
	return &UserHandler{service: service}
}

func (h *UserHandler) Login(c *gin.Context) {
	var authInput AuthInput
	if err := c.ShouldBindJSON(&authInput); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := h.service.Login(authInput.Username, authInput.Password)
	if err != nil {
		if err.Error() == "invalid password" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid password"})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": "user not found"})
		return
	}
	generateToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":           user.ID,
		"username":     user.Username,
		"account_type": user.AccountType,
		"exp":          time.Now().Add(time.Hour * 24).Unix(),
	})

	token, err := generateToken.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to generate token"})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"token": token,
	})
}

func (h *UserHandler) Register(c *gin.Context) {
	var req models.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}
	user, err := h.service.Register(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{
		"message":      "user created",
		"user":         user.Username,
		"account_type": user.AccountType,
	})
}
