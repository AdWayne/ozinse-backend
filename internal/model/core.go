package model

import (
	"time"
	"encoding/json"
)

type Category struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type AdminProject struct {
	ID          int       `json:"id"`
	Title       string    `json:"title"`
	Type        string    `json:"type"`
	ReleaseYear int       `json:"release_year"`
	CreatedAt   time.Time `json:"created_at"`
}

type Movie struct {
	ID            int       `json:"id"`
	Title         string    `json:"title"`
	Description   string    `json:"description"`
	ReleaseYear   int       `json:"release_year"`
	Director      string    `json:"director"`
	Producer      string    `json:"producer"`
	CoverImageURL string    `json:"cover_image_url"`
	CategoryID    int       `json:"category_id"`
	AgeRatingID   int       `json:"age_rating_id"`
	YoutubeVideoID string   `json:"youtube_video_id"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

// Series отражает таблицу series
type Series struct {
	ID            int       `json:"id"`
	Title         string    `json:"title"`
	Description   string    `json:"description"`
	ReleaseYear   int       `json:"release_year"`
	Director      string    `json:"director"`
	Producer      string    `json:"producer"`
	CoverImageURL string    `json:"cover_image_url"`
	CategoryID    int       `json:"category_id"`
	AgeRatingID   int       `json:"age_rating_id"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

type Genre struct {
	ID      int    `json:"id"`
	Name    string `json:"name"`
	IconURL string `json:"icon_url,omitempty"`
}

type AgeRating struct {
	ID    int    `json:"id"`
	Range string `json:"range"`
}

type UserResponse struct {
	ID        int        `json:"id"`
	Email     string     `json:"email"`
	FullName  *string    `json:"full_name"`
	Phone     *string    `json:"phone"`
	BirthDate *time.Time `json:"birth_date"`
	RoleID    *int       `json:"role_id"`
	CreatedAt time.Time  `json:"created_at"`
}

type Role struct {
	ID          int             `json:"id"`
	Name        string          `json:"name"`
	Permissions json.RawMessage `json:"permissions"`
}

type AssignRoleRequest struct {
	RoleID int `json:"role_id" binding:"required"`
}