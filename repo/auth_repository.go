package repo

import (
	"go-blog/models"
)

type AuthRepository interface {
	CheckPostOwnership(postID int, userID int) error
}

type authRepository struct {
	postRepo PostRepository
}

func NewAuthRepository(postRepo PostRepository) AuthRepository {
	return &authRepository{postRepo: postRepo}
}

func (r *authRepository) CheckPostOwnership(postID int, userID int) error {
	post, err := r.postRepo.GetPost(postID)
	if err != nil {
		return models.ErrPostNotFound
	}

	if post.UserID != userID {
		return models.ErrPostUnauthorized
	}

	return nil
}
