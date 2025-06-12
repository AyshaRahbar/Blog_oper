package routes

import (
	"go-blog/handlers"
	"go-blog/middleware"
	"go-blog/models"
	"go-blog/repo"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(postHandler *handlers.PostHandler, userHandler *handlers.UserHandler, likeHandler *handlers.LikeHandler, commentHandler *handlers.CommentHandler, authRepo repo.AuthRepository) *gin.Engine {
	router := gin.Default()

	api := router.Group("/api")
	{
		api.POST("/register", userHandler.Register)
		api.POST("/login", userHandler.Login)
		api.GET("/users", userHandler.GetUsers)
		api.GET("/posts", postHandler.GetPosts)
		api.GET("/posts/:id", postHandler.GetPostByID)

		api.GET("/posts/:id/likes", likeHandler.GetPostLikes)
		api.GET("/posts/:id/likes/detail", likeHandler.GetPostLikesCount)
		api.GET("/posts/:id/comments", commentHandler.GetPostComments)
		api.GET("/posts/:id/with-comments", commentHandler.GetPostWithComments)

		bloggerType := models.AccountTypeBlogger

		api.POST("/posts", middleware.JWTAuthMiddleware(&bloggerType), postHandler.CreatePost)
		api.PUT("/posts/:id", middleware.JWTAuthMiddleware(nil), middleware.CheckPostOwnership(authRepo), postHandler.UpdatePost)
		api.DELETE("/posts/:id", middleware.JWTAuthMiddleware(nil), middleware.CheckPostOwnership(authRepo), postHandler.DeletePost)

		api.POST("/posts/:id/like", middleware.JWTAuthMiddleware(nil), likeHandler.LikePost)
		api.DELETE("/posts/:id/like", middleware.JWTAuthMiddleware(nil), likeHandler.UnlikePost)
		api.GET("/users/me/likes", middleware.JWTAuthMiddleware(nil), likeHandler.GetUserLikes)

		api.POST("/posts/:id/comments", middleware.JWTAuthMiddleware(nil), commentHandler.CreateComment)
		api.PUT("/comments/:id", middleware.JWTAuthMiddleware(nil), commentHandler.UpdateComment)
		api.DELETE("/comments/:id", middleware.JWTAuthMiddleware(nil), commentHandler.DeleteComment)
		api.GET("/users/me/comments", middleware.JWTAuthMiddleware(nil), commentHandler.GetUserComments)
	}
	return router
}
