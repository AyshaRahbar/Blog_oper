package service

import (
	"errors"
	"go-blog/models"
	"go-blog/repo"
)

type UserService struct {
	repo *repo.UserRepository
}

func NewUserService(repo *repo.UserRepository) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) Register(req *models.RegisterRequest) (*models.User, error) {
	exists, err := s.repo.UsernameExists(req.Username)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errors.New("username exists")
	}

	user := &models.User{
		Username:    req.Username,
		Password:    req.Password,
		AccountType: req.AccountType,
	}

	return s.repo.CreateUser(user)
}

func (s *UserService) Login(username, password string) (*models.User, error) {
	user, err := s.repo.GetUserByUsername(username)
	if err != nil {
		return nil, err
	}
	if user.Password != password {
		return nil, errors.New("invalid password")
	}
	return user, nil
}
