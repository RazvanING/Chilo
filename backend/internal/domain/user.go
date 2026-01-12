package domain

import (
	"time"
)

type User struct {
	ID              int64     `json:"id" db:"id"`
	Email           string    `json:"email" db:"email"`
	Username        string    `json:"username" db:"username"`
	PasswordHash    string    `json:"-" db:"password_hash"`
	IsAdmin         bool      `json:"is_admin" db:"is_admin"`
	GoogleID        string    `json:"-" db:"google_id"`
	TwoFactorSecret string    `json:"-" db:"two_factor_secret"`
	TwoFactorEnabled bool     `json:"two_factor_enabled" db:"two_factor_enabled"`
	EmailVerified   bool      `json:"email_verified" db:"email_verified"`
	CreatedAt       time.Time `json:"created_at" db:"created_at"`
	UpdatedAt       time.Time `json:"updated_at" db:"updated_at"`
}

type UserRegistration struct {
	Email    string `json:"email" validate:"required,email"`
	Username string `json:"username" validate:"required,min=3,max=50"`
	Password string `json:"password" validate:"required,min=8"`
}

type UserLogin struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type TwoFactorSetup struct {
	Secret string `json:"secret"`
	QRCode string `json:"qr_code"`
}

type TwoFactorVerify struct {
	Code string `json:"code" validate:"required,len=6"`
}
