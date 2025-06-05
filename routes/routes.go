package routes

import (
	"go-blog/handlers"
	"go-blog/models"
	"net/http"
	"os"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func SetupRoutes(postHandler *handlers.PostHandler, userHandler *handlers.UserHandler) *gin.Engine {
	router := gin.Default()
	api := router.Group("/api")
	{
		api.POST("/login", userHandler.Login)
		api.POST("/register", userHandler.Register)
		api.POST("/posts", JWTAuthMiddleware(models.AccountTypeBlogger), postHandler.CreatePost)
		api.DELETE("/posts/:id", JWTAuthMiddleware(models.AccountTypeBlogger), postHandler.DeletePost)
		api.PUT("/posts/:id/update", JWTAuthMiddleware(models.AccountTypeBlogger), postHandler.UpdatePost)
		api.GET("/posts", JWTAuthMiddleware(models.AccountType("")), postHandler.GetPosts)
		api.GET("/posts/:id", JWTAuthMiddleware(models.AccountType("")), postHandler.GetPostByID)
	}
	return router
}

func JWTAuthMiddleware(requiredType models.AccountType) gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.GetHeader("Authorization")
		if header == "" || len(header) < 8 || header[:7] != "Bearer " {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "missing or invalid Authorization header"})
			c.Abort()
			return
		}
		tokenString := header[7:]
		secret := os.Getenv("JWT_SECRET")
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
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
		accountType, ok := claims["account_type"].(string)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid account type in token"})
			c.Abort()
			return
		}
		if requiredType != "" && models.AccountType(accountType) != requiredType {
			c.JSON(http.StatusForbidden, gin.H{"error": "you dont permission to execute this operation"})
			c.Abort()
			return
		}
		c.Set("username", claims["username"])
		c.Set("account_type", accountType)
		c.Next()
	}
}
