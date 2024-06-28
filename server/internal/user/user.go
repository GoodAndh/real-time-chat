package user

import (
	"context"
	"time"
)

type User struct {
	ID          int64     `json:"id"`
	Username    string    `json:"username"`
	Email       string    `json:"email"`
	Password    string    `json:"password"`
	CreatedAt   time.Time `json:"created_at"`
	LastUpdated time.Time `json:"last_updated"`
}

type RegisUserRequest struct {
	Username    string    `json:"username" validate:"required,min=8"`
	Email       string    `json:"email" validate:"email,required"`
	Password    string    `json:"password" validate:"required,min=8"`
	CreatedAt   time.Time `json:"created_at"`
	LastUpdated time.Time `json:"last_updated"`
}

type RegisUserResponse struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
}

type LoginUserRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginUserResponse struct {
	accessToken string 
	ID string `json:"id"`
	Username string `json:"username"`
}

type Repository interface {
	CreateUsers(ctx context.Context, user *User) (int, error)
	GetUserByUsername(ctx context.Context, username string) (*User, error)
	GetUserByEmail(ctx context.Context, email string) (*User, error)
}

type Service interface {
	CreateUsers(c context.Context, req *RegisUserRequest) (*RegisUserResponse, error)
	LoginUser(c context.Context, req *LoginUserRequest) (*LoginUserResponse, error)
}
