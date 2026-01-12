package service

import (
	"fmt"

	"github.com/razvan/library-app/internal/domain"
	"github.com/razvan/library-app/internal/repository"
)

type UserBookService struct {
	userBookRepo *repository.UserBookRepository
	bookRepo     *repository.BookRepository
}

func NewUserBookService(userBookRepo *repository.UserBookRepository, bookRepo *repository.BookRepository) *UserBookService {
	return &UserBookService{
		userBookRepo: userBookRepo,
		bookRepo:     bookRepo,
	}
}

// Reading list methods
func (s *UserBookService) AddToReadingList(userID, bookID int64, status domain.ReadingStatus) error {
	// Verify book exists
	_, err := s.bookRepo.GetByID(bookID)
	if err != nil {
		return fmt.Errorf("book not found")
	}

	return s.userBookRepo.AddToReadingList(userID, bookID, status)
}

func (s *UserBookService) RemoveFromReadingList(userID, bookID int64) error {
	return s.userBookRepo.RemoveFromReadingList(userID, bookID)
}

func (s *UserBookService) GetUserReadingList(userID int64, status string) ([]domain.UserBook, error) {
	return s.userBookRepo.GetUserBooks(userID, status)
}

// Favorite methods
func (s *UserBookService) AddToFavorites(userID, bookID int64) error {
	// Verify book exists
	_, err := s.bookRepo.GetByID(bookID)
	if err != nil {
		return fmt.Errorf("book not found")
	}

	return s.userBookRepo.AddToFavorites(userID, bookID)
}

func (s *UserBookService) RemoveFromFavorites(userID, bookID int64) error {
	return s.userBookRepo.RemoveFromFavorites(userID, bookID)
}

func (s *UserBookService) GetUserFavorites(userID int64) ([]domain.Favorite, error) {
	return s.userBookRepo.GetUserFavorites(userID)
}

func (s *UserBookService) IsFavorite(userID, bookID int64) (bool, error) {
	return s.userBookRepo.IsFavorite(userID, bookID)
}

// Comment methods
func (s *UserBookService) CreateComment(userID, bookID int64, content string) (*domain.Comment, error) {
	// Verify book exists
	_, err := s.bookRepo.GetByID(bookID)
	if err != nil {
		return nil, fmt.Errorf("book not found")
	}

	comment := &domain.Comment{
		UserID:  userID,
		BookID:  bookID,
		Content: content,
	}

	err = s.userBookRepo.CreateComment(comment)
	if err != nil {
		return nil, fmt.Errorf("failed to create comment: %w", err)
	}

	return comment, nil
}

func (s *UserBookService) GetBookComments(bookID int64) ([]domain.Comment, error) {
	return s.userBookRepo.GetBookComments(bookID)
}

func (s *UserBookService) UpdateComment(commentID, userID int64, content string) error {
	return s.userBookRepo.UpdateComment(commentID, userID, content)
}

func (s *UserBookService) DeleteComment(commentID, userID int64) error {
	return s.userBookRepo.DeleteComment(commentID, userID)
}
