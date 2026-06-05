package repository

import (
	"context"
	"database/sql"
	"ozinse-backend/internal/model"
)

type ReferenceRepository interface {
	GetCategories(ctx context.Context) ([]model.Category, error)
	GetGenres(ctx context.Context) ([]model.Genre, error)
	GetAgeRatings(ctx context.Context) ([]model.AgeRating, error)
}

type referenceRepository struct {
	db *sql.DB
}

func NewReferenceRepository(db *sql.DB) ReferenceRepository {
	return &referenceRepository{db: db}
}

func (r *referenceRepository) GetCategories(ctx context.Context) ([]model.Category, error) {
	var categories []model.Category
	rows, err := r.db.QueryContext(ctx, "SELECT id, name FROM categories ORDER BY id ASC;")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var c model.Category
		if err := rows.Scan(&c.ID, &c.Name); err != nil {
			return nil, err
		}
		categories = append(categories, c)
	}
	return categories, nil
}

func (r *referenceRepository) GetGenres(ctx context.Context) ([]model.Genre, error) {
	var genres []model.Genre
	rows, err := r.db.QueryContext(ctx, "SELECT id, name, COALESCE(icon_url, '') FROM genres ORDER BY id ASC;")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var g model.Genre
		if err := rows.Scan(&g.ID, &g.Name, &g.IconURL); err != nil {
			return nil, err
		}
		genres = append(genres, g)
	}
	return genres, nil
}

func (r *referenceRepository) GetAgeRatings(ctx context.Context) ([]model.AgeRating, error) {
	var ratings []model.AgeRating
	rows, err := r.db.QueryContext(ctx, "SELECT id, range FROM age_ratings ORDER BY id ASC;")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var ar model.AgeRating
		if err := rows.Scan(&ar.ID, &ar.Range); err != nil {
			return nil, err
		}
		ratings = append(ratings, ar)
	}
	return ratings, nil
}