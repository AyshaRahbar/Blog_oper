package main

import (
	"go-blog/handlers"
	"go-blog/models"
	"go-blog/repo"
	"go-blog/routes"
	"go-blog/service"
	"os"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

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
	db.AutoMigrate(&models.Post{}, &models.User{})

	postRepo := repo.NewPostRepository(db)
	postService := service.NewPostService(postRepo)
	postHandler := handlers.NewPostHandler(postService)
	userRepo := repo.NewUserRepository(db)
	userService := service.NewUserService(userRepo)
	userHandler := handlers.NewUserHandler(userService)

	r := routes.SetupRoutes(postHandler, userHandler)
	r.Run(":8080")
}
