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

type UserBookHandler struct {
	userBookService *service.UserBookService
}

func NewUserBookHandler(userBookService *service.UserBookService) *UserBookHandler {
	return &UserBookHandler{userBookService: userBookService}
}

// Reading list handlers
func (h *UserBookHandler) AddToReadingList(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())
	vars := mux.Vars(r)
	bookID, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, "invalid book ID")
		return
	}

	var req domain.UserBookRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if err := validator.Validate(&req); err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	err = h.userBookService.AddToReadingList(userID, bookID, req.Status)
	if err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	utils.SuccessResponseWithMessage(w, "Book added to reading list")
}

func (h *UserBookHandler) RemoveFromReadingList(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())
	vars := mux.Vars(r)
	bookID, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, "invalid book ID")
		return
	}

	err = h.userBookService.RemoveFromReadingList(userID, bookID)
	if err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	utils.SuccessResponseWithMessage(w, "Book removed from reading list")
}

func (h *UserBookHandler) GetReadingList(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())
	status := r.URL.Query().Get("status")

	books, err := h.userBookService.GetUserReadingList(userID, status)
	if err != nil {
		utils.ErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SuccessResponseWithData(w, books)
}

// Favorite handlers
func (h *UserBookHandler) AddToFavorites(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())
	vars := mux.Vars(r)
	bookID, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, "invalid book ID")
		return
	}

	err = h.userBookService.AddToFavorites(userID, bookID)
	if err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	utils.SuccessResponseWithMessage(w, "Book added to favorites")
}

func (h *UserBookHandler) RemoveFromFavorites(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())
	vars := mux.Vars(r)
	bookID, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, "invalid book ID")
		return
	}

	err = h.userBookService.RemoveFromFavorites(userID, bookID)
	if err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	utils.SuccessResponseWithMessage(w, "Book removed from favorites")
}

func (h *UserBookHandler) GetFavorites(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())

	favorites, err := h.userBookService.GetUserFavorites(userID)
	if err != nil {
		utils.ErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SuccessResponseWithData(w, favorites)
}

// Comment handlers
func (h *UserBookHandler) CreateComment(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())
	vars := mux.Vars(r)
	bookID, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, "invalid book ID")
		return
	}

	var req domain.CommentCreate
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if err := validator.Validate(&req); err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	comment, err := h.userBookService.CreateComment(userID, bookID, req.Content)
	if err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	utils.JSONResponse(w, http.StatusCreated, comment)
}

func (h *UserBookHandler) GetBookComments(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	bookID, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, "invalid book ID")
		return
	}

	comments, err := h.userBookService.GetBookComments(bookID)
	if err != nil {
		utils.ErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SuccessResponseWithData(w, comments)
}

func (h *UserBookHandler) UpdateComment(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())
	vars := mux.Vars(r)
	commentID, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, "invalid comment ID")
		return
	}

	var req domain.CommentCreate
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if err := validator.Validate(&req); err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	err = h.userBookService.UpdateComment(commentID, userID, req.Content)
	if err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	utils.SuccessResponseWithMessage(w, "Comment updated successfully")
}

func (h *UserBookHandler) DeleteComment(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())
	vars := mux.Vars(r)
	commentID, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, "invalid comment ID")
		return
	}

	err = h.userBookService.DeleteComment(commentID, userID)
	if err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	utils.SuccessResponseWithMessage(w, "Comment deleted successfully")
}
