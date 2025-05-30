package service

import (
	"go-blog/models"
	"go-blog/repo"
)

type PostService interface {
	GetPosts() []models.Post
	GetPost(id string) *models.Post
	CreatePost(post *models.Post) *models.Post
	UpdatePost(id string, post *models.Post)
	DeletePost(id string)
}

type postService struct {
	repo repo.PostRepository
}

func NewPostService(repo repo.PostRepository) PostService {
	return &postService{repo: repo}
}

func (s *postService) GetPosts() []models.Post {
	return s.repo.ListPosts()
}

func (s *postService) GetPost(id string) *models.Post {
	return s.repo.GetPost(id)
}

func (s *postService) CreatePost(post *models.Post) *models.Post {
	return s.repo.CreatePost(post)
}

func (s *postService) UpdatePost(id string, post *models.Post) {
	s.repo.Update(post)
}

func (s *postService) DeletePost(id string) {
	s.repo.DeletePost(id)
}
