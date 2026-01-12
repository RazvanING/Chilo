package domain

import (
	"time"
)

type Book struct {
	ID          int64     `json:"id" db:"id"`
	Title       string    `json:"title" db:"title"`
	Description string    `json:"description" db:"description"`
	CoverURL    string    `json:"cover_url" db:"cover_url"`
	ISBN        string    `json:"isbn" db:"isbn"`
	PublishedAt time.Time `json:"published_at" db:"published_at"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
	Authors     []Author  `json:"authors,omitempty"`
}

type Author struct {
	ID        int64     `json:"id" db:"id"`
	Name      string    `json:"name" db:"name"`
	Bio       string    `json:"bio" db:"bio"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

type BookCreate struct {
	Title       string   `json:"title" validate:"required,min=1,max=255"`
	Description string   `json:"description" validate:"required"`
	ISBN        string   `json:"isbn" validate:"omitempty,len=13"`
	PublishedAt string   `json:"published_at" validate:"required"`
	AuthorIDs   []int64  `json:"author_ids" validate:"required,min=1"`
}

type BookUpdate struct {
	Title       string   `json:"title" validate:"omitempty,min=1,max=255"`
	Description string   `json:"description" validate:"omitempty"`
	ISBN        string   `json:"isbn" validate:"omitempty,len=13"`
	PublishedAt string   `json:"published_at" validate:"omitempty"`
	AuthorIDs   []int64  `json:"author_ids" validate:"omitempty,min=1"`
}

type AuthorCreate struct {
	Name string `json:"name" validate:"required,min=1,max=255"`
	Bio  string `json:"bio" validate:"omitempty"`
}
