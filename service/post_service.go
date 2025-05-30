package service

import (
	"go-blog/models"
	"go-blog/repo"
)

type PostService interface {
	GetAllPosts() []models.Post
	CreatePost(post *models.Post) *models.Post
	UpdatePost(id int, post *models.Post) *models.Post
	DeletePost(id int)
}

type postService struct {
	repo repo.PostRepository
}

func NewPostService(repo repo.PostRepository) PostService {
	return &postService{repo: repo}
}

func (s *postService) GetAllPosts() []models.Post {
	return s.repo.ListPosts()
}

func (s *postService) CreatePost(post *models.Post) *models.Post {
	return s.repo.CreatePost(post)
}

func (s *postService) UpdatePost(id int, post *models.Post) *models.Post {
	post.ID = id
	updatedPost, _ := s.repo.Update(post)
	return updatedPost
}

func (s *postService) DeletePost(id int) {
	s.repo.DeletePost(id)
}
