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

	if w.Code != http.StatusCreated {
		return "", fmt.Errorf("expected status %d, got %d", http.StatusCreated, w.Code)
	}

	responseBody, _ := io.ReadAll(w.Body)
	var post models.Post
	json.Unmarshal(responseBody, &post)
	createdID := fmt.Sprintf("%d", post.ID)

	posts, err := s.getPosts()
	if err != nil {
		return "", fmt.Errorf("verification failed: %v", err)
	}

	postExists := false
	for _, check := range posts {
		if fmt.Sprintf("%d", check.ID) == createdID && check.Title == title && check.Content == content {
			postExists = true
			break
		}
	}
	if !postExists {
		return "", fmt.Errorf("post not found in database")
	}
	return createdID, nil
}

func (s *TestSuite) getPosts() ([]models.Post, error) {
	w := s.makeRequest("GET", "/api/posts", nil)

	if w.Code != http.StatusOK {
		return nil, fmt.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
	}

	responseBody, _ := io.ReadAll(w.Body)
	if len(responseBody) == 0 {
		fmt.Println("Empty response body, returning empty posts slice")
		return []models.Post{}, nil
	}

	var posts []models.Post
	err := json.Unmarshal(responseBody, &posts)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %v", err)
	}

	if posts == nil {
		posts = []models.Post{}
	}

	fmt.Printf("Successfully retrieved %d posts\n", len(posts))
	return posts, nil
}

func (s *TestSuite) deletePost(id string) error {
	fmt.Printf("delete post with ID: %s\n", id)

	Prevpost, err := s.getPosts()
	if err != nil {
		return fmt.Errorf("Could not get posts before deletion: %v", err)
	}

	postExists := false
	for _, post := range Prevpost {
		if fmt.Sprintf("%d", post.ID) == id {
			postExists = true
			break
		}
	}
	if !postExists {
		return fmt.Errorf("post with ID %s does not exist before deletion", id)
	}

	w := s.makeRequest("DELETE", fmt.Sprintf("/api/posts/%s", id), nil)

	if w.Code != http.StatusOK {
		return fmt.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
	}

	afterPosts, err := s.getPosts()
	if err != nil {
		return fmt.Errorf("verification failed - could not get posts after deletion: %v", err)
	}
	postStillPresent := false
	for _, post := range afterPosts {
		if fmt.Sprintf("%d", post.ID) == id {
			postStillPresent = true
			break
		}
	}
	if postStillPresent {
		return fmt.Errorf("post with ID %s still exists after deletion", id)
	}

	fmt.Printf("Successfully deleted post with ID: %s (database now has %d posts)\n", id, len(afterPosts))
	return nil
}

func (s *TestSuite) updatePost(id string, title, content string) error {
	fmt.Printf("Update Post for ID: %s with title: %s\n", id, title)

	beforePosts, err := s.getPosts()
	if err != nil {
		return fmt.Errorf("could not get posts: %v", err)
	}

	AlrExisting := false
	for _, post := range beforePosts {
		if fmt.Sprintf("%d", post.ID) == id {
			AlrExisting = true
			break
		}
	}
	if !AlrExisting {
		return fmt.Errorf("post with ID %s does not exist for update", id)
	}

	payload := map[string]string{"title": title, "content": content}
	body, _ := json.Marshal(payload)

	w := s.makeRequest("PUT", fmt.Sprintf("/api/posts/%s/update", id), bytes.NewBuffer(body))

	if w.Code != http.StatusOK {
		return fmt.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
	}

	afterPosts, err := s.getPosts()
	if err != nil {
		return fmt.Errorf("could not get posts: %v", err)
	}

	postUpdated := false
	for _, post := range afterPosts {
		if fmt.Sprintf("%d", post.ID) == id && post.Title == title && post.Content == content {
			postUpdated = true
			break
		}
	}
	if postUpdated == false {
		return fmt.Errorf("post with ID %s was not updated ", id)
	}

	fmt.Printf("post updated")
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
