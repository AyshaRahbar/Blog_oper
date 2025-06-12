package repo

import (
	"errors"
	"go-blog/models"
	"gorm.io/gorm"
)

type LikeRepository interface {
	CreateLike(like *models.Like) (*models.Like, error)
	DeleteLike(userID, postID int) error
	GetLike(userID, postID int) (*models.Like, error)
	GetLikesByPost(postID int) ([]models.Like, error)
	GetLikesByUser(userID int) ([]models.Like, error)
	GetLikeCount(postID int) (int, error)
	IsPostLikedByUser(userID, postID int) (bool, error)
}

type likeRepository struct {
	db *gorm.DB
}

func NewLikeRepository(db *gorm.DB) LikeRepository {
	return &likeRepository{db: db}
}

func (r *likeRepository) CreateLike(like *models.Like) (*models.Like, error) {
	if err := r.db.Create(like).Error; err != nil {
		return nil, err
	}
	return like, nil
}

func (r *likeRepository) DeleteLike(userID, postID int) error {
	if err := r.db.Where("user_id = ? AND post_id = ?", userID, postID).Delete(&models.Like{}).Error; err != nil {
		return err
	}
	return nil
}

func (r *likeRepository) GetLike(userID, postID int) (*models.Like, error) {
	var like models.Like
	if err := r.db.Where("user_id = ? AND post_id = ?", userID, postID).First(&like).Error; err != nil {
		return nil, err
	}
	return &like, nil
}

func (r *likeRepository) GetLikesByPost(postID int) ([]models.Like, error) {
	var likes []models.Like
	if err := r.db.Where("post_id = ?", postID).Preload("User").Find(&likes).Error; err != nil {
		return nil, err
	}
	return likes, nil
}

func (r *likeRepository) GetLikesByUser(userID int) ([]models.Like, error) {
	var likes []models.Like
	if err := r.db.Where("user_id = ?", userID).Preload("Post").Find(&likes).Error; err != nil {
		return nil, err
	}
	return likes, nil
}

func (r *likeRepository) GetLikeCount(postID int) (int, error) {
	var count int64
	if err := r.db.Model(&models.Like{}).Where("post_id = ?", postID).Count(&count).Error; err != nil {
		return 0, err
	}
	return int(count), nil
}

func (r *likeRepository) IsPostLikedByUser(userID, postID int) (bool, error) {
	var like models.Like
	err := r.db.Where("user_id = ? AND post_id = ?", userID, postID).First(&like).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}
