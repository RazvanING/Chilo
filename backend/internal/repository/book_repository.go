package repository

import (
	"database/sql"
	"fmt"

	"github.com/razvan/library-app/internal/domain"
)

type BookRepository struct {
	db *sql.DB
}

func NewBookRepository(db *sql.DB) *BookRepository {
	return &BookRepository{db: db}
}

func (r *BookRepository) Create(book *domain.Book, authorIDs []int64) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	query := `
		INSERT INTO books (title, description, cover_url, isbn, published_at)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, created_at, updated_at`

	err = tx.QueryRow(
		query,
		book.Title,
		book.Description,
		book.CoverURL,
		book.ISBN,
		book.PublishedAt,
	).Scan(&book.ID, &book.CreatedAt, &book.UpdatedAt)

	if err != nil {
		return err
	}

	// Link authors
	for _, authorID := range authorIDs {
		_, err = tx.Exec(
			"INSERT INTO book_authors (book_id, author_id) VALUES ($1, $2)",
			book.ID, authorID,
		)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

func (r *BookRepository) GetByID(id int64) (*domain.Book, error) {
	book := &domain.Book{}
	query := `
		SELECT id, title, description, cover_url, isbn, published_at,
		       created_at, updated_at
		FROM books WHERE id = $1`

	var coverURL, isbn sql.NullString
	err := r.db.QueryRow(query, id).Scan(
		&book.ID, &book.Title, &book.Description, &coverURL, &isbn,
		&book.PublishedAt, &book.CreatedAt, &book.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("book not found")
	}
	if err != nil {
		return nil, err
	}

	if coverURL.Valid {
		book.CoverURL = coverURL.String
	}
	if isbn.Valid {
		book.ISBN = isbn.String
	}

	// Get authors
	book.Authors, err = r.GetBookAuthors(id)
	if err != nil {
		return nil, err
	}

	return book, nil
}

func (r *BookRepository) GetAll(limit, offset int) ([]domain.Book, error) {
	query := `
		SELECT id, title, description, cover_url, isbn, published_at,
		       created_at, updated_at
		FROM books
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2`

	rows, err := r.db.Query(query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var books []domain.Book
	for rows.Next() {
		var book domain.Book
		var coverURL, isbn sql.NullString

		err := rows.Scan(
			&book.ID, &book.Title, &book.Description, &coverURL, &isbn,
			&book.PublishedAt, &book.CreatedAt, &book.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		if coverURL.Valid {
			book.CoverURL = coverURL.String
		}
		if isbn.Valid {
			book.ISBN = isbn.String
		}

		// Get authors for each book
		book.Authors, _ = r.GetBookAuthors(book.ID)
		books = append(books, book)
	}

	return books, nil
}

func (r *BookRepository) Search(query string, limit, offset int) ([]domain.Book, error) {
	sqlQuery := `
		SELECT DISTINCT b.id, b.title, b.description, b.cover_url, b.isbn,
		       b.published_at, b.created_at, b.updated_at
		FROM books b
		LEFT JOIN book_authors ba ON b.id = ba.book_id
		LEFT JOIN authors a ON ba.author_id = a.id
		WHERE b.title ILIKE $1 OR b.description ILIKE $1 OR a.name ILIKE $1
		ORDER BY b.created_at DESC
		LIMIT $2 OFFSET $3`

	searchPattern := "%" + query + "%"
	rows, err := r.db.Query(sqlQuery, searchPattern, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var books []domain.Book
	for rows.Next() {
		var book domain.Book
		var coverURL, isbn sql.NullString

		err := rows.Scan(
			&book.ID, &book.Title, &book.Description, &coverURL, &isbn,
			&book.PublishedAt, &book.CreatedAt, &book.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		if coverURL.Valid {
			book.CoverURL = coverURL.String
		}
		if isbn.Valid {
			book.ISBN = isbn.String
		}

		book.Authors, _ = r.GetBookAuthors(book.ID)
		books = append(books, book)
	}

	return books, nil
}

func (r *BookRepository) Update(book *domain.Book, authorIDs []int64) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	query := `
		UPDATE books
		SET title = $1, description = $2, cover_url = $3, isbn = $4, published_at = $5
		WHERE id = $6`

	_, err = tx.Exec(
		query,
		book.Title,
		book.Description,
		book.CoverURL,
		book.ISBN,
		book.PublishedAt,
		book.ID,
	)
	if err != nil {
		return err
	}

	// Update authors if provided
	if len(authorIDs) > 0 {
		// Remove existing authors
		_, err = tx.Exec("DELETE FROM book_authors WHERE book_id = $1", book.ID)
		if err != nil {
			return err
		}

		// Add new authors
		for _, authorID := range authorIDs {
			_, err = tx.Exec(
				"INSERT INTO book_authors (book_id, author_id) VALUES ($1, $2)",
				book.ID, authorID,
			)
			if err != nil {
				return err
			}
		}
	}

	return tx.Commit()
}

func (r *BookRepository) Delete(id int64) error {
	query := "DELETE FROM books WHERE id = $1"
	_, err := r.db.Exec(query, id)
	return err
}

func (r *BookRepository) GetBookAuthors(bookID int64) ([]domain.Author, error) {
	query := `
		SELECT a.id, a.name, a.bio, a.created_at, a.updated_at
		FROM authors a
		JOIN book_authors ba ON a.id = ba.author_id
		WHERE ba.book_id = $1`

	rows, err := r.db.Query(query, bookID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var authors []domain.Author
	for rows.Next() {
		var author domain.Author
		var bio sql.NullString

		err := rows.Scan(
			&author.ID, &author.Name, &bio,
			&author.CreatedAt, &author.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		if bio.Valid {
			author.Bio = bio.String
		}

		authors = append(authors, author)
	}

	return authors, nil
}

// Author methods
func (r *BookRepository) CreateAuthor(author *domain.Author) error {
	query := `
		INSERT INTO authors (name, bio)
		VALUES ($1, $2)
		RETURNING id, created_at, updated_at`

	return r.db.QueryRow(query, author.Name, author.Bio).Scan(
		&author.ID, &author.CreatedAt, &author.UpdatedAt,
	)
}

func (r *BookRepository) GetAuthorByID(id int64) (*domain.Author, error) {
	author := &domain.Author{}
	query := `SELECT id, name, bio, created_at, updated_at FROM authors WHERE id = $1`

	var bio sql.NullString
	err := r.db.QueryRow(query, id).Scan(
		&author.ID, &author.Name, &bio,
		&author.CreatedAt, &author.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("author not found")
	}
	if err != nil {
		return nil, err
	}

	if bio.Valid {
		author.Bio = bio.String
	}

	return author, nil
}

func (r *BookRepository) GetAllAuthors(limit, offset int) ([]domain.Author, error) {
	query := `
		SELECT id, name, bio, created_at, updated_at
		FROM authors
		ORDER BY name ASC
		LIMIT $1 OFFSET $2`

	rows, err := r.db.Query(query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var authors []domain.Author
	for rows.Next() {
		var author domain.Author
		var bio sql.NullString

		err := rows.Scan(
			&author.ID, &author.Name, &bio,
			&author.CreatedAt, &author.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		if bio.Valid {
			author.Bio = bio.String
		}

		authors = append(authors, author)
	}

	return authors, nil
}

func (r *BookRepository) UpdateAuthor(author *domain.Author) error {
	query := `UPDATE authors SET name = $1, bio = $2 WHERE id = $3`
	_, err := r.db.Exec(query, author.Name, author.Bio, author.ID)
	return err
}

func (r *BookRepository) DeleteAuthor(id int64) error {
	query := "DELETE FROM authors WHERE id = $1"
	_, err := r.db.Exec(query, id)
	return err
}
