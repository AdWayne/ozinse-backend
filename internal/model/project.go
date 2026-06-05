package model

import "time"

type Project struct {
	ID            int       `json:"id"`
	Type          string    `json:"type"`
	Title         string    `json:"title"`
	Description   string    `json:"description"`
	ReleaseYear   int       `json:"release_year"`
	Director      string    `json:"director"`
	Producer      string    `json:"producer"`
	CoverImageURL string    `json:"cover_image_url"`
	CategoryID    *int      `json:"category_id"`
	AgeRatingID   *int      `json:"age_rating_id"`
	YouTubeVideoID string   `json:"youtube_video_id,omitempty"`
	CreatedAt     time.Time `json:"created_at"`
}

type ProjectFilter struct {
	Search      string `form:"search"`
	CategoryID  int    `form:"category_id"`
	GenreID     int    `form:"genre_id"`
	AgeRatingID int    `form:"age_rating_id"`
	Page        int    `form:"page,default=1"`
	Size        int    `form:"size,default=10"`
}

type ProjectDetail struct {
	Project
	Actors       string   `json:"actors"`
	Runtime      int      `json:"runtime,omitempty"`
	CategoryName string   `json:"category_name"`
	AgeRating    string   `json:"age_rating"`
	Genres       []string `json:"genres"`
}

type Season struct {
	ID        int       `json:"id"`
	ProjectID int       `json:"project_id"`
	Number    int       `json:"number"`
	Title     string    `json:"title"`
	Episodes  []Episode `json:"episodes"`
}

type Episode struct {
	ID             int    `json:"id"`
	SeasonID       int    `json:"season_id"`
	Number         int    `json:"number"`
	Title          string `json:"title"`
	Runtime        int    `json:"runtime"`
	YouTubeVideoID string `json:"youtube_video_id"`
}