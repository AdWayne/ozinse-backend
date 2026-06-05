package model

import "time"

type User struct {
	ID           int       `json:"id"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"-"` // Скрываем из JSON
	FullName     string    `json:"full_name"`
	Phone        string    `json:"phone,omitempty"`
	BirthDate    string    `json:"birth_date,omitempty"` // Формат YYYY-MM-DD
	RoleName     string    `json:"role"`
	CreatedAt    time.Time `json:"created_at"`
}

type RegisterInput struct {
	Email          string `json:"email" binding:"required,email"`
	Password       string `json:"password" binding:"required,min=6"`
	RepeatPassword string `json:"repeat_password" binding:"required"`
}

type LoginInput struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type RefreshInput struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

type UpdateProfileInput struct {
	FullName  string `json:"full_name" binding:"required"`
	Phone     string `json:"phone"`
	BirthDate string `json:"birth_date"` // Будет проверяться как строка формата даты
}

type UpdatePasswordInput struct {
	OldPassword string `json:"old_password" binding:"required"`
	NewPassword string `json:"new_password" binding:"required,min=6"`
}

type ResetPasswordInput struct {
	Email string `json:"email" binding:"required,email"`
}