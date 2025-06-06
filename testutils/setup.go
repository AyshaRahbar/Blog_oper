package testutils

import (
	"bytes"
	"fmt"
	"go-blog/handlers"
	"go-blog/models"
	"go-blog/repo"
	"go-blog/routes"
	"go-blog/service"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"net/http"
	"net/http/httptest"
)

type TestSuite struct {
	Router      *gin.Engine
	PostService service.PostService
	PostRepo    repo.PostRepository
	UserHandler *handlers.UserHandler
}

func Setup() *TestSuite {
	dsn := "postgres://postgres:postgres@localhost:5432/blogdb?sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(fmt.Sprintf("couldn't connect to db: %v", err))
	}

	db.AutoMigrate(&models.Post{}, &models.User{})
	gin.SetMode(gin.TestMode)

	postRepository := repo.NewPostRepository(db)
	postService := service.NewPostService(postRepository)
	postHandler := handlers.NewPostHandler(postService)

	userRepository := repo.NewUserRepository(db)
	userService := service.NewUserService(userRepository)
	userHandler := handlers.NewUserHandler(userService)

	router := routes.SetupRoutes(postHandler, userHandler)

	return &TestSuite{
		Router:      router,
		PostService: postService,
		PostRepo:    postRepository,
		UserHandler: userHandler,
	}
}

func (s *TestSuite) MakeRequest(method, url string, body *bytes.Buffer, headers ...map[string]string) *httptest.ResponseRecorder {
	var req *http.Request
	if body != nil {
		req, _ = http.NewRequest(method, url, body)
		req.Header.Set("Content-Type", "application/json")
	} else {
		req, _ = http.NewRequest(method, url, nil)
	}
	if len(headers) > 0 {
		for k, v := range headers[0] {
			req.Header.Set(k, v)
		}
	}
	w := httptest.NewRecorder()
	s.Router.ServeHTTP(w, req)
	return w
}
