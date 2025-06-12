package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"go-blog/testutils"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestLikeFunctionality(t *testing.T) {
	suite := testutils.Setup()

	t.Run("LikePost", func(t *testing.T) {
		timestamp := time.Now().UnixNano()
		bloggerToken := registerUserWithTimestamp(t, suite, fmt.Sprintf("blogger_like_%d", timestamp), "password", "blogger")
		viewerToken := registerUserWithTimestamp(t, suite, fmt.Sprintf("viewer_like_%d", timestamp), "password", "viewer")
		postID := createPostWithToken(t, suite, bloggerToken, "Test Post", "Test Content")
		w := suite.MakeRequest("POST", fmt.Sprintf("/api/posts/%d/like", postID), nil, map[string]string{"Authorization": "Bearer " + viewerToken})

		assert.Equal(t, http.StatusCreated, w.Code)
		var response map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &response)
		assert.Equal(t, "Post liked successfully", response["message"])
		assert.NotNil(t, response["like"])
	})

	t.Run("UnlikePost", func(t *testing.T) {
		timestamp := time.Now().UnixNano()
		bloggerToken := registerUserWithTimestamp(t, suite, fmt.Sprintf("blogger_unlike_%d", timestamp), "password", "blogger")
		viewerToken := registerUserWithTimestamp(t, suite, fmt.Sprintf("viewer_unlike_%d", timestamp), "password", "viewer")

		postID := createPostWithToken(t, suite, bloggerToken, "Test Post", "Test Content")

		suite.MakeRequest("POST", fmt.Sprintf("/api/posts/%d/like", postID), nil, map[string]string{"Authorization": "Bearer " + viewerToken})
		w := suite.MakeRequest("DELETE", fmt.Sprintf("/api/posts/%d/like", postID), nil, map[string]string{"Authorization": "Bearer " + viewerToken})

		assert.Equal(t, http.StatusOK, w.Code)
		var response map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &response)
		assert.Equal(t, "Post unliked successfully", response["message"])
	})

	t.Run("GetPostLikes", func(t *testing.T) {
		timestamp := time.Now().UnixNano()
		bloggerToken := registerUserWithTimestamp(t, suite, fmt.Sprintf("blogger_getlikes_%d", timestamp), "password", "blogger")
		viewerToken := registerUserWithTimestamp(t, suite, fmt.Sprintf("viewer_getlikes_%d", timestamp), "password", "viewer")
		postID := createPostWithToken(t, suite, bloggerToken, "Test Post", "Test Content")
		suite.MakeRequest("POST", fmt.Sprintf("/api/posts/%d/like", postID), nil, map[string]string{"Authorization": "Bearer " + viewerToken})

		w := suite.MakeRequest("GET", fmt.Sprintf("/api/posts/%d/likes", postID), nil, map[string]string{"Authorization": "Bearer " + viewerToken})
		assert.Equal(t, http.StatusOK, w.Code)

		type GetLikesResponse struct {
			PostID        int  `json:"post_id"`
			LikeCount     int  `json:"like_count"`
			IsLikedByUser bool `json:"is_liked_by_user"`
		}

		var response GetLikesResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, postID, response.PostID)
		assert.GreaterOrEqual(t, response.LikeCount, 1)
	})

	t.Run("GetUserLikes", func(t *testing.T) {
		timestamp := time.Now().UnixNano()
		bloggerToken := registerUserWithTimestamp(t, suite, fmt.Sprintf("blogger_userlikes_%d", timestamp), "password", "blogger")
		viewerToken := registerUserWithTimestamp(t, suite, fmt.Sprintf("viewer_userlikes_%d", timestamp), "password", "viewer")
		postID := createPostWithToken(t, suite, bloggerToken, "Test Post", "Test Content")

		suite.MakeRequest("POST", fmt.Sprintf("/api/posts/%d/like", postID), nil, map[string]string{"Authorization": "Bearer " + viewerToken})
		w := suite.MakeRequest("GET", "/api/users/me/likes", nil, map[string]string{"Authorization": "Bearer " + viewerToken})

		assert.Equal(t, http.StatusOK, w.Code)
		var response map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &response)
		likes := response["likes"].([]interface{})
		assert.GreaterOrEqual(t, len(likes), 1) // At least 1 like

		found := false
		for _, like := range likes {
			likeMap := like.(map[string]interface{})
			if int(likeMap["post_id"].(float64)) == postID {
				found = true
				break
			}
		}
		assert.True(t, found, "Should find the like we created")
	})

	t.Run("CannotLikeOwnPost", func(t *testing.T) {
		timestamp := time.Now().UnixNano()
		bloggerToken := registerUserWithTimestamp(t, suite, fmt.Sprintf("blogger_ownlike_%d", timestamp), "password", "blogger")
		postID := createPostWithToken(t, suite, bloggerToken, "Test Post", "Test Content")
		w := suite.MakeRequest("POST", fmt.Sprintf("/api/posts/%d/like", postID), nil, map[string]string{"Authorization": "Bearer " + bloggerToken})

		assert.Equal(t, http.StatusForbidden, w.Code)
		var response map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &response)
		assert.Equal(t, "Cannot like your own post", response["error"])
	})

	t.Run("AlreadyLiked", func(t *testing.T) {
		timestamp := time.Now().UnixNano()
		bloggerToken := registerUserWithTimestamp(t, suite, fmt.Sprintf("blogger_already_%d", timestamp), "password", "blogger")
		viewerToken := registerUserWithTimestamp(t, suite, fmt.Sprintf("viewer_already_%d", timestamp), "password", "viewer")

		postID := createPostWithToken(t, suite, bloggerToken, "Test Post", "Test Content")
		suite.MakeRequest("POST", fmt.Sprintf("/api/posts/%d/like", postID), nil, map[string]string{"Authorization": "Bearer " + viewerToken})
		w := suite.MakeRequest("POST", fmt.Sprintf("/api/posts/%d/like", postID), nil, map[string]string{"Authorization": "Bearer " + viewerToken})

		assert.Equal(t, http.StatusConflict, w.Code)
		var response map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &response)
		assert.Equal(t, "You have already liked this post", response["error"])
	})
}

func registerUserWithTimestamp(t *testing.T, suite *testutils.TestSuite, username, password, accountType string) string {
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

func createPostWithToken(t *testing.T, suite *testutils.TestSuite, token string, title, content string) int {
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
