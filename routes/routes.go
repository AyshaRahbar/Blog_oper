package routes

import (
	"go-blog/handlers"
	"go-blog/middleware"
	"go-blog/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupRoutes(postHandler *handlers.PostHandler, userHandler *handlers.UserHandler, db *gorm.DB) *gin.Engine {
	router := gin.Default()

	api := router.Group("/api")
	{
		api.POST("/register", userHandler.Register)
		api.POST("/login", userHandler.Login)
		api.GET("/posts", postHandler.GetPosts)
		api.GET("/posts/:id", postHandler.GetPostByID)
		api.POST("/posts",middleware.JWTAuthMiddleware(),middleware.RequireAccountType(models.AccountTypeBlogger),postHandler.CreatePost)
		api.PUT("/posts/:id", middleware.JWTAuthMiddleware(),middleware.CheckPostOwnership(db),postHandler.UpdatePost)
		api.DELETE("/posts/:id",middleware.JWTAuthMiddleware(),	middleware.CheckPostOwnership(db),postHandler.DeletePost)
	}
	return router
}
