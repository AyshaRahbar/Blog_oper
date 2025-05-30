package handlers

import (
	"go-blog/models"
	"go-blog/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type PostHandler struct {
	service service.PostService
}

func NewPostHandler(service service.PostService) *PostHandler {
	return &PostHandler{service: service}
}

func (h *PostHandler) GetPosts(c *gin.Context) {
	posts := h.service.GetAllPosts()
	c.JSON(http.StatusOK, posts)
}

func (h *PostHandler) CreatePost(c *gin.Context) {
	var post models.Post
	c.ShouldBindJSON(&post)
	createdPost := h.service.CreatePost(&post)
	c.JSON(http.StatusCreated, createdPost)
}

func (h *PostHandler) GetPost(c *gin.Context) {
	id := c.Param("id")
	post := h.service.GetPostByID(id)
	c.JSON(http.StatusOK, post)
}

func (h *PostHandler) UpdatePost(c *gin.Context) {
	id := c.Param("id")
	var post models.Post
	c.ShouldBindJSON(&post)
	h.service.UpdatePost(id, &post)
	c.JSON(http.StatusOK, post)
}

func (h *PostHandler) DeletePost(c *gin.Context) {
	id := c.Param("id")
	h.service.DeletePost(id)
	c.JSON(http.StatusOK, gin.H{"message": "Post deleted successfully"})
}
