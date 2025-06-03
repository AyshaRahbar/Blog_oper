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
	"strconv"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type TestSuite struct {
	Router *gin.Engine
}

const (
	TestTitle      = "title 1"
	TestContent    = "random content..."
	UpdatedTitle   = " title updatedd"
	UpdatedContent = "updated random blog"
)

func setup() *TestSuite {
	dsn := "postgres://postgres:postgres@localhost:5432/blogdb?sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(fmt.Sprintf("couldn't connect to db: %v", err))
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

func NewPostResponseForTest(postId string, title string, content string) *models.Post {
	id, _ := strconv.Atoi(postId)
	return &models.Post{
		ID:      id,
		Title:   title,
		Content: content,
	}
}

func NewDeletePostResponseForTest() map[string]interface{} {
	return map[string]interface{}{
		"message": "Post deleted successfully",
	}
}

func (s *TestSuite) CreatePostTest(t *testing.T, title string, content string) string {
	payload := map[string]string{"title": title, "content": content}
	body, _ := json.Marshal(payload)
	w := s.makeRequest("POST", "/api/posts", bytes.NewBuffer(body))
	if w.Code != http.StatusCreated {
		t.Fatalf("Expected status 201 Created, got %d", w.Code)
	}
	responseBody, _ := io.ReadAll(w.Body)
	var actualResponse models.Post
	json.Unmarshal(responseBody, &actualResponse)
	expectedResponse := NewPostResponseForTest(
		fmt.Sprintf("%d", actualResponse.ID),
		title,
		content,
	)
	assert.Equal(t, expectedResponse.Title, actualResponse.Title, "title couldnt form correctly")
	assert.Equal(t, expectedResponse.Content, actualResponse.Content, "content is different from expectations")
	assert.NotZero(t, actualResponse.ID, "Post ID should be set")
	return fmt.Sprintf("%d", actualResponse.ID)
}

func (s *TestSuite) GetPostsTest(t *testing.T) []models.Post {
	w := s.makeRequest("GET", "/api/posts", nil)
	if w.Code != http.StatusOK {
		t.Fatalf("Expected status 200 OK, got %d", w.Code)
	}
	responseBody, _ := io.ReadAll(w.Body)
	var actualResponse []models.Post
	json.Unmarshal(responseBody, &actualResponse)
	assert.True(t, len(actualResponse) >= 0, "response can be empty or have posts")
	return actualResponse
}

func (s *TestSuite) UpdatePostTest(t *testing.T, id string, title string, content string) {
	payload := map[string]string{"title": title, "content": content}
	body, _ := json.Marshal(payload)
	w := s.makeRequest("PUT", fmt.Sprintf("/api/posts/%s/update", id), bytes.NewBuffer(body))
	if w.Code != http.StatusOK {
		t.Fatalf("Expected status 200 OK, got %d", w.Code)
	}
	responseBody, _ := io.ReadAll(w.Body)
	var actualResponse models.Post
	json.Unmarshal(responseBody, &actualResponse)
	expectedResponse := NewPostResponseForTest(id, title, content)
	assert.Equal(t, expectedResponse.Title, actualResponse.Title, "title could not update correctly")
	assert.Equal(t, expectedResponse.Content, actualResponse.Content, "content is different from expectations")
	assert.NotZero(t, actualResponse.ID, "Post ID should be preserved")
}

func (s *TestSuite) DeletePostTest(t *testing.T, id string) {
	w := s.makeRequest("DELETE", fmt.Sprintf("/api/posts/%s", id), nil)
	if w.Code != http.StatusOK {
		t.Fatalf("Expected status 200 OK, got %d", w.Code)
	}
	responseBody, _ := io.ReadAll(w.Body)
	var actualResponse map[string]interface{}
	json.Unmarshal(responseBody, &actualResponse)
	expectedResponse := NewDeletePostResponseForTest()
	assert.Equal(t, expectedResponse, actualResponse, "delete response did not match expected response")
}

func TestPostIntegration(t *testing.T) {
	t.Run("Complete Post Lifecycle", func(t *testing.T) {
		suite := setup()
		var createdPostID string

		t.Run("Create Post", func(t *testing.T) {
			createdPostID = suite.CreatePostTest(t, TestTitle, TestContent)
			require.NotEmpty(t, createdPostID, "Created post ID should not be empty")
		})

		t.Run("Get All Posts", func(t *testing.T) {
			posts := suite.GetPostsTest(t)
			require.NotEmpty(t, posts, "Posts list should not be empty")
		})

		t.Run("Update Post", func(t *testing.T) {
			suite.UpdatePostTest(t, createdPostID, UpdatedTitle, UpdatedContent)
		})

		t.Run("Delete Post", func(t *testing.T) {
			suite.DeletePostTest(t, createdPostID)
		})
	})
}

func TestCreatePost(t *testing.T) {
	suite := setup()
	suite.CreatePostTest(t, TestTitle, TestContent)
}

func TestGetPosts(t *testing.T) {
	suite := setup()
	suite.GetPostsTest(t)
}

func TestUpdatePost(t *testing.T) {
	suite := setup()
	createdID := suite.CreatePostTest(t, TestTitle, TestContent)
	suite.UpdatePostTest(t, createdID, UpdatedTitle, UpdatedContent)
}

func TestDeletePost(t *testing.T) {
	suite := setup()
	createdID := suite.CreatePostTest(t, TestTitle, TestContent)
	suite.DeletePostTest(t, createdID)
}
