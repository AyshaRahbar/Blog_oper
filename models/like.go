package models

import (
	"errors"
)

type Like struct {
	ID     int  `json:"id" gorm:"primaryKey"`
	UserID int  `json:"user_id" gorm:"not null;uniqueIndex:idx_user_post"`
	PostID int  `json:"post_id" gorm:"not null;uniqueIndex:idx_user_post"`
	User   User `json:"user" gorm:"foreignKey:UserID"`
	Post   Post `json:"post" gorm:"foreignKey:PostID"`
}
type LikeRequest struct {
	PostID int `json:"post_id"`
}
type LikeResponse struct {
	ID     int `json:"id"`
	UserID int `json:"user_id"`
	PostID int `json:"post_id"`
}
type PostLikesResponse struct {
	PostID        int  `json:"post_id"`
	LikeCount     int  `json:"like_count"`
	IsLikedByUser bool `json:"is_liked_by_user"`
}

var (
	ErrLikeNotFound        = errors.New("like not found")
	ErrLikeAlreadyExists   = errors.New("user has already liked this post")
	ErrCannotLikeOwnPost   = errors.New("cannot like your own post")
	ErrPostNotFoundForLike = errors.New("post not found for like")
	ErrUnauthorizedToLike  = errors.New("unauthorized to like posts")
	ErrLikeDatabaseError   = errors.New("like database error")
)
