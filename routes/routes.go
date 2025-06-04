package routes

import (
	"go-blog/handlers"
	"go-blog/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(postHandler *handlers.PostHandler, userHandler *handlers.UserHandler) *gin.Engine {
	router := gin.Default()
	api := router.Group("/api")
	{
		api.POST("/login", userHandler.Login)
		api.POST("/register", userHandler.Register)
		api.POST("/posts", PermissionMiddleware(models.AccountTypeBlogger), postHandler.CreatePost)
		api.DELETE("/posts/:id", PermissionMiddleware(models.AccountTypeBlogger), postHandler.DeletePost)
		api.PUT("/posts/:id/update", PermissionMiddleware(models.AccountTypeBlogger), postHandler.UpdatePost)
		api.GET("/posts", PermissionMiddleware(models.AccountType("")), postHandler.GetPosts)
		api.GET("/posts/:id", PermissionMiddleware(models.AccountType("")), postHandler.GetPostByID)
	}
	return router
}

func PermissionMiddleware(requiredType models.AccountType) gin.HandlerFunc {
	return func(c *gin.Context) {
		accountType := models.AccountType(c.GetHeader("Account-Type"))
		if requiredType != "" && accountType != requiredType {
			c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
			c.Abort()
			return
		}
		c.Next()
	}
}
