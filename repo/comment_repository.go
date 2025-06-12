package repo

import (
	"go-blog/models"

	"gorm.io/gorm"
)

type CommentRepository interface {
	CreateComment(comment *models.Comment) (*models.Comment, error)
	GetComment(commentID int) (*models.Comment, error)
	GetCommentsByPost(postID int) ([]models.Comment, error)
	GetCommentsByUser(userID int) ([]models.Comment, error)
	UpdateComment(commentID int, content string) (*models.Comment, error)
	DeleteComment(commentID int) error
	GetCommentCount(postID int) (int, error)
}

type commentRepository struct {
	db *gorm.DB
}

func NewCommentRepository(db *gorm.DB) CommentRepository {
	return &commentRepository{db: db}
}

func (r *commentRepository) CreateComment(comment *models.Comment) (*models.Comment, error) {
	if err := r.db.Create(comment).Error; err != nil {
		return nil, err
	}
	if err := r.db.Preload("User").First(comment, comment.ID).Error; err != nil {
		return nil, err
	}
	return comment, nil
}

func (r *commentRepository) GetComment(commentID int) (*models.Comment, error) {
	var comment models.Comment
	if err := r.db.Preload("User").Preload("Post").First(&comment, commentID).Error; err != nil {
		return nil, err
	}
	return &comment, nil
}

func (r *commentRepository) GetCommentsByPost(postID int) ([]models.Comment, error) {
	var comments []models.Comment
	if err := r.db.Where("post_id = ?", postID).Preload("User").Order("created_at ASC").Find(&comments).Error; err != nil {
		return nil, err
	}
	return comments, nil
}

func (r *commentRepository) GetCommentsByUser(userID int) ([]models.Comment, error) {
	var comments []models.Comment
	if err := r.db.Where("user_id = ?", userID).Preload("Post").Order("created_at DESC").Find(&comments).Error; err != nil {
		return nil, err
	}
	return comments, nil
}

func (r *commentRepository) UpdateComment(commentID int, content string) (*models.Comment, error) {
	if err := r.db.Model(&models.Comment{}).Where("id = ?", commentID).Update("content", content).Error; err != nil {
		return nil, err
	}
	var updatedComment models.Comment
	if err := r.db.Preload("User").First(&updatedComment, commentID).Error; err != nil {
		return nil, err
	}
	return &updatedComment, nil
}

func (r *commentRepository) DeleteComment(commentID int) error {
	if err := r.db.Delete(&models.Comment{}, commentID).Error; err != nil {
		return err
	}
	return nil
}

func (r *commentRepository) GetCommentCount(postID int) (int, error) {
	var comments []models.Comment
	if err := r.db.Where("post_id = ?", postID).Find(&comments).Error; err != nil {
		return 0, err
	}
	return len(comments), nil
}
