package routes

import (
	"go-blog/handlers"
	"go-blog/middleware"
	"go-blog/models"
	"go-blog/repo"
	"github.com/gin-gonic/gin"
)

func SetupRoutes(postHandler *handlers.PostHandler, userHandler *handlers.UserHandler, authRepo repo.AuthRepository) *gin.Engine {
	router := gin.Default()

	api := router.Group("/api")
	{
		api.POST("/register", userHandler.Register)
		api.POST("/login", userHandler.Login)
		api.GET("/posts", postHandler.GetPosts)
		api.GET("/posts/:id", postHandler.GetPostByID)

		bloggerType := models.AccountTypeBlogger
		api.POST("/posts", middleware.JWTAuthMiddleware(&bloggerType), postHandler.CreatePost)
		api.PUT("/posts/:id", middleware.JWTAuthMiddleware(nil), middleware.CheckPostOwnership(authRepo), postHandler.UpdatePost)
		api.DELETE("/posts/:id", middleware.JWTAuthMiddleware(nil), middleware.CheckPostOwnership(authRepo), postHandler.DeletePost)
	}
	return router
}
