package models

import (
	"errors"
	"time"
)

type Comment struct {
	ID        int       `json:"id" gorm:"primaryKey"`
	Content   string    `json:"content" gorm:"not null"`
	UserID    int       `json:"user_id" gorm:"not null"`
	PostID    int       `json:"post_id" gorm:"not null"`
	CreatedAt time.Time `json:"created_at"`
	User      User      `json:"user" gorm:"foreignKey:UserID"`
	Post      Post      `json:"post" gorm:"foreignKey:PostID"`
}

type CreateCommentRequest struct {
	Comment string `json:"comment"`
}

type UpdateCommentRequest struct {
	Comment string `json:"comment"`
}

type CommentResponse struct {
	ID        int       `json:"id"`
	Content   string    `json:"content"`
	UserID    int       `json:"user_id"`
	PostID    int       `json:"post_id"`
	Username  string    `json:"username"`
	CreatedAt time.Time `json:"created_at"`
}

type PostWithCommentsResponse struct {
	ID       int               `json:"id"`
	Title    string            `json:"title"`
	Content  string            `json:"content"`
	UserID   int               `json:"user_id"`
	Comments []CommentResponse `json:"comments"`
}

var (
	ErrCommentNotFound        = errors.New("comment not found")
	ErrCommentUnauthorized    = errors.New("unauthorized to access this comment")
	ErrCommentContentEmpty    = errors.New("comment content cannot be empty")
	ErrPostNotFoundForComment = errors.New("post not found for comment")
	ErrCommentDatabaseError   = errors.New("comment database error")
)
