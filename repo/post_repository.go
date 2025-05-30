package repo

import (
	"go-blog/models"
	"strconv"

	"gorm.io/gorm"
)

type PostRepository interface {
	ListPosts() []models.Post
	GetPost(postID string) *models.Post
	CreatePost(post *models.Post) *models.Post
	Update(id string, post *models.Post) (*models.Post, error)
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
	r.db.Raw("select * from posts").Scan(&posts)
	return posts
}

func (r *postRepository) GetPost(postID string) *models.Post {
	var post models.Post
	r.db.Raw("select * from posts where id = ?", postID).Scan(&post)
	return &post
}

func (r *postRepository) CreatePost(post *models.Post) *models.Post {
	var newPost models.Post
	r.db.Raw("insert into posts (title, content) values (?, ?) returning *", post.Title, post.Content).Scan(&newPost)
	return &newPost
}

func (r *postRepository) Update(id string, post *models.Post) (*models.Post, error) {
	intID, _ := strconv.Atoi(id)
	if err := r.db.Model(&models.Post{}).Where("id = ?", intID).Updates(models.Post{
		Title:   post.Title,
		Content: post.Content,
	}).Error; err != nil {
		return nil, err
	}

	var updatedPost models.Post
	r.db.First(&updatedPost, intID)
	return &updatedPost, nil
}

func (r *postRepository) DeletePost(postID string) {
	r.db.Exec("delete from posts where id = ?", postID)
}
