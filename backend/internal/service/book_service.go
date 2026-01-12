package service

import (
	"fmt"
	"time"

	"github.com/razvan/library-app/internal/domain"
	"github.com/razvan/library-app/internal/repository"
)

type BookService struct {
	bookRepo *repository.BookRepository
}

func NewBookService(bookRepo *repository.BookRepository) *BookService {
	return &BookService{bookRepo: bookRepo}
}

func (s *BookService) CreateBook(req *domain.BookCreate) (*domain.Book, error) {
	publishedAt, err := time.Parse("2006-01-02", req.PublishedAt)
	if err != nil {
		return nil, fmt.Errorf("invalid date format, use YYYY-MM-DD")
	}

	book := &domain.Book{
		Title:       req.Title,
		Description: req.Description,
		ISBN:        req.ISBN,
		PublishedAt: publishedAt,
	}

	err = s.bookRepo.Create(book, req.AuthorIDs)
	if err != nil {
		return nil, fmt.Errorf("failed to create book: %w", err)
	}

	return book, nil
}

func (s *BookService) GetBook(id int64) (*domain.Book, error) {
	return s.bookRepo.GetByID(id)
}

func (s *BookService) GetAllBooks(page, pageSize int) ([]domain.Book, error) {
	offset := (page - 1) * pageSize
	return s.bookRepo.GetAll(pageSize, offset)
}

func (s *BookService) SearchBooks(query string, page, pageSize int) ([]domain.Book, error) {
	offset := (page - 1) * pageSize
	return s.bookRepo.Search(query, pageSize, offset)
}

func (s *BookService) UpdateBook(id int64, req *domain.BookUpdate) (*domain.Book, error) {
	book, err := s.bookRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	if req.Title != "" {
		book.Title = req.Title
	}
	if req.Description != "" {
		book.Description = req.Description
	}
	if req.ISBN != "" {
		book.ISBN = req.ISBN
	}
	if req.PublishedAt != "" {
		publishedAt, err := time.Parse("2006-01-02", req.PublishedAt)
		if err != nil {
			return nil, fmt.Errorf("invalid date format, use YYYY-MM-DD")
		}
		book.PublishedAt = publishedAt
	}

	err = s.bookRepo.Update(book, req.AuthorIDs)
	if err != nil {
		return nil, fmt.Errorf("failed to update book: %w", err)
	}

	return book, nil
}

func (s *BookService) DeleteBook(id int64) error {
	return s.bookRepo.Delete(id)
}

func (s *BookService) UpdateBookCover(bookID int64, coverURL string) error {
	book, err := s.bookRepo.GetByID(bookID)
	if err != nil {
		return err
	}

	book.CoverURL = coverURL
	return s.bookRepo.Update(book, nil)
}

// Author methods
func (s *BookService) CreateAuthor(req *domain.AuthorCreate) (*domain.Author, error) {
	author := &domain.Author{
		Name: req.Name,
		Bio:  req.Bio,
	}

	err := s.bookRepo.CreateAuthor(author)
	if err != nil {
		return nil, fmt.Errorf("failed to create author: %w", err)
	}

	return author, nil
}

func (s *BookService) GetAuthor(id int64) (*domain.Author, error) {
	return s.bookRepo.GetAuthorByID(id)
}

func (s *BookService) GetAllAuthors(page, pageSize int) ([]domain.Author, error) {
	offset := (page - 1) * pageSize
	return s.bookRepo.GetAllAuthors(pageSize, offset)
}

func (s *BookService) UpdateAuthor(id int64, req *domain.AuthorCreate) (*domain.Author, error) {
	author, err := s.bookRepo.GetAuthorByID(id)
	if err != nil {
		return nil, err
	}

	author.Name = req.Name
	author.Bio = req.Bio

	err = s.bookRepo.UpdateAuthor(author)
	if err != nil {
		return nil, fmt.Errorf("failed to update author: %w", err)
	}

	return author, nil
}

func (s *BookService) DeleteAuthor(id int64) error {
	return s.bookRepo.DeleteAuthor(id)
}
