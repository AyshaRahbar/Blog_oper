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

func initPostgreSQL() (*gorm.DB, error) {
	dsn := os.Getenv("DB_DSN")
	if dsn == "" {
		dsn = "host=localhost user=postgres password=postgres dbname=blogdb port=5432 sslmode=disable"
	}
	return gorm.Open(postgres.Open(dsn), &gorm.Config{})
}

func main() {
	db, _ := initPostgreSQL()
	db.AutoMigrate(&models.Post{})

	postRepo := repo.NewPostRepository(db)
	postService := service.NewPostService(postRepo)
	postHandler := handlers.NewPostHandler(postService)

	r := routes.SetupRoutes(postHandler)
	r.Run(":8080")
}
