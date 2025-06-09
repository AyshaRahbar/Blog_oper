package middleware

import (
	"fmt"
	"go-blog/models"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"
)

func JWTAuthMiddleware(requiredType *models.AccountType) gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.GetHeader("Authorization")
		if header == "" || len(header) < 8 || header[:7] != "Bearer " {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "missing or invalid Authorization header"})
			c.Abort()
			return
		}

		tokenString := header[7:]
		secret := os.Getenv("JWT_SECRET")
		if secret == "" {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "JWT secret not configured"})
			c.Abort()
			return
		}
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(secret), nil
		})
		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			c.Abort()
			return
		}
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token claims"})
			c.Abort()
			return
		}
		userID, ok := claims["id"].(float64)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid user ID in token"})
			c.Abort()
			return
		}
		accountType, ok := claims["account_type"].(string)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid account type in token"})
			c.Abort()
			return
		}
		if requiredType != nil && models.AccountType(accountType) != *requiredType {
			c.JSON(http.StatusForbidden, gin.H{"error": "insufficient permissions"})
			c.Abort()
			return
		}
		c.Set("user_id", int(userID))
		c.Set("account_type", accountType)
		if username, exists := claims["username"]; exists {
			c.Set("username", username)
		}
		c.Next()
	}
}

func JWTAuth() gin.HandlerFunc {
	return JWTAuthMiddleware(nil)
}

func JWTAuthBlogger() gin.HandlerFunc {
	bloggerType := models.AccountTypeBlogger
	return JWTAuthMiddleware(&bloggerType)
}

func JWTAuthViewer() gin.HandlerFunc {
	viewerType := models.AccountTypeViewer
	return JWTAuthMiddleware(&viewerType)
}

func CheckPostOwnership(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, exists := c.Get("user_id")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
			c.Abort()
			return
		}

		postIDStr := c.Param("id")
		postID, err := strconv.Atoi(postIDStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid post ID"})
			c.Abort()
			return
		}

		var post models.Post
		if err := db.First(&post, "id = ?", postID).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				c.JSON(http.StatusNotFound, gin.H{"error": "post not found"})
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "database error"})
			}
			c.Abort()
			return
		}

		if post.UserID != userID.(int) {
			c.JSON(http.StatusForbidden, gin.H{"error": "you can only update your own posts"})
			c.Abort()
			return
		}
		c.Next()
	}
}
