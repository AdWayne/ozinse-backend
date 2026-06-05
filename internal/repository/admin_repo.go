package repository

import (
	"context"
	"database/sql"
	"fmt"
	"ozinse-backend/internal/model"
)

type AdminRepository struct {
	db *sql.DB
}

func NewAdminRepository(db *sql.DB) *AdminRepository {
	return &AdminRepository{db: db}
}

func (r *AdminRepository) GetProjects(ctx context.Context, limit, offset int, sortBy, sortOrder string) ([]model.ProjectAdmin, int, error) {
	allowedColumns := map[string]bool{"id": true, "title": true, "created_at": true, "release_year": true}
	if !allowedColumns[sortBy] {
		sortBy = "id"
	}
	if sortOrder != "desc" {
		sortOrder = "asc"
	}

	var totalCount int
	countQuery := `SELECT COUNT(*) FROM (SELECT id FROM movies UNION ALL SELECT id FROM series) AS total`
	err := r.db.QueryRowContext(ctx, countQuery).Scan(&totalCount)
	if err != nil {
		return nil, 0, err
	}

	query := fmt.Sprintf(`
		SELECT id, type, title, description, release_year, director, producer, cover_image_url, category_id, age_rating_id, youtube_video_id, created_at
		FROM (
			SELECT id, 'movie' as type, title, description, release_year, director, producer, cover_image_url, category_id, age_rating_id, youtube_video_id, created_at FROM movies
			UNION ALL
			SELECT id, 'series' as type, title, description, release_year, director, producer, cover_image_url, category_id, age_rating_id, '' as youtube_video_id, created_at FROM series
		) as projects
		ORDER BY %s %s
		LIMIT $1 OFFSET $2`, sortBy, sortOrder)

	rows, err := r.db.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var projects []model.ProjectAdmin
	for rows.Next() {
		var p model.ProjectAdmin
		err := rows.Scan(&p.ID, &p.Type, &p.Title, &p.Description, &p.ReleaseYear, &p.Director, &p.Producer, &p.CoverImageURL, &p.CategoryID, &p.AgeRatingID, &p.YoutubeVideoID, &p.CreatedAt)
		if err != nil {
			return nil, 0, err
		}
		projects = append(projects, p)
	}

	return projects, totalCount, nil
}

func (r *AdminRepository) CreateProject(ctx context.Context, p *model.ProjectAdmin) error {
	if p.Type == "movie" {
		query := `
			INSERT INTO movies (title, description, release_year, director, producer, cover_image_url, category_id, age_rating_id, youtube_video_id)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
			RETURNING id, created_at`
		return r.db.QueryRowContext(ctx, query, p.Title, p.Description, p.ReleaseYear, p.Director, p.Producer, p.CoverImageURL, p.CategoryID, p.AgeRatingID, p.YoutubeVideoID).Scan(&p.ID, &p.CreatedAt)
	}
	query := `
		INSERT INTO series (title, description, release_year, director, producer, cover_image_url, category_id, age_rating_id)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id, created_at`
	return r.db.QueryRowContext(ctx, query, p.Title, p.Description, p.ReleaseYear, p.Director, p.Producer, p.CoverImageURL, p.CategoryID, p.AgeRatingID).Scan(&p.ID, &p.CreatedAt)
}

func (r *AdminRepository) UpdateProject(ctx context.Context, p *model.ProjectAdmin) error {
	if p.Type == "movie" {
		query := `
			UPDATE movies
			SET title = $1, description = $2, release_year = $3, director = $4, producer = $5, cover_image_url = $6, category_id = $7, age_rating_id = $8, youtube_video_id = $9, updated_at = CURRENT_TIMESTAMP
			WHERE id = $10`
		_, err := r.db.ExecContext(ctx, query, p.Title, p.Description, p.ReleaseYear, p.Director, p.Producer, p.CoverImageURL, p.CategoryID, p.AgeRatingID, p.YoutubeVideoID, p.ID)
		return err
	}
	query := `
		UPDATE series
		SET title = $1, description = $2, release_year = $3, director = $4, producer = $5, cover_image_url = $6, category_id = $7, age_rating_id = $8, updated_at = CURRENT_TIMESTAMP
		WHERE id = $9`
	_, err := r.db.ExecContext(ctx, query, p.Title, p.Description, p.ReleaseYear, p.Director, p.Producer, p.CoverImageURL, p.CategoryID, p.AgeRatingID, p.ID)
	return err
}

func (r *AdminRepository) DeleteProject(ctx context.Context, id int, projType string) error {
	if projType == "movie" {
		_, err := r.db.ExecContext(ctx, "DELETE FROM movies WHERE id = $1", id)
		return err
	}
	_, err := r.db.ExecContext(ctx, "DELETE FROM series WHERE id = $1", id)
	return err
}

func (r *AdminRepository) CreateSeason(ctx context.Context, seriesID int, s *model.SeasonAdmin) error {
	query := `
		INSERT INTO seasons (series_id, season_number)
		VALUES ($1, $2)
		RETURNING id`
	s.SeriesID = seriesID
	return r.db.QueryRowContext(ctx, query, seriesID, s.SeasonNumber).Scan(&s.ID)
}

func (r *AdminRepository) CreateEpisode(ctx context.Context, seasonID int, e *model.EpisodeAdmin) error {
	query := `
		INSERT INTO episodes (season_id, episode_number, title, youtube_video_id, duration)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id`
	e.SeasonID = seasonID
	return r.db.QueryRowContext(ctx, query, seasonID, e.EpisodeNumber, e.Title, e.YoutubeVideoID, e.Duration).Scan(&e.ID)
}

func (r *AdminRepository) UpdateFeaturedOrder(ctx context.Context, blockType string, items []model.FeaturedOrderItem) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	_, err = tx.ExecContext(ctx, "DELETE FROM featured_content WHERE block_type = $1", blockType)
	if err != nil {
		return err
	}

	for index, item := range items {
		var query string
		if item.Type == "movie" {
			query = `INSERT INTO featured_content (movie_id, series_id, block_type, sort_order) VALUES ($1, NULL, $2, $3)`
		} else {
			query = `INSERT INTO featured_content (movie_id, series_id, block_type, sort_order) VALUES (NULL, $1, $2, $3)`
		}
		_, err = tx.ExecContext(ctx, query, item.ID, blockType, index+1)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

func (r *AdminRepository) GetUserRoleName(ctx context.Context, userID int) (string, error) {
	var roleName string
	query := `
		SELECT r.name 
		FROM roles r
		JOIN users u ON u.role_id = r.id
		WHERE u.id = $1;
	`
	err := r.db.QueryRowContext(ctx, query, userID).Scan(&roleName)
	return roleName, err
}

func (r *AdminRepository) GetUsers(ctx context.Context) ([]model.UserResponse, error) {
	rows, err := r.db.QueryContext(ctx, "SELECT id, email, full_name, phone, birth_date, role_id, created_at FROM users")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var users []model.UserResponse
	for rows.Next() {
		var u model.UserResponse
		if err := rows.Scan(&u.ID, &u.Email, &u.FullName, &u.Phone, &u.BirthDate, &u.RoleID, &u.CreatedAt); err != nil {
			return nil, err
		}
		users = append(users, u)
	}
	return users, nil
}

func (r *AdminRepository) GetRoles(ctx context.Context) ([]model.Role, error) {
	rows, err := r.db.QueryContext(ctx, "SELECT id, name, permissions FROM roles")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var roles []model.Role
	for rows.Next() {
		var rle model.Role
		if err := rows.Scan(&rle.ID, &rle.Name, &rle.Permissions); err != nil {
			return nil, err
		}
		roles = append(roles, rle)
	}
	return roles, nil
}

func (r *AdminRepository) CreateRole(ctx context.Context, role *model.Role) error {
	return r.db.QueryRowContext(ctx, "INSERT INTO roles (name, permissions) VALUES ($1, $2) RETURNING id", role.Name, role.Permissions).Scan(&role.ID)
}

func (r *AdminRepository) UpdateRole(ctx context.Context, role *model.Role) error {
	_, err := r.db.ExecContext(ctx, "UPDATE roles SET name = $1, permissions = $2 WHERE id = $3", role.Name, role.Permissions, role.ID)
	return err
}

func (r *AdminRepository) DeleteRole(ctx context.Context, id int) error {
	_, err := r.db.ExecContext(ctx, "DELETE FROM roles WHERE id = $1", id)
	return err
}

func (r *AdminRepository) AssignRole(ctx context.Context, userID int, roleID int) error {
	_, err := r.db.ExecContext(ctx, "UPDATE users SET role_id = $1 WHERE id = $2", roleID, userID)
	return err
}

func (r *AdminRepository) CreateCategory(ctx context.Context, cat *model.Category) error {
	return r.db.QueryRowContext(ctx, "INSERT INTO categories (name) VALUES ($1) RETURNING id", cat.Name).Scan(&cat.ID)
}

func (r *AdminRepository) UpdateCategory(ctx context.Context, cat *model.Category) error {
	_, err := r.db.ExecContext(ctx, "UPDATE categories SET name = $1 WHERE id = $2", cat.Name, cat.ID)
	return err
}

func (r *AdminRepository) DeleteCategory(ctx context.Context, id int) error {
	_, err := r.db.ExecContext(ctx, "DELETE FROM categories WHERE id = $1", id)
	return err
}

func (r *AdminRepository) CreateGenre(ctx context.Context, g *model.Genre) error {
	return r.db.QueryRowContext(ctx, "INSERT INTO genres (name, icon_url) VALUES ($1, $2) RETURNING id", g.Name, g.IconURL).Scan(&g.ID)
}

func (r *AdminRepository) UpdateGenre(ctx context.Context, g *model.Genre) error {
	_, err := r.db.ExecContext(ctx, "UPDATE genres SET name = $1, icon_url = $2 WHERE id = $3", g.Name, g.IconURL, g.ID)
	return err
}

func (r *AdminRepository) DeleteGenre(ctx context.Context, id int) error {
	_, err := r.db.ExecContext(ctx, "DELETE FROM genres WHERE id = $1", id)
	return err
}

func (r *AdminRepository) CreateAgeRating(ctx context.Context, ar *model.AgeRating) error {
	return r.db.QueryRowContext(ctx, "INSERT INTO age_ratings (range) VALUES ($1) RETURNING id", ar.Range).Scan(&ar.ID)
}

func (r *AdminRepository) UpdateAgeRating(ctx context.Context, ar *model.AgeRating) error {
	_, err := r.db.ExecContext(ctx, "UPDATE age_ratings SET range = $1 WHERE id = $2", ar.Range, ar.ID)
	return err
}

func (r *AdminRepository) DeleteAgeRating(ctx context.Context, id int) error {
	_, err := r.db.ExecContext(ctx, "DELETE FROM age_ratings WHERE id = $1", id)
	return err
}