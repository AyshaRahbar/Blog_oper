package routes

import (
	"go-blog/handlers"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(postHandler *handlers.PostHandler) *gin.Engine {
	router := gin.Default()
	api := router.Group("/api")
	{
		api.GET("/posts", postHandler.GetPosts)
		api.POST("/posts", postHandler.CreatePost)
		api.GET("/posts/:id", postHandler.GetPost)
		api.PUT("/posts/:id", postHandler.UpdatePost)
		api.DELETE("/posts/:id", postHandler.DeletePost)
	}
	return router
}
