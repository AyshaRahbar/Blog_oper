package models

type Post struct {
	ID      int    `json:"id" gorm:"primaryKey"`
	Title   string `json:"title"`
	Content string `json:"content"`
}

type UpdatePostRequest struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}
