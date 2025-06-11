package handlers

import (
	"go-blog/models"
	"go-blog/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type LikeHandler struct {
	likeService service.LikeService
}

func NewLikeHandler(likeService service.LikeService) *LikeHandler {
	return &LikeHandler{likeService: likeService}
}

func (h *LikeHandler) LikePost(c *gin.Context) {
	postIDStr := c.Param("id")
	postID, err := strconv.Atoi(postIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid post ID"})
		return
	}
	userIDInterface, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}
	userID, ok := userIDInterface.(int)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user ID format"})
		return
	}
	likeResponse, err := h.likeService.LikePost(userID, postID)
	if err != nil {
		switch err {
		case models.ErrPostNotFoundForLike:
			c.JSON(http.StatusNotFound, gin.H{"error": "Post not found"})
		case models.ErrCannotLikeOwnPost:
			c.JSON(http.StatusForbidden, gin.H{"error": "Cannot like your own post"})
		case models.ErrLikeAlreadyExists:
			c.JSON(http.StatusConflict, gin.H{"error": "You have already liked this post"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to like post"})
		}
		return
	}
	c.JSON(http.StatusCreated, gin.H{
		"message": "Post liked successfully",
		"like":    likeResponse,
	})
}

func (h *LikeHandler) UnlikePost(c *gin.Context) {
	postIDStr := c.Param("id")
	postID, err := strconv.Atoi(postIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid post ID"})
		return
	}
	userIDInterface, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}
	userID, ok := userIDInterface.(int)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user ID format"})
		return
	}
	err = h.likeService.UnlikePost(userID, postID)
	if err != nil {
		switch err {
		case models.ErrPostNotFoundForLike:
			c.JSON(http.StatusNotFound, gin.H{"error": "Post not found"})
		case models.ErrLikeNotFound:
			c.JSON(http.StatusNotFound, gin.H{"error": "Like not found"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to unlike post"})
		}
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Post unliked successfully"})
}

func (h *LikeHandler) GetPostLikes(c *gin.Context) {
	postIDStr := c.Param("id")
	postID, err := strconv.Atoi(postIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid post ID"})
		return
	}
	var userID *int
	if userIDInterface, exists := c.Get("user_id"); exists {
		if uid, ok := userIDInterface.(int); ok {
			userID = &uid
		}
	}

	likesResponse, err := h.likeService.GetPostLikes(postID, userID)
	if err != nil {
		switch err {
		case models.ErrPostNotFoundForLike:
			c.JSON(http.StatusNotFound, gin.H{"error": "Post not found"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get post likes"})
		}
		return
	}
	c.JSON(http.StatusOK, likesResponse)
}

func (h *LikeHandler) GetUserLikes(c *gin.Context) {
	userIDInterface, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}
	userID, ok := userIDInterface.(int)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user ID format"})
		return
	}
	likes, err := h.likeService.GetUserLikes(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user likes"})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"likes": likes,
	})
}
