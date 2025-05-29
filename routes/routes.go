package routes

import (
	"go-blog/handlers"
	"github.com/gin-gonic/gin"
)

func SetupRoutes(postHandler *handlers.PostHandler) *gin.Engine {
	router := gin.Default()
	api := router.Group("/api")
	{ //demo
		api.GET("/posts", postHandler.ListPosts)
		api.POST("/posts", postHandler.CreatePost)
		api.POST("/posts/:id/reaction", postHandler.AddReaction)
		api.DELETE("/posts/:id", postHandler.DeletePost)
	}
	router.GET("/posts", postHandler.ListPosts)
	router.POST("/posts", postHandler.CreatePost)
	router.DELETE("/posts/:id", postHandler.DeletePost)

	return router
}
