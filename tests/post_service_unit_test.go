package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"go-blog/models"
	"io"
	"net/http"
	"testing"
	"go-blog/testutils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	TestTitle      = "title 1"
	TestContent    = "random content..."
	UpdatedTitle   = " title updatedd"
	UpdatedContent = "updated random blog"
)

func newPostResponseForTest(id int, title string, content string) *models.Post {
	return &models.Post{
		ID:      id,
		Title:   title,
		Content: content,
	}
}

func newDeletePostResponseForTest() map[string]interface{} {
	return map[string]interface{}{
		"message": "Post deleted successfully",
	}
}

func createPostTest(t *testing.T, suite *testutils.TestSuite, title string, content string) int {
	payload := map[string]string{"title": title, "content": content}
	body, _ := json.Marshal(payload)
	w := suite.MakeRequest("POST", "/api/posts", bytes.NewBuffer(body))
	if w.Code != http.StatusCreated {
		t.Fatalf("Expected status 201 Created, got %d", w.Code)
	}
	responseBody, _ := io.ReadAll(w.Body)
	var actualResponse models.Post
	json.Unmarshal(responseBody, &actualResponse)
	expectedResponse := newPostResponseForTest(
		actualResponse.ID,
		title,
		content,
	)
	assert.Equal(t, expectedResponse.Title, actualResponse.Title, "title couldnt form correctly")
	assert.Equal(t, expectedResponse.Content, actualResponse.Content, "content is different from expectations")
	assert.NotZero(t, actualResponse.ID, "Post ID should be set")

	fetchedPost, err := suite.PostRepo.GetPost(actualResponse.ID)
	require.NoError(t, err, "no error while fetching from db")
	require.NotNil(t, fetchedPost, "post not nil in db")
	assert.Equal(t, actualResponse.ID, fetchedPost.ID, "postId matches")
	assert.Equal(t, title, fetchedPost.Title, "title should match")
	assert.Equal(t, content, fetchedPost.Content, "content should match")

	return actualResponse.ID
}

func getPostsTest(t *testing.T, suite *testutils.TestSuite) []models.Post {
	w := suite.MakeRequest("GET", "/api/posts", nil)
	if w.Code != http.StatusOK {
		t.Fatalf("Expected status 200 OK, got %d", w.Code)
	}
	responseBody, _ := io.ReadAll(w.Body)
	var actualResponse []models.Post
	json.Unmarshal(responseBody, &actualResponse)
	assert.True(t, len(actualResponse) >= 0, "response can be empty or have posts")
	return actualResponse
}

func updatePostTest(t *testing.T, suite *testutils.TestSuite, id int, title string, content string) {
	payload := map[string]string{"title": title, "content": content}
	body, _ := json.Marshal(payload)
	w := suite.MakeRequest("PUT", fmt.Sprintf("/api/posts/%d/update", id), bytes.NewBuffer(body))
	if w.Code != http.StatusOK {
		t.Fatalf("Expected status 200 OK, got %d", w.Code)
	}
	responseBody, _ := io.ReadAll(w.Body)
	var actualResponse models.Post
	json.Unmarshal(responseBody, &actualResponse)
	expectedResponse := newPostResponseForTest(id, title, content)
	assert.Equal(t, expectedResponse.Title, actualResponse.Title, "title could not update correctly")
	assert.Equal(t, expectedResponse.Content, actualResponse.Content, "content is different from expectations")
	assert.NotZero(t, actualResponse.ID, "Post ID should be preserved")
}

func deletePostTest(t *testing.T, suite *testutils.TestSuite, id int) {
	w := suite.MakeRequest("DELETE", fmt.Sprintf("/api/posts/%d", id), nil)
	if w.Code != http.StatusOK {
		t.Fatalf("Expected status 200 OK, got %d", w.Code)
	}
	responseBody, _ := io.ReadAll(w.Body)
	var actualResponse map[string]interface{}
	json.Unmarshal(responseBody, &actualResponse)
	expectedResponse := newDeletePostResponseForTest()
	assert.Equal(t, expectedResponse, actualResponse, "delete response did not match expected response")
}

func getPostByIDTest(t *testing.T, suite *testutils.TestSuite, id int, expectedTitle, expectedContent string) {
	w := suite.MakeRequest("GET", fmt.Sprintf("/api/posts/%d", id), nil)
	if w.Code != http.StatusOK {
		t.Fatalf("Expected status 200 OK, got %d", w.Code)
	}
	responseBody, _ := io.ReadAll(w.Body)
	var actualResponse models.Post
	json.Unmarshal(responseBody, &actualResponse)
	assert.Equal(t, id, actualResponse.ID, "id should match")
	assert.Equal(t, expectedTitle, actualResponse.Title, "title should match")
	assert.Equal(t, expectedContent, actualResponse.Content, "content should match")
}

func TestPostIntegration(t *testing.T) {
	suite := testutils.Setup()
	var createdPostID int

	t.Run("Create Post", func(t *testing.T) {
		createdPostID = createPostTest(t, suite, TestTitle, TestContent)
		require.NotZero(t, createdPostID, "Created post ID should not be zero")
	})

	t.Run("Get All Posts", func(t *testing.T) {
		posts := getPostsTest(t, suite)
		require.NotEmpty(t, posts, "Posts list should not be empty")
	})

	t.Run("Get Post by ID", func(t *testing.T) {
		getPostByIDTest(t, suite, createdPostID, TestTitle, TestContent)
	})

	t.Run("Update Post", func(t *testing.T) {
		updatePostTest(t, suite, createdPostID, UpdatedTitle, UpdatedContent)
		getPostByIDTest(t, suite, createdPostID, UpdatedTitle, UpdatedContent)
	})

	t.Run("Delete Post", func(t *testing.T) {
		deletePostTest(t, suite, createdPostID)
		w := suite.MakeRequest("GET", fmt.Sprintf("/api/posts/%d", createdPostID), nil)
		if w.Code != http.StatusNotFound {
			t.Fatalf("should get not found after delete, got %d", w.Code)
		}
	})
}

func TestCreatePost(t *testing.T) {
	suite := testutils.Setup()
	createPostTest(t, suite, TestTitle, TestContent)
}

func TestGetPosts(t *testing.T) {
	suite := testutils.Setup()
	getPostsTest(t, suite)
}

func TestUpdatePost(t *testing.T) {
	suite := testutils.Setup()
	createdID := createPostTest(t, suite, TestTitle, TestContent)
	updatePostTest(t, suite, createdID, UpdatedTitle, UpdatedContent)
	getPostByIDTest(t, suite, createdID, UpdatedTitle, UpdatedContent)
}

func TestDeletePost(t *testing.T) {
	suite := testutils.Setup()
	createdID := createPostTest(t, suite, TestTitle, TestContent)
	deletePostTest(t, suite, createdID)
	w := suite.MakeRequest("GET", fmt.Sprintf("/api/posts/%d", createdID), nil)
	if w.Code != http.StatusNotFound {
		t.Fatalf("should get not found after delete, got %d", w.Code)
	}
}

func TestGetPostByID(t *testing.T) {
	suite := testutils.Setup()
	newID := createPostTest(t, suite, TestTitle, TestContent)
	getPostByIDTest(t, suite, newID, TestTitle, TestContent)
}
