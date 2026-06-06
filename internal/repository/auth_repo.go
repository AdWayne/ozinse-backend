package repository

import (
	"context"
	"database/sql"
	"errors"
	"time"
	"ozinse-backend/internal/model"
)

type AuthRepository interface {
	CreateUser(ctx context.Context, email, passwordHash, fullName string) error
	GetUserByEmail(ctx context.Context, email string) (*model.User, error)
	GetUserByID(ctx context.Context, id int) (*model.User, error)
	UpdateProfile(ctx context.Context, id int, fullName, phone, birthDate string) error
	UpdatePassword(ctx context.Context, id int, newPasswordHash string) error
	SaveRefreshToken(ctx context.Context, userID int, token string, expiresAt time.Time) error
	GetRefreshToken(ctx context.Context, token string) (int, time.Time, error)
	DeleteRefreshToken(ctx context.Context, token string) error
}

type authRepository struct {
	db *sql.DB
}

func NewAuthRepository(db *sql.DB) AuthRepository {
	return &authRepository{db: db}
}

func (r *authRepository) CreateUser(ctx context.Context, email, passwordHash, fullName string) error {
	query := `
		INSERT INTO users (email, password_hash, full_name, role_id) 
		VALUES ($1, $2, $3, (SELECT id FROM roles WHERE name = 'Пользователь' LIMIT 1));`
	_, err := r.db.ExecContext(ctx, query, email, passwordHash, fullName)
	return err
}

func (r *authRepository) GetUserByEmail(ctx context.Context, email string) (*model.User, error) {
	var u model.User
	query := `
		SELECT u.id, u.email, u.password_hash, u.full_name, COALESCE(u.phone, ''), COALESCE(u.birth_date::text, ''), r.name, u.created_at
		FROM users u
		LEFT JOIN roles r ON u.role_id = r.id
		WHERE u.email = $1;`
	
	err := r.db.QueryRowContext(ctx, query, email).Scan(&u.ID, &u.Email, &u.PasswordHash, &u.FullName, &u.Phone, &u.BirthDate, &u.RoleName, &u.CreatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	return &u, err
}

func (r *authRepository) GetUserByID(ctx context.Context, id int) (*model.User, error) {
	var u model.User
	query := `
		SELECT u.id, u.email, u.password_hash, u.full_name, COALESCE(u.phone, ''), COALESCE(u.birth_date::text, ''), r.name, u.created_at
		FROM users u
		LEFT JOIN roles r ON u.role_id = r.id
		WHERE u.id = $1;`
	
	err := r.db.QueryRowContext(ctx, query, id).Scan(&u.ID, &u.Email, &u.PasswordHash, &u.FullName, &u.Phone, &u.BirthDate, &u.RoleName, &u.CreatedAt)
	return &u, err
}

func (r *authRepository) UpdateProfile(ctx context.Context, id int, fullName, phone, birthDate string) error {
	var err error
	if birthDate == "" {
		query := `UPDATE users SET full_name = $1, phone = $2, birth_date = NULL WHERE id = $3;`
		_, err = r.db.ExecContext(ctx, query, fullName, phone, id)
	} else {
		query := `UPDATE users SET full_name = $1, phone = $2, birth_date = $3 WHERE id = $4;`
		_, err = r.db.ExecContext(ctx, query, fullName, phone, birthDate, id)
	}
	return err
}

func (r *authRepository) UpdatePassword(ctx context.Context, id int, newPasswordHash string) error {
	query := `UPDATE users SET password_hash = $1 WHERE id = $2;`
	_, err := r.db.ExecContext(ctx, query, newPasswordHash, id)
	return err
}

func (r *authRepository) SaveRefreshToken(ctx context.Context, userID int, token string, expiresAt time.Time) error {
	// Удаляем старые токены пользователя, чтобы не забивать базу
	_, _ = r.db.ExecContext(ctx, "DELETE FROM refresh_tokens WHERE user_id = $1;", userID)
	
	query := `INSERT INTO refresh_tokens (user_id, token, expires_at) VALUES ($1, $2, $3);`
	_, err := r.db.ExecContext(ctx, query, userID, token, expiresAt)
	return err
}

func (r *authRepository) GetRefreshToken(ctx context.Context, token string) (int, time.Time, error) {
	var userID int
	var expiresAt time.Time
	query := `SELECT user_id, expires_at FROM refresh_tokens WHERE token = $1;`
	err := r.db.QueryRowContext(ctx, query, token).Scan(&userID, &expiresAt)
	return userID, expiresAt, err
}

func (r *authRepository) DeleteRefreshToken(ctx context.Context, token string) error {
	query := `DELETE FROM refresh_tokens WHERE token = $1;`
	_, err := r.db.ExecContext(ctx, query, token)
	return err
}