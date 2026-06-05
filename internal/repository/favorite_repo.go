package repository

import (
	"context"
	"database/sql"
	"ozinse-backend/internal/model"
)

type FavoriteRepository interface {
	GetFavorites(ctx context.Context, userID int) ([]model.Project, error)
	AddFavorite(ctx context.Context, userID int, projectID int, projectType string) error
	DeleteFavorite(ctx context.Context, userID int, projectID int, projectType string) error
}

type favoriteRepository struct {
	db *sql.DB
}

func NewFavoriteRepository(db *sql.DB) FavoriteRepository {
	return &favoriteRepository{db: db}
}

func (r *favoriteRepository) GetFavorites(ctx context.Context, userID int) ([]model.Project, error) {
	var projects []model.Project

	query := `
		SELECT p.id, p.type, p.title, p.description, p.release_year, p.director, p.producer, p.cover_image_url, p.category_id, p.age_rating_id, f.added_at
		FROM (
			SELECT movie_id as project_id, 'movie' as project_type, user_id, added_at FROM favorite_movies
			UNION ALL
			SELECT series_id as project_id, 'series' as project_type, user_id, added_at FROM favorite_series
		) f
		JOIN (
			SELECT id, 'movie' as type, title, description, release_year, director, producer, cover_image_url, category_id, age_rating_id FROM movies
			UNION ALL
			SELECT id, 'series' as type, title, description, release_year, director, producer, cover_image_url, category_id, age_rating_id FROM series
		) p ON f.project_id = p.id AND f.project_type = p.type
		WHERE f.user_id = $1
		ORDER BY f.added_at DESC;
	`

	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var p model.Project
		var addedAt string
		err := rows.Scan(
			&p.ID, &p.Type, &p.Title, &p.Description, &p.ReleaseYear,
			&p.Director, &p.Producer, &p.CoverImageURL, &p.CategoryID, &p.AgeRatingID, &addedAt,
		)
		if err != nil {
			return nil, err
		}
		projects = append(projects, p)
	}

	return projects, nil
}

func (r *favoriteRepository) AddFavorite(ctx context.Context, userID int, projectID int, projectType string) error {
	var query string
	if projectType == "movie" {
		query = `INSERT INTO favorite_movies (user_id, movie_id) VALUES ($1, $2) ON CONFLICT (user_id, movie_id) DO NOTHING;`
	} else {
		query = `INSERT INTO favorite_series (user_id, series_id) VALUES ($1, $2) ON CONFLICT (user_id, series_id) DO NOTHING;`
	}
	
	_, err := r.db.ExecContext(ctx, query, userID, projectID)
	return err
}

func (r *favoriteRepository) DeleteFavorite(ctx context.Context, userID int, projectID int, projectType string) error {
	var query string
	if projectType == "movie" {
		query = `DELETE FROM favorite_movies WHERE user_id = $1 AND movie_id = $2;`
	} else {
		query = `DELETE FROM favorite_series WHERE user_id = $1 AND series_id = $2;`
	}

	_, err := r.db.ExecContext(ctx, query, userID, projectID)
	return err
}