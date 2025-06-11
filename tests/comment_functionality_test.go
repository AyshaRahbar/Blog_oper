package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"go-blog/testutils"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCommentFunctionality(t *testing.T) {
	suite := testutils.Setup()

	t.Run("CreateComment", func(t *testing.T) {
		bloggerToken := registerUser(t, suite, "blogger_comment", "password", "blogger")
		viewerToken := registerUser(t, suite, "viewer_comment", "password", "viewer")

		postID := createPost(t, suite, bloggerToken, "Test Post", "Test Content")

		commentPayload := map[string]interface{}{
			"comment": "This is a test comment",
		}
		commentBody, _ := json.Marshal(commentPayload)
		w := suite.MakeRequest("POST", fmt.Sprintf("/api/posts/%d/comments", postID), bytes.NewBuffer(commentBody), map[string]string{"Authorization": "Bearer " + viewerToken})

		assert.Equal(t, http.StatusCreated, w.Code)
		var response map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &response)
		assert.Equal(t, "Comment created successfully", response["message"])
		assert.NotNil(t, response["comment"])
	})

	t.Run("GetPostComments", func(t *testing.T) {
		bloggerToken := registerUser(t, suite, "blogger_get", "password", "blogger")
		viewerToken := registerUser(t, suite, "viewer_get", "password", "viewer")

		postID := createPost(t, suite, bloggerToken, "Test Post", "Test Content")

		commentPayload := map[string]interface{}{
			"comment": "Comment to retrieve",
		}
		commentBody, _ := json.Marshal(commentPayload)
		suite.MakeRequest("POST", fmt.Sprintf("/api/posts/%d/comments", postID), bytes.NewBuffer(commentBody), map[string]string{"Authorization": "Bearer " + viewerToken})

		w := suite.MakeRequest("GET", fmt.Sprintf("/api/posts/%d/comments", postID), nil)

		assert.Equal(t, http.StatusOK, w.Code)
		var response map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &response)
		comments := response["comments"].([]interface{})
		assert.GreaterOrEqual(t, len(comments), 1)
	})

	t.Run("UpdateComment", func(t *testing.T) {
		bloggerToken := registerUser(t, suite, "blogger_update", "password", "blogger")
		viewerToken := registerUser(t, suite, "viewer_update", "password", "viewer")

		postID := createPost(t, suite, bloggerToken, "Test Post", "Test Content")

		commentPayload := map[string]interface{}{
			"comment": "Original comment",
		}
		commentBody, _ := json.Marshal(commentPayload)
		createW := suite.MakeRequest("POST", fmt.Sprintf("/api/posts/%d/comments", postID), bytes.NewBuffer(commentBody), map[string]string{"Authorization": "Bearer " + viewerToken})

		var createResponse map[string]interface{}
		json.Unmarshal(createW.Body.Bytes(), &createResponse)
		comment := createResponse["comment"].(map[string]interface{})
		commentID := int(comment["id"].(float64))

		updatePayload := map[string]interface{}{
			"comment": "Updated comment",
		}
		updateBody, _ := json.Marshal(updatePayload)
		w := suite.MakeRequest("PUT", fmt.Sprintf("/api/comments/%d", commentID), bytes.NewBuffer(updateBody), map[string]string{"Authorization": "Bearer " + viewerToken})

		assert.Equal(t, http.StatusOK, w.Code)
		var response map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &response)
		assert.Equal(t, "Comment updated successfully", response["message"])
	})

	t.Run("DeleteComment", func(t *testing.T) {
		bloggerToken := registerUser(t, suite, "blogger_delete", "password", "blogger")
		viewerToken := registerUser(t, suite, "viewer_delete", "password", "viewer")

		postID := createPost(t, suite, bloggerToken, "Test Post", "Test Content")

		commentPayload := map[string]interface{}{
			"comment": "Comment to delete",
		}
		commentBody, _ := json.Marshal(commentPayload)
		createW := suite.MakeRequest("POST", fmt.Sprintf("/api/posts/%d/comments", postID), bytes.NewBuffer(commentBody), map[string]string{"Authorization": "Bearer " + viewerToken})

		var createResponse map[string]interface{}
		json.Unmarshal(createW.Body.Bytes(), &createResponse)
		comment := createResponse["comment"].(map[string]interface{})
		commentID := int(comment["id"].(float64))

		w := suite.MakeRequest("DELETE", fmt.Sprintf("/api/comments/%d", commentID), nil, map[string]string{"Authorization": "Bearer " + viewerToken})

		assert.Equal(t, http.StatusOK, w.Code)
		var response map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &response)
		assert.Equal(t, "Comment deleted successfully", response["message"])
	})
}

func registerUser(t *testing.T, suite *testutils.TestSuite, username, password, accountType string) string {
	payload := map[string]interface{}{
		"username":     username,
		"password":     password,
		"account_type": accountType,
	}
	body, _ := json.Marshal(payload)
	w := suite.MakeRequest("POST", "/api/register", bytes.NewBuffer(body))

	if w.Code != http.StatusCreated {
		loginPayload := map[string]interface{}{
			"username": username,
			"password": password,
		}
		loginBody, _ := json.Marshal(loginPayload)
		w = suite.MakeRequest("POST", "/api/login", bytes.NewBuffer(loginBody))
	} else {
		loginPayload := map[string]interface{}{
			"username": username,
			"password": password,
		}
		loginBody, _ := json.Marshal(loginPayload)
		w = suite.MakeRequest("POST", "/api/login", bytes.NewBuffer(loginBody))
	}

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	return response["token"].(string)
}

func createPost(t *testing.T, suite *testutils.TestSuite, token string, title, content string) int {
	payload := map[string]string{"title": title, "content": content}
	body, _ := json.Marshal(payload)
	w := suite.MakeRequest("POST", "/api/posts", bytes.NewBuffer(body), map[string]string{"Authorization": "Bearer " + token})

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	if postResp, ok := response["post"].(map[string]interface{}); ok {
		return int(postResp["id"].(float64))
	}
	return int(response["id"].(float64))
}
