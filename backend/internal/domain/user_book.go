package domain

import (
	"time"
)

type ReadingStatus string

const (
	StatusWantToRead ReadingStatus = "want_to_read"
	StatusReading    ReadingStatus = "reading"
	StatusRead       ReadingStatus = "read"
)

type UserBook struct {
	ID        int64         `json:"id" db:"id"`
	UserID    int64         `json:"user_id" db:"user_id"`
	BookID    int64         `json:"book_id" db:"book_id"`
	Status    ReadingStatus `json:"status" db:"status"`
	CreatedAt time.Time     `json:"created_at" db:"created_at"`
	UpdatedAt time.Time     `json:"updated_at" db:"updated_at"`
	Book      *Book         `json:"book,omitempty"`
}

type Favorite struct {
	ID        int64     `json:"id" db:"id"`
	UserID    int64     `json:"user_id" db:"user_id"`
	BookID    int64     `json:"book_id" db:"book_id"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	Book      *Book     `json:"book,omitempty"`
}

type Comment struct {
	ID        int64     `json:"id" db:"id"`
	UserID    int64     `json:"user_id" db:"user_id"`
	BookID    int64     `json:"book_id" db:"book_id"`
	Content   string    `json:"content" db:"content"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
	Username  string    `json:"username,omitempty" db:"username"`
}

type UserBookRequest struct {
	Status ReadingStatus `json:"status" validate:"required,oneof=want_to_read reading read"`
}

type CommentCreate struct {
	Content string `json:"content" validate:"required,min=1,max=1000"`
}
