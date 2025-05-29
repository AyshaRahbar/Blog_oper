package routes

import (
	"go-blog/handlers"
	"github.com/gin-gonic/gin"
)

func SetupRoutes(postHandler *handlers.PostHandler) *gin.Engine {
	router := gin.Default()

	router.GET("/posts", postHandler.GetAllPosts)
	router.POST("/posts", postHandler.CreatePost)
	router.GET("/posts/:id", postHandler.GetPostByID)
	router.PUT("/posts/:id", postHandler.UpdatePost)
	router.DELETE("/posts/:id", postHandler.DeletePost)
	router.POST("/posts/:id/reactions", postHandler.AddReaction)

	return router
}
