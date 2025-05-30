package repo

import (
	"go-blog/models"

	"gorm.io/gorm"
)

type PostRepository interface {
	ListPosts() []models.Post
	GetByID(postID string) *models.Post
	CreatePost(post *models.Post) *models.Post
	Update(post *models.Post)
	DeletePost(postID string)
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

func (r *postRepository) GetByID(postID string) *models.Post {
	var post models.Post
	r.db.Raw("SELECT * FROM posts WHERE id = ?", postID).Scan(&post)
	return &post
}

func (r *postRepository) CreatePost(post *models.Post) *models.Post {
	r.db.Create(post)
	return post
}

func (r *postRepository) Update(post *models.Post) {
	r.db.Exec("UPDATE posts SET title = ?, content = ? WHERE id = ?", post.Title, post.Content, post.ID)
}

func (r *postRepository) DeletePost(postID string) {
	r.db.Exec("DELETE FROM posts WHERE id = ?", postID)
}
