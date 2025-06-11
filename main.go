package main

import (
	"go-blog/handlers"
	"go-blog/models"
	"go-blog/repo"
	"go-blog/routes"
	"go-blog/service"
	"log"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}
}

func initPostgreSQL() *gorm.DB {
	dsn := os.Getenv("DB_DSN")
	if dsn == "" {
		dsn = "host=localhost user=postgres password=postgres dbname=blogdb port=5432 sslmode=disable"
	}
	db, _ := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	return db
}

func main() {
	db := initPostgreSQL()
	db.AutoMigrate(&models.Post{}, &models.User{}, &models.Like{}, &models.Comment{})

	postRepo := repo.NewPostRepository(db)
	authRepo := repo.NewAuthRepository(postRepo)
	postService := service.NewPostService(postRepo)
	postHandler := handlers.NewPostHandler(postService)

	userRepo := repo.NewUserRepository(db)
	userService := service.NewUserService(userRepo)
	userHandler := handlers.NewUserHandler(userService)

	likeRepo := repo.NewLikeRepository(db)
	likeService := service.NewLikeService(likeRepo, postRepo)
	likeHandler := handlers.NewLikeHandler(likeService)

	commentRepo := repo.NewCommentRepository(db)
	commentService := service.NewCommentService(commentRepo, postRepo)
	commentHandler := handlers.NewCommentHandler(commentService)

	r := routes.SetupRoutes(postHandler, userHandler, likeHandler, commentHandler, authRepo)
	r.Run(":8080")
}
