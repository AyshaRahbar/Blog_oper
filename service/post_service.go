package service

import (
	"errors"
	"go-blog/models"
	"go-blog/repo"
	"strings"
)

type PostService interface {
	GetAllPosts() ([]models.Post, error)
	CreatePost(post *models.Post) (*models.Post, error)
	UpdatePost(id string, post *models.Post) (*models.Post, error)
	DeletePost(id string) error
}

type postService struct {
	repo repo.PostRepository
}

func NewPostService(repo repo.PostRepository) PostService {
	return &postService{repo: repo}
}

func (s *postService) GetAllPosts() ([]models.Post, error) {
	return s.repo.ListPosts()
}

func (s *postService) CreatePost(post *models.Post) (*models.Post, error) {
	if strings.TrimSpace(post.Title) == "" {
		return nil, errors.New("title cannot be empty")
	}
	if strings.TrimSpace(post.Content) == "" {
		return nil, errors.New("content cannot be empty")
	}

	return s.repo.CreatePost(post)
}

func (s *postService) UpdatePost(id string, post *models.Post) (*models.Post, error) {
	if strings.TrimSpace(post.Title) == "" {
		return nil, errors.New("title cannot be empty")
	}
	if strings.TrimSpace(post.Content) == "" {
		return nil, errors.New("content cannot be empty")
	}

	return s.repo.Update(id, post)
}

func (s *postService) DeletePost(id string) error {
	return s.repo.DeletePost(id)
}
