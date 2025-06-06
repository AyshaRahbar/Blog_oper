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

func registerAndLogin(t *testing.T, suite *testutils.TestSuite, username, password, accountType string) string {
	registerPayload := map[string]interface{}{
		"username":     username,
		"password":     password,
		"account_type": accountType,
	}
	registerBody, _ := json.Marshal(registerPayload)
	suite.MakeRequest("POST", "/api/register", bytes.NewBuffer(registerBody))
	loginPayload := map[string]interface{}{
		"username": username,
		"password": password,
	}
	loginBody, _ := json.Marshal(loginPayload)
	loginW := suite.MakeRequest("POST", "/api/login", bytes.NewBuffer(loginBody))
	var loginResp map[string]interface{}
	_ = json.Unmarshal(loginW.Body.Bytes(), &loginResp)
	token, ok := loginResp["token"].(string)
	if !ok || token == "" {
		t.Fatalf("expected token in login response")
	}
	return token
}

func createPostAs(t *testing.T, suite *testutils.TestSuite, token string, title, content string) int {
	postPayload := map[string]interface{}{
		"title":   title,
		"content": content,
	}
	postBody, _ := json.Marshal(postPayload)
	createW := suite.MakeRequest("POST", "/api/posts", bytes.NewBuffer(postBody), map[string]string{"Authorization": "Bearer " + token})
	assert.Equal(t, http.StatusCreated, createW.Code)
	var postResp map[string]interface{}
	_ = json.Unmarshal(createW.Body.Bytes(), &postResp)
	postID, ok := postResp["id"].(float64)
	if !ok {
		if postMap, ok2 := postResp["post"].(map[string]interface{}); ok2 {
			postID, ok = postMap["id"].(float64)
		}
	}
	if !ok {
		t.Fatalf("could not get post ID from create response")
	}
	return int(postID)
}

func TestBloggerDeletePost(t *testing.T) {
	suite := testutils.Setup()
	bloggerToken := registerAndLogin(t, suite, "blogdel", "bloggerpass", "blogger")
	postID := createPostAs(t, suite, bloggerToken, "Blogger Delete", "Blogger can delete this post.")
	deleteW := suite.MakeRequest("DELETE", fmt.Sprintf("/api/posts/%d", postID), nil, map[string]string{"Authorization": "Bearer " + bloggerToken})
	assert.Equal(t, http.StatusOK, deleteW.Code)
}

func TestViewerCannotCreateOrUpdatePost(t *testing.T) {
	suite := testutils.Setup()
	viewerToken := registerAndLogin(t, suite, "viewMake", "viewerpass", "viewer")
	postPayload := map[string]interface{}{"title": "Should Not Work", "content": "Viewers cannot create posts."}
	postBody, _ := json.Marshal(postPayload)
	createW := suite.MakeRequest("POST", "/api/posts", bytes.NewBuffer(postBody), map[string]string{"Authorization": "Bearer " + viewerToken})
	assert.Equal(t, http.StatusForbidden, createW.Code)
	updatePayload := map[string]interface{}{"title": "Updated Title", "content": "Updated Content"}
	updateBody, _ := json.Marshal(updatePayload)
	updateW := suite.MakeRequest("PUT", "/api/posts/1/update", bytes.NewBuffer(updateBody), map[string]string{"Authorization": "Bearer " + viewerToken})
	assert.Equal(t, http.StatusForbidden, updateW.Code)
}

func TestViewerCanGetPosts(t *testing.T) {
	suite := testutils.Setup()
	viewerToken := registerAndLogin(t, suite, "viewCanSee", "viewerpass", "viewer")
	getW := suite.MakeRequest("GET", "/api/posts", nil, map[string]string{"Authorization": "Bearer " + viewerToken})
	assert.Equal(t, http.StatusOK, getW.Code)
}

func TestBloggerCanCreateAndUpdatePost(t *testing.T) {
	suite := testutils.Setup()
	bloggerToken := registerAndLogin(t, suite, "blogWriter1", "bloggerpass", "blogger")
	postID := createPostAs(t, suite, bloggerToken, "Blogger Create", "Blogger can create this post.")
	updatePayload := map[string]interface{}{"title": "Updated by Blogger", "content": "Updated Content"}
	updateBody, _ := json.Marshal(updatePayload)
	updateW := suite.MakeRequest("PUT", fmt.Sprintf("/api/posts/%d/update", postID), bytes.NewBuffer(updateBody), map[string]string{"Authorization": "Bearer " + bloggerToken})
	assert.Equal(t, http.StatusOK, updateW.Code)
}

func TestDeleteNonExistentPost(t *testing.T) {
	suite := testutils.Setup()
	bloggerToken := registerAndLogin(t, suite, "blogNoPost", "bloggerpass", "blogger")
	deleteW := suite.MakeRequest("DELETE", "/api/posts/A8921", nil, map[string]string{"Authorization": "Bearer " + bloggerToken})
	assert.True(t, deleteW.Code == http.StatusNotFound || deleteW.Code == http.StatusForbidden)
}

func TestViewerCannotDeletePost(t *testing.T) {
	suite := testutils.Setup()
	viewerToken := registerAndLogin(t, suite, "viewerdelete", "viewerpass", "viewer")
	bloggerToken := registerAndLogin(t, suite, "bloggerforviewdel", "bloggerpass", "blogger")
	postID := createPostAs(t, suite, bloggerToken, "Post to be deleted", "This post will be deleted by viewer.")
	deleteW := suite.MakeRequest("DELETE", fmt.Sprintf("/api/posts/%d", postID), nil, map[string]string{"Authorization": "Bearer " + viewerToken})
	assert.Equal(t, http.StatusForbidden, deleteW.Code)
}

func TestInvalidToken(t *testing.T) {
	suite := testutils.Setup()
	invalidToken := "should.make.no.sense.buffer.value.with.no.meaning."
	getW := suite.MakeRequest("GET", "/api/posts", nil, map[string]string{"Authorization": "Bearer " + invalidToken})
	assert.Equal(t, http.StatusUnauthorized, getW.Code)
}

func TestEmptyToken(t *testing.T) {
	suite := testutils.Setup()
	getW := suite.MakeRequest("GET", "/api/posts", nil, map[string]string{"Authorization": ""})
	assert.Equal(t, http.StatusUnauthorized, getW.Code)
}

func TestInvalidPostID(t *testing.T) {
	suite := testutils.Setup()
	bloggerToken := registerAndLogin(t, suite, "bloggerinvalid", "bloggerpass", "blogger")
	getW := suite.MakeRequest("GET", "/api/posts/abc", nil, map[string]string{"Authorization": "Bearer " + bloggerToken})
	assert.Equal(t, http.StatusBadRequest, getW.Code)
}

func TestEmptyPostContent(t *testing.T) {
	suite := testutils.Setup()
	bloggerToken := registerAndLogin(t, suite, "bloggerempty", "bloggerpass", "blogger")
	postPayload := map[string]interface{}{
		"title":   "null, empty",
		"content": "",
	}
	postBody, _ := json.Marshal(postPayload)
	createW := suite.MakeRequest("POST", "/api/posts", bytes.NewBuffer(postBody), map[string]string{"Authorization": "Bearer " + bloggerToken})
	assert.Equal(t, http.StatusInternalServerError, createW.Code)
}

func TestDuplicateUsername(t *testing.T) {           
	suite := testutils.Setup()
	registerPayload := map[string]interface{}{
		"username":     "duplicateuser",
		"password":     "testpass",
		"account_type": "blogger",
	}
	registerBody, _ := json.Marshal(registerPayload)

	suite.MakeRequest("POST", "/api/register", bytes.NewBuffer(registerBody))

	secondW := suite.MakeRequest("POST", "/api/register", bytes.NewBuffer(registerBody))
	assert.Equal(t, http.StatusBadRequest, secondW.Code)

	var resp map[string]interface{}
	_ = json.Unmarshal(secondW.Body.Bytes(), &resp)
	assert.Equal(t, "username exists", resp["error"])
}

