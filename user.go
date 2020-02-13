package app

import (
	"net/http"
	"time"
)

type User struct {
	ID        uint32    `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	Password  string    `json:"password"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type UserService interface {
	CreateToken(userId uint32) (string, error)
	ExtractAuthenticationToken(r *http.Request) (uint32, error)
	Save(user *User) error
	GetById(userId uint32) (*User, error)
	Login(u *User) (string, error)
}
