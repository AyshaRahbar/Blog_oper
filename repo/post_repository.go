package repo

import (
	"go-blog/models"

	"gorm.io/gorm"
)

type PostRepository interface {
	ListPosts() []models.Post
	GetPost(postID int) *models.Post
	CreatePost(post *models.Post) *models.Post
	Update(post *models.Post) (*models.Post, error)
	DeletePost(postID int)
}

type postRepository struct {
	db *gorm.DB
}

func NewPostRepository(db *gorm.DB) PostRepository {
	return &postRepository{db: db}
}

func (r *postRepository) ListPosts() []models.Post {
	var posts []models.Post
	r.db.Raw("SELECT * FROM posts").Scan(&posts)
	return posts
}

func (r *postRepository) GetPost(postID int) *models.Post {
	var post models.Post
	r.db.Raw("SELECT * FROM posts WHERE id = ?", postID).Scan(&post)
	return &post
}

func (r *postRepository) CreatePost(post *models.Post) *models.Post {
	var newPost models.Post
	r.db.Raw("INSERT INTO posts (title, content) VALUES (?, ?) RETURNING *", post.Title, post.Content).Scan(&newPost)
	return &newPost
}

func (r *postRepository) Update(post *models.Post) (*models.Post, error) {
	if err := r.db.Model(&models.Post{}).Where("id = ?", post.ID).Updates(models.Post{
		Title:   post.Title,
		Content: post.Content,
	}).Error; err != nil {
		return nil, err
	}
	var updatedPost models.Post
	r.db.First(&updatedPost, post.ID)
	return &updatedPost, nil
}

func (r *postRepository) DeletePost(postID int) {
	r.db.Exec("DELETE FROM posts WHERE id = ?", postID)
}
