package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
	"os"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"github.com/joho/godotenv"
	"log"
)

func init() {
    err := godotenv.Load()
    if err != nil {
        log.Println("No .env file found, using system environment")
    }
}

type Post struct {
	ID        int       
	Title     string    
	Content   string    
	Created   time.Time 
	Updated   time.Time 	
	Reactions []string
}

type Emoji struct {
	Slug    string 
	Unicode string 
}

var db *gorm.DB

func initDB() {
	var err error
	
	dbHost := ("DB_HOST")
	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASS")
	dbName := os.Getenv("DB_NAME")
	dbPort := os.Getenv("DB_PORT")
	dbSSLMode := os.Getenv("DB_SSLMODE")
	
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		dbHost, dbUser, dbPass, dbName, dbPort, dbSSLMode)
	
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		fmt.Println("Failed to connect to the database:", err)
		return
	}
	db.AutoMigrate(&Post{})
}

func createPost(title, content string) {
	post := Post{
		Title:     title,
		Content:   content,
		Created:   time.Now(),
		Reactions: []string{},
	}
	db.Create(&post)
	fmt.Println("New post created")
}

func ListPost() {
	var posts []Post
	if err := db.Find(&posts).Error; err != nil {
		fmt.Println("Error fetching posts:", err)
		return
	}
	if len(posts) == 0 {
		fmt.Println("No posts created yet.")
	} else {
		fmt.Println("The posts are:")
		for i, post := range posts {
			fmt.Printf("%d. %s - %s (Created on: %s)\n", i+1, post.Title, post.Content, post.Created.Format("January 2, 2006"))
			if len(post.Reactions) > 0 {
				fmt.Println("Reactions:")
				for _, emoji := range post.Reactions {
					fmt.Printf(" - %s\n", emoji)
				}
			} else {
				fmt.Println("No reactions yet.")
			}
		}
	}
}

func addReactions(postID int, emoji string) {
	var post Post
	if err := db.First(&post, postID).Error; err != nil {
		fmt.Println("Post not found.")
		return
	}
	post.Reactions = append(post.Reactions, emoji)
	if err := db.Save(&post).Error; err != nil {
		fmt.Println("Error updating post reactions:", err)
	} else {
		fmt.Println("Reaction added:", emoji)
	}
}

func fetchEmojis() ([]Emoji, error) {
	apiURL := os.Getenv("API_URL")
	resp, err := http.Get(apiURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var emojis []Emoji
	if err := json.Unmarshal(body, &emojis); err != nil {
		return nil, err
	}
	return emojis, nil
}

func main() {
	initDB()

	for {
		fmt.Println("This is a blog post")
		fmt.Println("1. Create new post")
		fmt.Println("2. View posts")
		fmt.Println("3. React to a post")

		var choice int
		fmt.Print("Choose an option: ")
		_, err := fmt.Scan(&choice)

		if err != nil {
			fmt.Println("You entered the wrong choice.")
			continue
		} else {
			fmt.Println("You have entered:", choice)
		}

		switch choice {
		case 1:
			var title, content string
			fmt.Println("Enter post title:")
			fmt.Scan(&title)
			fmt.Println("Enter post content:")
			fmt.Scan(&content)
			createPost(title, content)
		case 2:
			ListPost()
		case 3:
			var postID int
			fmt.Println("Enter the post ID to react to:")
			fmt.Scan(&postID)

			emojis, err := fetchEmojis()
			if err != nil {
				fmt.Println("Error fetching emojis:", err)
				break
			}

			if len(emojis) > 0 {
				emoji := emojis[0].Unicode
				addReactions(postID, emoji)
			} else {
				fmt.Println("No emojis available.")
			}
		case 4:
			var postID int
			fmt.Println("Enter the post ID to delete:")
			fmt.Scan(&postID)
			deletePost(postID)
		default:
			fmt.Println("Please enter a valid choice.")
		}
	}
}
func deletePost(postID int) {
	var post Post
	if err := db.First(&post, postID).Error; err != nil {
		fmt.Println("Post not found.")
		return
	}

	fmt.Printf("Are you sure you want to delete '%s'? (y/n): ", post.Title)
	var confirm string
	fmt.Scan(&confirm)

	if confirm == "y" || confirm == "Y" {
		if err := db.Delete(&post).Error; err != nil {
			fmt.Println("Error deleting post:", err)
		} else {
			fmt.Println("Post deleted successfully")
		}
	} else {
		fmt.Println("Delete cancelled")
	}
}
