package handlers

import (
	"go-blog/models"
	"go-blog/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type CommentHandler struct {
	commentService service.CommentService
}

func NewCommentHandler(commentService service.CommentService) *CommentHandler {
	return &CommentHandler{commentService: commentService}
}

func (h *CommentHandler) CreateComment(c *gin.Context) {
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

	var req models.CreateCommentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}

	commentResponse, err := h.commentService.CreateComment(userID, postID, req.Comment)
	if err != nil {
		switch err {
		case models.ErrCommentContentEmpty:
			c.JSON(http.StatusBadRequest, gin.H{"error": "Comment content cannot be empty"})
		case models.ErrPostNotFoundForComment:
			c.JSON(http.StatusNotFound, gin.H{"error": "Post not found"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create comment"})
		}
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Comment created successfully",
		"comment": commentResponse,
	})
}

func (h *CommentHandler) GetPostComments(c *gin.Context) {
	postIDStr := c.Param("id")
	postID, err := strconv.Atoi(postIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid post ID"})
		return
	}

	comments, err := h.commentService.GetCommentsByPost(postID)
	if err != nil {
		switch err {
		case models.ErrPostNotFoundForComment:
			c.JSON(http.StatusNotFound, gin.H{"error": "Post not found"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get comments"})
		}
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"comments": comments,
	})
}

func (h *CommentHandler) GetPostWithComments(c *gin.Context) {
	postIDStr := c.Param("id")
	postID, err := strconv.Atoi(postIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid post ID"})
		return
	}

	postWithComments, err := h.commentService.GetPostWithComments(postID)
	if err != nil {
		switch err {
		case models.ErrPostNotFoundForComment:
			c.JSON(http.StatusNotFound, gin.H{"error": "Post not found"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get post with comments"})
		}
		return
	}
	c.JSON(http.StatusOK, postWithComments)
}

func (h *CommentHandler) GetUserComments(c *gin.Context) {
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
	comments, err := h.commentService.GetCommentsByUser(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user comments"})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"comments": comments,
	})
}

func (h *CommentHandler) UpdateComment(c *gin.Context) {
	commentIDStr := c.Param("id")
	commentID, err := strconv.Atoi(commentIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid comment ID"})
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
	var req models.UpdateCommentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}

	commentResponse, err := h.commentService.UpdateComment(commentID, userID, req.Comment)
	if err != nil {
		switch err {
		case models.ErrCommentNotFound:
			c.JSON(http.StatusNotFound, gin.H{"error": "Comment not found"})
		case models.ErrCommentUnauthorized:
			c.JSON(http.StatusForbidden, gin.H{"error": "Cannot update comment - not your comment"})
		case models.ErrCommentContentEmpty:
			c.JSON(http.StatusBadRequest, gin.H{"error": "Comment content cannot be empty"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update comment"})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Comment updated successfully",
		"comment": commentResponse,
	})
}

func (h *CommentHandler) DeleteComment(c *gin.Context) {
	commentIDStr := c.Param("id")
	commentID, err := strconv.Atoi(commentIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid comment ID"})
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
	err = h.commentService.DeleteComment(commentID, userID)
	if err != nil {
		switch err {
		case models.ErrCommentNotFound:
			c.JSON(http.StatusNotFound, gin.H{"error": "Comment not found"})
		case models.ErrCommentUnauthorized:
			c.JSON(http.StatusForbidden, gin.H{"error": "Cannot delete comment - not your comment"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete comment"})
		}
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Comment deleted successfully"})
}
