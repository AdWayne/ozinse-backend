package model

import "time"

type ProjectAdmin struct {
	ID             int       `json:"id"`
	Type           string    `json:"type" binding:"required"`
	Title          string    `json:"title" binding:"required"`
	Description    string    `json:"description"`
	ReleaseYear    int       `json:"release_year"`
	Director       string    `json:"director"`
	Producer       string    `json:"producer"`
	CoverImageURL  string    `json:"cover_image_url"`
	CategoryID     int       `json:"category_id"`
	AgeRatingID    int       `json:"age_rating_id"`
	YoutubeVideoID string    `json:"youtube_video_id,omitempty"`
	CreatedAt      time.Time `json:"created_at"`
}

type SeasonAdmin struct {
	ID           int `json:"id"`
	SeriesID     int `json:"series_id"`
	SeasonNumber int `json:"season_number" binding:"required"`
}

type EpisodeAdmin struct {
	ID             int    `json:"id"`
	SeasonID       int    `json:"season_id"`
	EpisodeNumber  int    `json:"episode_number" binding:"required"`
	Title          string `json:"title"`
	YoutubeVideoID string `json:"youtube_video_id" binding:"required"`
	Duration       int    `json:"duration"`
}

type FeaturedOrderItem struct {
	ID   int    `json:"id" binding:"required"`
	Type string `json:"type" binding:"required"`
}

type FeaturedOrderRequest struct {
	BlockType string              `json:"block_type" binding:"required"`
	Items     []FeaturedOrderItem `json:"items" binding:"required"`
}

type ProjectsListResponse struct {
	Items      []ProjectAdmin `json:"items"`
	TotalCount int            `json:"total_count"`
	Page       int            `json:"page"`
	Limit      int            `json:"limit"`
}