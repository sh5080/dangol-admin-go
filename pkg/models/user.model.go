package models

// User는 사용자 모델입니다.
type User struct {
	ID    string `json:"id" db:"id"`
	Email string `json:"email" db:"email"`
	Name  string `json:"name" db:"name"`
}

// Role은 사용자 역할입니다.
type Role string

const (
	CUSTOMER Role = "CUSTOMER"
	OWNER    Role = "OWNER"
	ADMIN    Role = "ADMIN"
)
