package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/razvan/library-app/internal/domain"
	"github.com/razvan/library-app/internal/middleware"
	"github.com/razvan/library-app/internal/service"
	"github.com/razvan/library-app/internal/utils"
	"github.com/razvan/library-app/pkg/validator"
)

type AuthHandler struct {
	authService *service.AuthService
}

func NewAuthHandler(authService *service.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req domain.UserRegistration
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if err := validator.Validate(&req); err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	user, err := h.authService.Register(&req)
	if err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	utils.SuccessResponseWithData(w, map[string]interface{}{
		"user":    user,
		"message": "Registration successful. Please verify your email.",
	})
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req domain.UserLogin
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if err := validator.Validate(&req); err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	tokens, user, err := h.authService.Login(&req)
	if err != nil {
		utils.ErrorResponse(w, http.StatusUnauthorized, err.Error())
		return
	}

	// If 2FA is enabled, don't return tokens yet
	if user.TwoFactorEnabled {
		utils.SuccessResponseWithData(w, map[string]interface{}{
			"requires_2fa": true,
			"user_id":      user.ID,
			"message":      "Please provide 2FA code",
		})
		return
	}

	utils.SuccessResponseWithData(w, map[string]interface{}{
		"tokens": tokens,
		"user":   user,
	})
}

func (h *AuthHandler) SetupTwoFactor(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())

	setup, err := h.authService.SetupTwoFactor(userID)
	if err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	utils.SuccessResponseWithData(w, setup)
}

func (h *AuthHandler) VerifyTwoFactorSetup(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())

	var req domain.TwoFactorVerify
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if err := validator.Validate(&req); err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	err := h.authService.VerifyAndEnableTwoFactor(userID, req.Code)
	if err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	utils.SuccessResponseWithMessage(w, "2FA enabled successfully")
}

func (h *AuthHandler) VerifyTwoFactorLogin(w http.ResponseWriter, r *http.Request) {
	var req struct {
		UserID int64  `json:"user_id" validate:"required"`
		Code   string `json:"code" validate:"required,len=6"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if err := validator.Validate(&req); err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	tokens, err := h.authService.VerifyTwoFactorLogin(req.UserID, req.Code)
	if err != nil {
		utils.ErrorResponse(w, http.StatusUnauthorized, err.Error())
		return
	}

	utils.SuccessResponseWithData(w, map[string]interface{}{
		"tokens": tokens,
	})
}

func (h *AuthHandler) DisableTwoFactor(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())

	var req domain.TwoFactorVerify
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if err := validator.Validate(&req); err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	err := h.authService.DisableTwoFactor(userID, req.Code)
	if err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	utils.SuccessResponseWithMessage(w, "2FA disabled successfully")
}

func (h *AuthHandler) MakeAdmin(w http.ResponseWriter, r *http.Request) {
	adminUserID := middleware.GetUserID(r.Context())
	vars := mux.Vars(r)
	targetUserID, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, "invalid user ID")
		return
	}

	err = h.authService.MakeAdmin(adminUserID, targetUserID)
	if err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	utils.SuccessResponseWithMessage(w, "User is now an admin")
}

func (h *AuthHandler) GetMe(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())
	email := middleware.GetEmail(r.Context())
	isAdmin := middleware.IsAdmin(r.Context())

	utils.SuccessResponseWithData(w, map[string]interface{}{
		"id":       userID,
		"email":    email,
		"is_admin": isAdmin,
	})
}
