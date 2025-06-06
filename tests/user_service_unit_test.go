package tests

import (
	"bytes"
	"encoding/json"
	"go-blog/testutils"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRegisterUser(t *testing.T) {
	suite := testutils.Setup()
	router := suite.Router
	payload := map[string]interface{}{
		"username":     "testuser1",
		"password":     "testpass1",
		"account_type": "blogger",
	}
	body, _ := json.Marshal(payload)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/register", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Equal(t, "user created", resp["message"])
	user := resp["user"]
	assert.NotContains(t, user, "password")
}

func TestLoginUser(t *testing.T) {
	suite := testutils.Setup()
	router := suite.Router
	payload := map[string]interface{}{
		"username":     "loginuser",
		"password":     "loginpass",
		"account_type": "viewer",
	}
	body, _ := json.Marshal(payload)
	req, _ := http.NewRequest("POST", "/api/register", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	loginPayload := map[string]interface{}{
		"username": "loginuser",
		"password": "loginpass",
	}
	loginBody, _ := json.Marshal(loginPayload)
	loginReq, _ := http.NewRequest("POST", "/api/login", bytes.NewBuffer(loginBody))
	loginReq.Header.Set("Content-Type", "application/json")
	loginW := httptest.NewRecorder()
	router.ServeHTTP(loginW, loginReq)
	assert.Equal(t, http.StatusOK, loginW.Code)
	var resp map[string]interface{}
	json.Unmarshal(loginW.Body.Bytes(), &resp)
	assert.NotEmpty(t, resp["token"])

	wrongLoginPayload := map[string]interface{}{
		"username": "loginuser",
		"password": "wrongpass",
	}
	wrongLoginBody, _ := json.Marshal(wrongLoginPayload)
	wrongLoginReq, _ := http.NewRequest("POST", "/api/login", bytes.NewBuffer(wrongLoginBody))
	wrongLoginReq.Header.Set("Content-Type", "application/json")
	wrongLoginW := httptest.NewRecorder()
	router.ServeHTTP(wrongLoginW, wrongLoginReq)
	assert.Equal(t, http.StatusBadRequest, wrongLoginW.Code)

	nonExistLoginPayload := map[string]interface{}{
		"username": "nouser",
		"password": "nopass",
	}
	nonExistLoginBody, _ := json.Marshal(nonExistLoginPayload)
	nonExistLoginReq, _ := http.NewRequest("POST", "/api/login", bytes.NewBuffer(nonExistLoginBody))
	nonExistLoginReq.Header.Set("Content-Type", "application/json")
	nonExistLoginW := httptest.NewRecorder()
	router.ServeHTTP(nonExistLoginW, nonExistLoginReq)
	assert.Equal(t, http.StatusBadRequest, nonExistLoginW.Code)
}
