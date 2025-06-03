package routes

import (
	"go-blog/handlers"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(postHandler *handlers.PostHandler) *gin.Engine {
	router := gin.Default()
	api := router.Group("/api")
	{ //demo
		api.POST("/posts", postHandler.CreatePost)
		api.DELETE("/posts/:id", postHandler.DeletePost)
		api.GET("/posts", postHandler.GetPosts)
		api.PUT("/posts/:id/update", postHandler.UpdatePost)
		api.GET("/posts/:id",postHandler.GetPostByID)
	}

	return router
}
