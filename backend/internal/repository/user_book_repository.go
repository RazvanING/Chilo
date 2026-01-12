package repository

import (
	"database/sql"
	"fmt"

	"github.com/razvan/library-app/internal/domain"
)

type UserBookRepository struct {
	db *sql.DB
}

func NewUserBookRepository(db *sql.DB) *UserBookRepository {
	return &UserBookRepository{db: db}
}

// UserBook methods (reading lists)
func (r *UserBookRepository) AddToReadingList(userID, bookID int64, status domain.ReadingStatus) error {
	query := `
		INSERT INTO user_books (user_id, book_id, status)
		VALUES ($1, $2, $3)
		ON CONFLICT (user_id, book_id)
		DO UPDATE SET status = $3, updated_at = CURRENT_TIMESTAMP`

	_, err := r.db.Exec(query, userID, bookID, status)
	return err
}

func (r *UserBookRepository) RemoveFromReadingList(userID, bookID int64) error {
	query := "DELETE FROM user_books WHERE user_id = $1 AND book_id = $2"
	_, err := r.db.Exec(query, userID, bookID)
	return err
}

func (r *UserBookRepository) GetUserBooks(userID int64, status string) ([]domain.UserBook, error) {
	var query string
	var args []interface{}

	if status != "" {
		query = `
			SELECT ub.id, ub.user_id, ub.book_id, ub.status, ub.created_at, ub.updated_at,
			       b.id, b.title, b.description, b.cover_url, b.isbn, b.published_at,
			       b.created_at, b.updated_at
			FROM user_books ub
			JOIN books b ON ub.book_id = b.id
			WHERE ub.user_id = $1 AND ub.status = $2
			ORDER BY ub.updated_at DESC`
		args = []interface{}{userID, status}
	} else {
		query = `
			SELECT ub.id, ub.user_id, ub.book_id, ub.status, ub.created_at, ub.updated_at,
			       b.id, b.title, b.description, b.cover_url, b.isbn, b.published_at,
			       b.created_at, b.updated_at
			FROM user_books ub
			JOIN books b ON ub.book_id = b.id
			WHERE ub.user_id = $1
			ORDER BY ub.updated_at DESC`
		args = []interface{}{userID}
	}

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var userBooks []domain.UserBook
	for rows.Next() {
		var ub domain.UserBook
		ub.Book = &domain.Book{}
		var coverURL, isbn sql.NullString

		err := rows.Scan(
			&ub.ID, &ub.UserID, &ub.BookID, &ub.Status, &ub.CreatedAt, &ub.UpdatedAt,
			&ub.Book.ID, &ub.Book.Title, &ub.Book.Description, &coverURL, &isbn,
			&ub.Book.PublishedAt, &ub.Book.CreatedAt, &ub.Book.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		if coverURL.Valid {
			ub.Book.CoverURL = coverURL.String
		}
		if isbn.Valid {
			ub.Book.ISBN = isbn.String
		}

		userBooks = append(userBooks, ub)
	}

	return userBooks, nil
}

// Favorite methods
func (r *UserBookRepository) AddToFavorites(userID, bookID int64) error {
	query := `
		INSERT INTO favorites (user_id, book_id)
		VALUES ($1, $2)
		ON CONFLICT (user_id, book_id) DO NOTHING`

	_, err := r.db.Exec(query, userID, bookID)
	return err
}

func (r *UserBookRepository) RemoveFromFavorites(userID, bookID int64) error {
	query := "DELETE FROM favorites WHERE user_id = $1 AND book_id = $2"
	result, err := r.db.Exec(query, userID, bookID)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return fmt.Errorf("favorite not found")
	}

	return nil
}

func (r *UserBookRepository) GetUserFavorites(userID int64) ([]domain.Favorite, error) {
	query := `
		SELECT f.id, f.user_id, f.book_id, f.created_at,
		       b.id, b.title, b.description, b.cover_url, b.isbn, b.published_at,
		       b.created_at, b.updated_at
		FROM favorites f
		JOIN books b ON f.book_id = b.id
		WHERE f.user_id = $1
		ORDER BY f.created_at DESC`

	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var favorites []domain.Favorite
	for rows.Next() {
		var fav domain.Favorite
		fav.Book = &domain.Book{}
		var coverURL, isbn sql.NullString

		err := rows.Scan(
			&fav.ID, &fav.UserID, &fav.BookID, &fav.CreatedAt,
			&fav.Book.ID, &fav.Book.Title, &fav.Book.Description, &coverURL, &isbn,
			&fav.Book.PublishedAt, &fav.Book.CreatedAt, &fav.Book.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		if coverURL.Valid {
			fav.Book.CoverURL = coverURL.String
		}
		if isbn.Valid {
			fav.Book.ISBN = isbn.String
		}

		favorites = append(favorites, fav)
	}

	return favorites, nil
}

func (r *UserBookRepository) IsFavorite(userID, bookID int64) (bool, error) {
	query := "SELECT EXISTS(SELECT 1 FROM favorites WHERE user_id = $1 AND book_id = $2)"
	var exists bool
	err := r.db.QueryRow(query, userID, bookID).Scan(&exists)
	return exists, err
}

// Comment methods
func (r *UserBookRepository) CreateComment(comment *domain.Comment) error {
	query := `
		INSERT INTO comments (user_id, book_id, content)
		VALUES ($1, $2, $3)
		RETURNING id, created_at, updated_at`

	return r.db.QueryRow(query, comment.UserID, comment.BookID, comment.Content).Scan(
		&comment.ID, &comment.CreatedAt, &comment.UpdatedAt,
	)
}

func (r *UserBookRepository) GetBookComments(bookID int64) ([]domain.Comment, error) {
	query := `
		SELECT c.id, c.user_id, c.book_id, c.content, c.created_at, c.updated_at, u.username
		FROM comments c
		JOIN users u ON c.user_id = u.id
		WHERE c.book_id = $1
		ORDER BY c.created_at DESC`

	rows, err := r.db.Query(query, bookID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var comments []domain.Comment
	for rows.Next() {
		var comment domain.Comment
		err := rows.Scan(
			&comment.ID, &comment.UserID, &comment.BookID, &comment.Content,
			&comment.CreatedAt, &comment.UpdatedAt, &comment.Username,
		)
		if err != nil {
			return nil, err
		}
		comments = append(comments, comment)
	}

	return comments, nil
}

func (r *UserBookRepository) UpdateComment(commentID, userID int64, content string) error {
	query := `UPDATE comments SET content = $1 WHERE id = $2 AND user_id = $3`
	result, err := r.db.Exec(query, content, commentID, userID)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return fmt.Errorf("comment not found or unauthorized")
	}

	return nil
}

func (r *UserBookRepository) DeleteComment(commentID, userID int64) error {
	query := `DELETE FROM comments WHERE id = $1 AND user_id = $2`
	result, err := r.db.Exec(query, commentID, userID)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return fmt.Errorf("comment not found or unauthorized")
	}

	return nil
}
