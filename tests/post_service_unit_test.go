package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"go-blog/handlers"
	"go-blog/models"
	"go-blog/repo"
	"go-blog/routes"
	"go-blog/service"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type TestSuite struct {
	Router *gin.Engine
}

const (
	TestTitle      = "tile 1"
	TestContent    = "content for title 1"
	UpdatedTitle   = "title - updated"
	UpdatedContent = "content is updated"
	PostID         = "2"
)

func setup() *TestSuite {
	dsn := "postgres://postgres:postgres@localhost:5432/blogdb?sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(fmt.Sprintf("couldnt connect to db : %v", err))
	}

	db.AutoMigrate(&models.Post{})
	gin.SetMode(gin.TestMode)
	repository := repo.NewPostRepository(db)
	postService := service.NewPostService(repository)
	postHandler := handlers.NewPostHandler(postService)
	router := routes.SetupRoutes(postHandler)

	return &TestSuite{Router: router}
}

func (s *TestSuite) makeRequest(method, url string, body *bytes.Buffer) *httptest.ResponseRecorder {
	var req *http.Request
	if body != nil {
		req, _ = http.NewRequest(method, url, body)
		req.Header.Set("Content-Type", "application/json")
	} else {
		req, _ = http.NewRequest(method, url, nil)
	}

	w := httptest.NewRecorder()
	s.Router.ServeHTTP(w, req)
	return w
}

func TestPostIntegration(t *testing.T) {
	t.Run("api working check", func(t *testing.T) {
		suite := setup()
		var createdPostID string

		t.Run("Create Post", func(t *testing.T) {
			postID, err := suite.createPost(TestTitle, TestContent)
			require.NoError(t, err, "post enoucnting problems")
			createdPostID = postID
		})

		t.Run("Update Post", func(t *testing.T) {
			err := suite.updatePost(createdPostID, UpdatedTitle, UpdatedContent)
			require.NoError(t, err, "post updated")
		})

		t.Run("Get Posts", func(t *testing.T) {
			posts, err := suite.getPosts()
			require.NoError(t, err, "get posts failed")
			require.NotEmpty(t, posts, "posts should not be empty")
		})

		t.Run("Delete Post", func(t *testing.T) {
			err := suite.deletePost(createdPostID)
			require.NoError(t, err, "post deleted")
		})
	})
}


func (s *TestSuite) createPost(title, content string) (string, error) {
	payload := map[string]string{"title": title, "content": content}
	body, _ := json.Marshal(payload)

	w := s.makeRequest("POST", "/api/posts", bytes.NewBuffer(body))

	if w.Code >= 400 {
		return "", fmt.Errorf("status %d", w.Code)
	}

	responseBody, _ := io.ReadAll(w.Body)
	var post models.Post
	json.Unmarshal(responseBody, &post)
	return fmt.Sprintf("%d", post.ID), nil
}

func (s *TestSuite) getPosts() ([]models.Post, error) {
	w := s.makeRequest("GET", "/api/posts", nil)

	if w.Code >= 400 {
		return nil, fmt.Errorf("status %d", w.Code)
	}

	responseBody, _ := io.ReadAll(w.Body)
	var posts []models.Post
	json.Unmarshal(responseBody, &posts)
	return posts, nil
}

func (s *TestSuite) deletePost(id string) error {
	w := s.makeRequest("DELETE", fmt.Sprintf("/api/posts/%s", id), nil)

	if w.Code >= 400 {
		return fmt.Errorf("status %d", w.Code)
	}
	return nil
}

func (s *TestSuite) updatePost(id string, title, content string) error {
	payload := map[string]string{"title": title, "content": content}
	body, _ := json.Marshal(payload)

	w := s.makeRequest("PUT", fmt.Sprintf("/api/posts/%s/update", id), bytes.NewBuffer(body))

	if w.Code >= 400 {
		return fmt.Errorf("status %d", w.Code)
	}
	return nil
}


/* individual test cases */

func TestCreatePost(t *testing.T) {
	suite := setup()
	_, err := suite.createPost(TestTitle, TestContent)
	require.NoError(t, err, "post creation failed")
}

func TestUpdatePost(t *testing.T) {
	suite := setup()
	err := suite.updatePost(PostID, UpdatedTitle, UpdatedContent)
	require.NoError(t, err, "post update failed")
}

func TestGetPosts(t *testing.T) {
	suite := setup()
	posts, err := suite.getPosts()
	require.NoError(t, err, "get posts failed")
	require.NotNil(t, posts, "posts should not be nil")
}

func TestDeletePost(t *testing.T) {
	suite := setup()
	err := suite.deletePost(PostID)
	require.NoError(t, err, "post deletion failed")
}
