package models

type AccountType string

const (
	AccountTypeBlogger AccountType = "blogger"
	AccountTypeViewer  AccountType = "viewer"
)

type User struct {
	ID          int         `json:"id" gorm:"primaryKey"`
	Username    string      `json:"username"`
	Password    string      `json:"password"`
	AccountType AccountType `json:"account_type"`
}

type RegisterRequest struct {
	Username    string      `json:"username"`
	Password    string      `json:"password"`
	AccountType AccountType `json:"account_type"`
}
