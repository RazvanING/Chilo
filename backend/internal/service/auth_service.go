package service

import (
	"fmt"

	"github.com/razvan/library-app/internal/domain"
	"github.com/razvan/library-app/internal/repository"
	"github.com/razvan/library-app/pkg/auth"
)

type AuthService struct {
	userRepo  *repository.UserRepository
	jwtSecret string
	appName   string
}

func NewAuthService(userRepo *repository.UserRepository, jwtSecret, appName string) *AuthService {
	return &AuthService{
		userRepo:  userRepo,
		jwtSecret: jwtSecret,
		appName:   appName,
	}
}

func (s *AuthService) Register(reg *domain.UserRegistration) (*domain.User, error) {
	// Check if user already exists
	existing, _ := s.userRepo.GetByEmail(reg.Email)
	if existing != nil {
		return nil, fmt.Errorf("user with this email already exists")
	}

	// Hash password
	passwordHash, err := auth.HashPassword(reg.Password)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	user := &domain.User{
		Email:         reg.Email,
		Username:      reg.Username,
		PasswordHash:  passwordHash,
		EmailVerified: false,
		IsAdmin:       false,
	}

	err = s.userRepo.Create(user)
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return user, nil
}

func (s *AuthService) Login(login *domain.UserLogin) (*auth.TokenPair, *domain.User, error) {
	user, err := s.userRepo.GetByEmail(login.Email)
	if err != nil {
		return nil, nil, fmt.Errorf("invalid credentials")
	}

	// Check password
	if !auth.CheckPassword(login.Password, user.PasswordHash) {
		return nil, nil, fmt.Errorf("invalid credentials")
	}

	// Don't generate token yet if 2FA is enabled
	if user.TwoFactorEnabled {
		return nil, user, nil
	}

	// Generate JWT tokens
	tokens, err := auth.GenerateTokenPair(user.ID, user.Email, user.IsAdmin, s.jwtSecret)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to generate tokens: %w", err)
	}

	return tokens, user, nil
}

func (s *AuthService) LoginWithGoogle(googleID, email, username string) (*auth.TokenPair, *domain.User, error) {
	// Try to find user by Google ID
	user, err := s.userRepo.GetByGoogleID(googleID)
	if err != nil {
		// User doesn't exist, create new user
		user = &domain.User{
			Email:         email,
			Username:      username,
			GoogleID:      googleID,
			EmailVerified: true,
			IsAdmin:       false,
		}

		err = s.userRepo.Create(user)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to create user: %w", err)
		}
	}

	// Generate JWT tokens (no 2FA for Google OAuth users in this implementation)
	tokens, err := auth.GenerateTokenPair(user.ID, user.Email, user.IsAdmin, s.jwtSecret)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to generate tokens: %w", err)
	}

	return tokens, user, nil
}

func (s *AuthService) SetupTwoFactor(userID int64) (*domain.TwoFactorSetup, error) {
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return nil, fmt.Errorf("user not found")
	}

	if user.TwoFactorEnabled {
		return nil, fmt.Errorf("2FA is already enabled")
	}

	// Generate TOTP secret
	secret, qrCode, err := auth.GenerateTOTPSecret(user.Email, s.appName)
	if err != nil {
		return nil, fmt.Errorf("failed to generate TOTP secret: %w", err)
	}

	// Save secret (but don't enable yet)
	user.TwoFactorSecret = secret
	err = s.userRepo.Update(user)
	if err != nil {
		return nil, fmt.Errorf("failed to save 2FA secret: %w", err)
	}

	return &domain.TwoFactorSetup{
		Secret: secret,
		QRCode: qrCode,
	}, nil
}

func (s *AuthService) VerifyAndEnableTwoFactor(userID int64, code string) error {
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return fmt.Errorf("user not found")
	}

	if user.TwoFactorSecret == "" {
		return fmt.Errorf("2FA setup not initiated")
	}

	// Verify the code
	if !auth.ValidateTOTP(code, user.TwoFactorSecret) {
		return fmt.Errorf("invalid verification code")
	}

	// Enable 2FA
	user.TwoFactorEnabled = true
	err = s.userRepo.Update(user)
	if err != nil {
		return fmt.Errorf("failed to enable 2FA: %w", err)
	}

	return nil
}

func (s *AuthService) VerifyTwoFactorLogin(userID int64, code string) (*auth.TokenPair, error) {
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return nil, fmt.Errorf("user not found")
	}

	if !user.TwoFactorEnabled {
		return nil, fmt.Errorf("2FA is not enabled for this user")
	}

	// Verify the code
	if !auth.ValidateTOTP(code, user.TwoFactorSecret) {
		return nil, fmt.Errorf("invalid verification code")
	}

	// Generate JWT tokens
	tokens, err := auth.GenerateTokenPair(user.ID, user.Email, user.IsAdmin, s.jwtSecret)
	if err != nil {
		return nil, fmt.Errorf("failed to generate tokens: %w", err)
	}

	return tokens, nil
}

func (s *AuthService) DisableTwoFactor(userID int64, code string) error {
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return fmt.Errorf("user not found")
	}

	if !user.TwoFactorEnabled {
		return fmt.Errorf("2FA is not enabled")
	}

	// Verify the code before disabling
	if !auth.ValidateTOTP(code, user.TwoFactorSecret) {
		return fmt.Errorf("invalid verification code")
	}

	// Disable 2FA
	user.TwoFactorEnabled = false
	user.TwoFactorSecret = ""
	err = s.userRepo.Update(user)
	if err != nil {
		return fmt.Errorf("failed to disable 2FA: %w", err)
	}

	return nil
}

func (s *AuthService) MakeAdmin(adminUserID, targetUserID int64) error {
	// Verify that the requesting user is an admin
	admin, err := s.userRepo.GetByID(adminUserID)
	if err != nil {
		return fmt.Errorf("admin user not found")
	}

	if !admin.IsAdmin {
		return fmt.Errorf("unauthorized: only admins can make other users admins")
	}

	// Get target user
	targetUser, err := s.userRepo.GetByID(targetUserID)
	if err != nil {
		return fmt.Errorf("target user not found")
	}

	// Make admin
	targetUser.IsAdmin = true
	err = s.userRepo.Update(targetUser)
	if err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}

	return nil
}
