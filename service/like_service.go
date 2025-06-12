package service

import (
	"errors"
	"fmt"
	"go-blog/models"
	"go-blog/repo"

	"gorm.io/gorm"
)

type LikeService interface {
	LikePost(userID, postID int) (*models.LikeResponse, error)
	UnlikePost(userID, postID int) error
	GetPostLikes(postID int, userID *int) (*models.PostLikesResponse, error)
	GetPostLikesCount(postID int) (*models.PostLikesDetailResponse, error)
	GetUserLikes(userID int) ([]models.LikeResponse, error)
}

type likeService struct {
	likeRepo repo.LikeRepository
	postRepo repo.PostRepository
}

func NewLikeService(likeRepo repo.LikeRepository, postRepo repo.PostRepository) LikeService {
	return &likeService{
		likeRepo: likeRepo,
		postRepo: postRepo,
	}
}

func (s *likeService) LikePost(userID, postID int) (*models.LikeResponse, error) {
	post, err := s.postRepo.GetPost(postID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, models.ErrPostNotFoundForLike
		}
		return nil, fmt.Errorf("error checking post existence: %w", err)
	}
	if post == nil {
		return nil, models.ErrPostNotFoundForLike
	}
	if post.UserID == userID {
		return nil, models.ErrCannotLikeOwnPost
	}
	existingLike, err := s.likeRepo.GetLike(userID, postID)
	if err == nil && existingLike != nil {
		return nil, models.ErrLikeAlreadyExists
	}
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("error checking existing like: %w", err)
	}
	like := &models.Like{
		UserID: userID,
		PostID: postID,
	}
	createdLike, err := s.likeRepo.CreateLike(like)
	if err != nil {
		return nil, fmt.Errorf("failed to create like: %w", err)
	}
	return &models.LikeResponse{
		ID:     createdLike.ID,
		UserID: createdLike.UserID,
		PostID: createdLike.PostID,
	}, nil
}

func (s *likeService) UnlikePost(userID, postID int) error {
	_, err := s.postRepo.GetPost(postID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return models.ErrPostNotFoundForLike
		}
		return fmt.Errorf("error checking post existence: %w", err)
	}
	existingLike, err := s.likeRepo.GetLike(userID, postID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return models.ErrLikeNotFound
		}
		return fmt.Errorf("error checking existing like: %w", err)
	}
	if existingLike == nil {
		return models.ErrLikeNotFound
	}
	err = s.likeRepo.DeleteLike(userID, postID)
	if err != nil {
		return fmt.Errorf("failed to delete like: %w", err)
	}
	return nil
}

func (s *likeService) GetPostLikes(postID int, userID *int) (*models.PostLikesResponse, error) {
	_, err := s.postRepo.GetPost(postID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, models.ErrPostNotFoundForLike
		}
		return nil, fmt.Errorf("error checking post existence: %w", err)
	}
	likeCount, err := s.likeRepo.GetLikeCount(postID)
	if err != nil {
		return nil, fmt.Errorf("failed to get like count: %w", err)
	}
	isLikedByUser := false
	if userID != nil {
		isLikedByUser, err = s.likeRepo.IsPostLikedByUser(*userID, postID)
		if err != nil {
			return nil, fmt.Errorf("failed to check if post is liked by user: %w", err)
		}
	}

	return &models.PostLikesResponse{
		PostID:        postID,
		LikeCount:     int(likeCount),
		IsLikedByUser: isLikedByUser,
	}, nil
}

func (s *likeService) GetUserLikes(userID int) ([]models.LikeResponse, error) {
	likes, err := s.likeRepo.GetLikesByUser(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user likes: %w", err)
	}

	likeResponses := make([]models.LikeResponse, len(likes))
	for i, like := range likes {
		likeResponses[i] = models.LikeResponse{
			ID:     like.ID,
			UserID: like.UserID,
			PostID: like.PostID,
		}
	}

	return likeResponses, nil
}

func (s *likeService) GetPostLikesCount(postID int) (*models.PostLikesDetailResponse, error) {
	_, err := s.postRepo.GetPost(postID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, models.ErrPostNotFoundForLike
		}
		return nil, fmt.Errorf("error checking post existence: %w", err)
	}
	likes, err := s.likeRepo.GetLikesByPost(postID)
	if err != nil {
		return nil, fmt.Errorf("failed to get post likes: %w", err)
	}
	likeResponses := make([]models.LikeWithUserResponse, len(likes))
	for i, like := range likes {
		likeResponses[i] = models.LikeWithUserResponse{
			ID:       like.ID,
			UserID:   like.UserID,
			PostID:   like.PostID,
			Username: like.User.Username,
		}
	}
	return &models.PostLikesDetailResponse{
		PostID:    postID,
		LikeCount: len(likes),
		Likes:     likeResponses,
	}, nil
}
