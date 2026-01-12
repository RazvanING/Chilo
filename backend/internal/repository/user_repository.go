package repository

import (
	"database/sql"
	"fmt"

	"github.com/razvan/library-app/internal/domain"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(user *domain.User) error {
	query := `
		INSERT INTO users (email, username, password_hash, google_id, email_verified)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, created_at, updated_at`

	err := r.db.QueryRow(
		query,
		user.Email,
		user.Username,
		user.PasswordHash,
		sql.NullString{String: user.GoogleID, Valid: user.GoogleID != ""},
		user.EmailVerified,
	).Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt)

	return err
}

func (r *UserRepository) GetByEmail(email string) (*domain.User, error) {
	user := &domain.User{}
	query := `
		SELECT id, email, username, password_hash, is_admin, google_id,
		       two_factor_secret, two_factor_enabled, email_verified,
		       created_at, updated_at
		FROM users WHERE email = $1`

	var googleID sql.NullString
	err := r.db.QueryRow(query, email).Scan(
		&user.ID, &user.Email, &user.Username, &user.PasswordHash,
		&user.IsAdmin, &googleID, &user.TwoFactorSecret,
		&user.TwoFactorEnabled, &user.EmailVerified,
		&user.CreatedAt, &user.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("user not found")
	}
	if err != nil {
		return nil, err
	}

	if googleID.Valid {
		user.GoogleID = googleID.String
	}

	return user, nil
}

func (r *UserRepository) GetByID(id int64) (*domain.User, error) {
	user := &domain.User{}
	query := `
		SELECT id, email, username, password_hash, is_admin, google_id,
		       two_factor_secret, two_factor_enabled, email_verified,
		       created_at, updated_at
		FROM users WHERE id = $1`

	var googleID sql.NullString
	err := r.db.QueryRow(query, id).Scan(
		&user.ID, &user.Email, &user.Username, &user.PasswordHash,
		&user.IsAdmin, &googleID, &user.TwoFactorSecret,
		&user.TwoFactorEnabled, &user.EmailVerified,
		&user.CreatedAt, &user.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("user not found")
	}
	if err != nil {
		return nil, err
	}

	if googleID.Valid {
		user.GoogleID = googleID.String
	}

	return user, nil
}

func (r *UserRepository) GetByGoogleID(googleID string) (*domain.User, error) {
	user := &domain.User{}
	query := `
		SELECT id, email, username, password_hash, is_admin, google_id,
		       two_factor_secret, two_factor_enabled, email_verified,
		       created_at, updated_at
		FROM users WHERE google_id = $1`

	var gID sql.NullString
	err := r.db.QueryRow(query, googleID).Scan(
		&user.ID, &user.Email, &user.Username, &user.PasswordHash,
		&user.IsAdmin, &gID, &user.TwoFactorSecret,
		&user.TwoFactorEnabled, &user.EmailVerified,
		&user.CreatedAt, &user.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("user not found")
	}
	if err != nil {
		return nil, err
	}

	if gID.Valid {
		user.GoogleID = gID.String
	}

	return user, nil
}

func (r *UserRepository) Update(user *domain.User) error {
	query := `
		UPDATE users
		SET username = $1, is_admin = $2, two_factor_secret = $3,
		    two_factor_enabled = $4, email_verified = $5
		WHERE id = $6`

	_, err := r.db.Exec(
		query,
		user.Username,
		user.IsAdmin,
		user.TwoFactorSecret,
		user.TwoFactorEnabled,
		user.EmailVerified,
		user.ID,
	)

	return err
}

func (r *UserRepository) UpdatePassword(userID int64, passwordHash string) error {
	query := `UPDATE users SET password_hash = $1 WHERE id = $2`
	_, err := r.db.Exec(query, passwordHash, userID)
	return err
}
