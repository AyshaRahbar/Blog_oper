package service

import (
	"errors"
	"fmt"
	"go-blog/models"
	"go-blog/repo"
	"strings"
)

type PostService interface {
	GetAllPosts() ([]models.Post, error)
	GetPostByID(id int) (*models.Post, error)
	CreatePost(post *models.Post) (*models.Post, error)
	UpdatePost(id int, post *models.Post) (*models.Post, error)
	DeletePost(id int) error
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

	createdPost, err := s.repo.CreatePost(post)
	if err != nil {
		return nil, fmt.Errorf("failed to create post: %w", err)
	}

	if createdPost == nil {
		return nil, errors.New("post creation failed - no post returned")
	}

	if createdPost.Title != post.Title || createdPost.Content != post.Content {
		return nil, errors.New("post creation failed - data mismatch")
	}

	return createdPost, nil
}

func (s *postService) UpdatePost(id int, post *models.Post) (*models.Post, error) {
	if strings.TrimSpace(post.Title) == "" {
		return nil, errors.New("title cannot be empty")
	}
	if strings.TrimSpace(post.Content) == "" {
		return nil, errors.New("content cannot be empty")
	}

	beforePosts, err := s.repo.GetPost(id)
	if err != nil {
		return nil, fmt.Errorf("post with ID %d does not exist", id)
	}
	if beforePosts == nil {
		return nil, fmt.Errorf("post with ID %d not found", id)
	}

	updatedPost, err := s.repo.Update(id, post)
	if err != nil {
		return nil, fmt.Errorf("failed to update post: %w", err)
	}

	if updatedPost == nil {
		return nil, errors.New("post update failed - no post returned")
	}

	if updatedPost.Title != post.Title || updatedPost.Content != post.Content {
		return nil, fmt.Errorf("post with ID %d was not updated", id)
	}

	if updatedPost.ID != id {
		return nil, errors.New("post update failed - ID mismatch")
	}

	return updatedPost, nil
}

func (s *postService) DeletePost(id int) error {
	Prevpost, err := s.repo.GetPost(id)
	if err != nil {
		return fmt.Errorf("post with ID %d does not exist to delete", id)
	}
	if Prevpost == nil {
		return fmt.Errorf("post with ID %d not found", id)
	}

	err = s.repo.DeletePost(id)
	if err != nil {
		return fmt.Errorf("failed to delete post: %w", err)
	}

	afterPosts, err := s.repo.GetPost(id)
	if err == nil && afterPosts != nil {
		return fmt.Errorf("deletion failed for %d", id)
	}

	return nil
}

func (s *postService) GetPostByID(id int) (*models.Post, error) {
	post, err := s.repo.GetPost(id)
	if err != nil {
		return nil, fmt.Errorf("post not found %d", id)
	}
	if post == nil {
		return nil, fmt.Errorf("post doesnt exist %d", id)
	}
	return post, nil
}
