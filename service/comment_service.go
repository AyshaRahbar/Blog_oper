package service

import (
	"errors"
	"fmt"
	"go-blog/models"
	"go-blog/repo"
	"strings"

	"gorm.io/gorm"
)

type CommentService interface {
	CreateComment(userID, postID int, content string) (*models.CommentResponse, error)
	GetCommentsByPost(postID int) ([]models.CommentResponse, error)
	GetCommentsByUser(userID int) ([]models.CommentResponse, error)
	UpdateComment(commentID, userID int, content string) (*models.CommentResponse, error)
	DeleteComment(commentID, userID int) error
	GetPostWithComments(postID int) (*models.PostWithCommentsResponse, error)
}

type commentService struct {
	commentRepo repo.CommentRepository
	postRepo    repo.PostRepository
}

func NewCommentService(commentRepo repo.CommentRepository, postRepo repo.PostRepository) CommentService {
	return &commentService{
		commentRepo: commentRepo,
		postRepo:    postRepo,
	}
}

func (s *commentService) CreateComment(userID, postID int, content string) (*models.CommentResponse, error) {
	if strings.TrimSpace(content) == "" {
		return nil, models.ErrCommentContentEmpty
	}

	_, err := s.postRepo.GetPost(postID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, models.ErrPostNotFoundForComment
		}
		return nil, fmt.Errorf("error checking post existence: %w", err)
	}

	comment := &models.Comment{
		Content: strings.TrimSpace(content),
		UserID:  userID,
		PostID:  postID,
	}

	createdComment, err := s.commentRepo.CreateComment(comment)
	if err != nil {
		return nil, fmt.Errorf("failed to create comment: %w", err)
	}

	return &models.CommentResponse{
		ID:        createdComment.ID,
		Content:   createdComment.Content,
		UserID:    createdComment.UserID,
		PostID:    createdComment.PostID,
		Username:  createdComment.User.Username,
		CreatedAt: createdComment.CreatedAt,
	}, nil
}

func (s *commentService) GetCommentsByPost(postID int) ([]models.CommentResponse, error) {
	_, err := s.postRepo.GetPost(postID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, models.ErrPostNotFoundForComment
		}
		return nil, fmt.Errorf("error checking post existence: %w", err)
	}
	comments, err := s.commentRepo.GetCommentsByPost(postID)
	if err != nil {
		return nil, fmt.Errorf("failed to get comments: %w", err)
	}
	commentResponses := make([]models.CommentResponse, len(comments))
	for i, comment := range comments {
		commentResponses[i] = models.CommentResponse{
			ID:        comment.ID,
			Content:   comment.Content,
			UserID:    comment.UserID,
			PostID:    comment.PostID,
			Username:  comment.User.Username,
			CreatedAt: comment.CreatedAt,
		}
	}
	return commentResponses, nil
}

func (s *commentService) GetCommentsByUser(userID int) ([]models.CommentResponse, error) {
	comments, err := s.commentRepo.GetCommentsByUser(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user comments: %w", err)
	}

	commentResponses := make([]models.CommentResponse, len(comments))
	for i, comment := range comments {
		commentResponses[i] = models.CommentResponse{
			ID:        comment.ID,
			Content:   comment.Content,
			UserID:    comment.UserID,
			PostID:    comment.PostID,
			Username:  comment.User.Username,
			CreatedAt: comment.CreatedAt,
		}
	}

	return commentResponses, nil
}

func (s *commentService) UpdateComment(commentID, userID int, content string) (*models.CommentResponse, error) {
	if strings.TrimSpace(content) == "" {
		return nil, models.ErrCommentContentEmpty
	}
	comment, err := s.commentRepo.GetComment(commentID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, models.ErrCommentNotFound
		}
		return nil, fmt.Errorf("error getting comment: %w", err)
	}
	if comment.UserID != userID {
		return nil, models.ErrCommentUnauthorized
	}
	updatedComment, err := s.commentRepo.UpdateComment(commentID, strings.TrimSpace(content))
	if err != nil {
		return nil, fmt.Errorf("failed to update comment: %w", err)
	}

	return &models.CommentResponse{
		ID:        updatedComment.ID,
		Content:   updatedComment.Content,
		UserID:    updatedComment.UserID,
		PostID:    updatedComment.PostID,
		Username:  updatedComment.User.Username,
		CreatedAt: updatedComment.CreatedAt,
	}, nil
}

func (s *commentService) DeleteComment(commentID, userID int) error {
	comment, err := s.commentRepo.GetComment(commentID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return models.ErrCommentNotFound
		}
		return fmt.Errorf("error getting comment: %w", err)
	}
	if comment.UserID != userID {
		return models.ErrCommentUnauthorized
	}
	if err := s.commentRepo.DeleteComment(commentID); err != nil {
		return fmt.Errorf("failed to delete comment: %w", err)
	}

	return nil
}

func (s *commentService) GetPostWithComments(postID int) (*models.PostWithCommentsResponse, error) {
	post, err := s.postRepo.GetPost(postID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, models.ErrPostNotFoundForComment
		}
		return nil, fmt.Errorf("error getting post: %w", err)
	}
	comments, err := s.GetCommentsByPost(postID)
	if err != nil {
		return nil, err
	}

	return &models.PostWithCommentsResponse{
		ID:       post.ID,
		Title:    post.Title,
		Content:  post.Content,
		UserID:   post.UserID,
		Comments: comments,
	}, nil
}
