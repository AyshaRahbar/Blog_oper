package models

import (
	"errors"
)

type Post struct {
	ID      int    `json:"id" gorm:"primaryKey"`
	Title   string `json:"title"`
	Content string `json:"content"`
	UserID  int    `json:"user_id" gorm:"not null"`
}

type UpdatePostRequest struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}

var (
	ErrPostNotFound = errors.New("post not found")
	ErrInvalidPost  = errors.New("invalid post")
)
