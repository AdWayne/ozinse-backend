package repository

import (
	"context"
	"database/sql"
	"fmt"
	"ozinse-backend/internal/model"
)

type ProjectRepository interface {
	GetProjects(ctx context.Context, filter model.ProjectFilter) ([]model.Project, error)
	GetProjectTypeAndTable(ctx context.Context, id int) (string, string, error)
	GetProjectByID(ctx context.Context, id int, tableName string) (*model.ProjectDetail, error)
	GetSeasonsWithEpisodes(ctx context.Context, seriesID int) ([]model.Season, error)
}

type projectRepository struct {
	db *sql.DB
}

func NewProjectRepository(db *sql.DB) ProjectRepository {
	return &projectRepository{db: db}
}

func (r *projectRepository) GetProjects(ctx context.Context, filter model.ProjectFilter) ([]model.Project, error) {
	var projects []model.Project
	var args []interface{}
	argCounter := 1

	query := `
		SELECT id, type, title, description, release_year, director, producer, cover_image_url, category_id, age_rating_id, created_at
		FROM (
			SELECT id, 'movie' as type, title, description, release_year, director, producer, cover_image_url, category_id, age_rating_id, created_at,
			       ARRAY(SELECT genre_id FROM movie_genres mg WHERE mg.movie_id = movies.id) as genre_ids
			FROM movies
			UNION ALL
			SELECT id, 'series' as type, title, description, release_year, director, producer, cover_image_url, category_id, age_rating_id, created_at,
			       ARRAY(SELECT genre_id FROM series_genres sg WHERE sg.series_id = series.id) as genre_ids
			FROM series
		) as combined
		WHERE 1=1
	`

	if filter.Search != "" {
		query += fmt.Sprintf(" AND (LOWER(title) LIKE LOWER($%d) OR LOWER(description) LIKE LOWER($%d))", argCounter, argCounter)
		args = append(args, "%"+filter.Search+"%")
		argCounter++
	}

	if filter.CategoryID > 0 {
		query += fmt.Sprintf(" AND category_id = $%d", argCounter)
		args = append(args, filter.CategoryID)
		argCounter++
	}

	if filter.AgeRatingID > 0 {
		query += fmt.Sprintf(" AND age_rating_id = $%d", argCounter)
		args = append(args, filter.AgeRatingID)
		argCounter++
	}

	if filter.GenreID > 0 {
		query += fmt.Sprintf(" AND $%d = ANY(genre_ids)", argCounter)
		args = append(args, filter.GenreID)
		argCounter++
	}

	offset := (filter.Page - 1) * filter.Size
	query += fmt.Sprintf(" ORDER BY created_at DESC LIMIT $%d OFFSET $%d;", argCounter, argCounter+1)
	args = append(args, filter.Size, offset)

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var p model.Project
		err := rows.Scan(
			&p.ID, &p.Type, &p.Title, &p.Description, &p.ReleaseYear,
			&p.Director, &p.Producer, &p.CoverImageURL, &p.CategoryID, &p.AgeRatingID, &p.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		projects = append(projects, p)
	}

	return projects, nil
}

func (r *projectRepository) GetProjectTypeAndTable(ctx context.Context, id int) (string, string, error) {
	var pType string
	err := r.db.QueryRowContext(ctx, "SELECT 'movie' FROM movies WHERE id = $1", id).Scan(&pType)
	if err == nil {
		return pType, "movies", nil
	}
	err = r.db.QueryRowContext(ctx, "SELECT 'series' FROM series WHERE id = $1", id).Scan(&pType)
	if err == nil {
		return pType, "series", nil
	}
	return "", "", fmt.Errorf("PROJECT_NOT_FOUND")
}

func (r *projectRepository) GetProjectByID(ctx context.Context, id int, tableName string) (*model.ProjectDetail, error) {
	var detail model.ProjectDetail
	detail.Type = "movie"
	if tableName == "series" {
		detail.Type = "series"
	}

	// Учитываем колонку range вместо name в age_ratings
	query := fmt.Sprintf(`
		SELECT p.id, p.title, p.description, p.release_year, p.director, p.producer, p.cover_image_url, p.category_id, p.age_rating_id, p.created_at,
		       c.name as category_name, COALESCE(ar.range, '') as age_rating
		FROM %s p
		LEFT JOIN categories c ON p.category_id = c.id
		LEFT JOIN age_ratings ar ON p.age_rating_id = ar.id
		WHERE p.id = $1;`, tableName)

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&detail.ID, &detail.Title, &detail.Description, &detail.ReleaseYear,
		&detail.Director, &detail.Producer, &detail.CoverImageURL, &detail.CategoryID, &detail.AgeRatingID, &detail.CreatedAt,
		&detail.CategoryName, &detail.AgeRating,
	)
	if err != nil {
		return nil, err
	}

	// Достаем жанры
	genreQuery := fmt.Sprintf(`
		SELECT g.name FROM genres g
		JOIN %s_genres pg ON g.id = pg.genre_id
		WHERE pg.%s_id = $1;`, detail.Type, detail.Type)

	rows, err := r.db.QueryContext(ctx, genreQuery, id)
	if err == nil {
		defer rows.Close()
		for rows.Next() {
			var gName string
			if err := rows.Scan(&gName); err == nil {
				detail.Genres = append(detail.Genres, gName)
			}
		}
	}

	return &detail, nil
}

func (r *projectRepository) GetSeasonsWithEpisodes(ctx context.Context, seriesID int) ([]model.Season, error) {
	var seasons []model.Season

	// Учитываем season_number
	sRows, err := r.db.QueryContext(ctx, "SELECT id, series_id, season_number FROM seasons WHERE series_id = $1 ORDER BY season_number ASC;", seriesID)
	if err != nil {
		return nil, err
	}
	defer sRows.Close()

	for sRows.Next() {
		var s model.Season
		// Если в твоей модели структуры другие имена полей, проследи, чтобы они мапились сюда
		if err := sRows.Scan(&s.ID, &s.ProjectID, &s.Number); err != nil {
			return nil, err
		}
		s.Title = fmt.Sprintf("%d сезон", s.Number)

		// Учитываем episode_number и duration
		eRows, err := r.db.QueryContext(ctx, "SELECT id, season_id, episode_number, title, duration, youtube_video_id FROM episodes WHERE season_id = $1 ORDER BY episode_number ASC;", s.ID)
		if err == nil {
			for eRows.Next() {
				var e model.Episode
				if err := eRows.Scan(&e.ID, &e.SeasonID, &e.Number, &e.Title, &e.Runtime, &e.YouTubeVideoID); err == nil {
					s.Episodes = append(s.Episodes, e)
				}
			}
			eRows.Close()
		}

		seasons = append(seasons, s)
	}

	return seasons, nil
}