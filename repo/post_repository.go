package repo

import (
	"go-blog/models"

	"gorm.io/gorm"
)

type PostRepository interface {
	ListPosts() ([]models.Post, error)
	GetPost(postID int) (*models.Post, error)
	CreatePost(post *models.Post) (*models.Post, error)
	Update(id int, post *models.Post) (*models.Post, error)
	DeletePost(postID int) error
}

type postRepository struct {
	db *gorm.DB
}

func NewPostRepository(db *gorm.DB) PostRepository {
	return &postRepository{db: db}
}
func (r *postRepository) ListPosts() ([]models.Post, error) {
	var posts []models.Post
	if err := r.db.Find(&posts).Error; err != nil {
		return nil, err
	}
	return posts, nil
}
func (r *postRepository) GetPost(postID int) (*models.Post, error) {
	var post models.Post
	if err := r.db.First(&post, "id = ?", postID).Error; err != nil {
		return nil, err
	}
	return &post, nil
}

func (r *postRepository) CreatePost(post *models.Post) (*models.Post, error) {
	if err := r.db.Create(post).Error; err != nil {
		return nil, err
	}
	return post, nil
}

func (r *postRepository) Update(id int, post *models.Post) (*models.Post, error) {
	if err := r.db.Model(&models.Post{}).Where("id = ?", id).Updates(models.Post{
		Title:   post.Title,
		Content: post.Content,
	}).Error; err != nil {
		return nil, err
	}

	var updatedPost models.Post
	if err := r.db.First(&updatedPost, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &updatedPost, nil
}

func (r *postRepository) DeletePost(postID int) error {
	if err := r.db.Delete(&models.Post{}, "id = ?", postID).Error; err != nil {
		return err
	}
	return nil
}
