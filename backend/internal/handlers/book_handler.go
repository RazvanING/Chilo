package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/razvan/library-app/internal/domain"
	"github.com/razvan/library-app/internal/service"
	"github.com/razvan/library-app/internal/utils"
	"github.com/razvan/library-app/pkg/validator"
)

type BookHandler struct {
	bookService *service.BookService
	uploadDir   string
}

func NewBookHandler(bookService *service.BookService, uploadDir string) *BookHandler {
	return &BookHandler{
		bookService: bookService,
		uploadDir:   uploadDir,
	}
}

func (h *BookHandler) CreateBook(w http.ResponseWriter, r *http.Request) {
	var req domain.BookCreate
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if err := validator.Validate(&req); err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	book, err := h.bookService.CreateBook(&req)
	if err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	utils.JSONResponse(w, http.StatusCreated, book)
}

func (h *BookHandler) GetBook(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, "invalid book ID")
		return
	}

	book, err := h.bookService.GetBook(id)
	if err != nil {
		utils.ErrorResponse(w, http.StatusNotFound, err.Error())
		return
	}

	utils.SuccessResponseWithData(w, book)
}

func (h *BookHandler) GetAllBooks(w http.ResponseWriter, r *http.Request) {
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	if page < 1 {
		page = 1
	}

	pageSize, _ := strconv.Atoi(r.URL.Query().Get("page_size"))
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	books, err := h.bookService.GetAllBooks(page, pageSize)
	if err != nil {
		utils.ErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SuccessResponseWithData(w, books)
}

func (h *BookHandler) SearchBooks(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("q")
	if query == "" {
		utils.ErrorResponse(w, http.StatusBadRequest, "search query is required")
		return
	}

	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	if page < 1 {
		page = 1
	}

	pageSize, _ := strconv.Atoi(r.URL.Query().Get("page_size"))
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	books, err := h.bookService.SearchBooks(query, page, pageSize)
	if err != nil {
		utils.ErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SuccessResponseWithData(w, books)
}

func (h *BookHandler) UpdateBook(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, "invalid book ID")
		return
	}

	var req domain.BookUpdate
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if err := validator.Validate(&req); err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	book, err := h.bookService.UpdateBook(id, &req)
	if err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	utils.SuccessResponseWithData(w, book)
}

func (h *BookHandler) DeleteBook(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, "invalid book ID")
		return
	}

	err = h.bookService.DeleteBook(id)
	if err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	utils.SuccessResponseWithMessage(w, "Book deleted successfully")
}

func (h *BookHandler) UploadCover(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	bookID, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, "invalid book ID")
		return
	}

	// Parse multipart form (10 MB max)
	err = r.ParseMultipartForm(10 << 20)
	if err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, "file too large")
		return
	}

	file, header, err := r.FormFile("cover")
	if err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, "cover file is required")
		return
	}
	defer file.Close()

	// Save file
	filename, err := utils.SaveUploadedFile(file, header, h.uploadDir)
	if err != nil {
		utils.ErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	// Update book cover URL
	coverURL := "/uploads/" + filename
	err = h.bookService.UpdateBookCover(bookID, coverURL)
	if err != nil {
		utils.ErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SuccessResponseWithData(w, map[string]string{
		"cover_url": coverURL,
	})
}

// Author handlers
func (h *BookHandler) CreateAuthor(w http.ResponseWriter, r *http.Request) {
	var req domain.AuthorCreate
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if err := validator.Validate(&req); err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	author, err := h.bookService.CreateAuthor(&req)
	if err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	utils.JSONResponse(w, http.StatusCreated, author)
}

func (h *BookHandler) GetAllAuthors(w http.ResponseWriter, r *http.Request) {
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	if page < 1 {
		page = 1
	}

	pageSize, _ := strconv.Atoi(r.URL.Query().Get("page_size"))
	if pageSize < 1 || pageSize > 100 {
		pageSize = 50
	}

	authors, err := h.bookService.GetAllAuthors(page, pageSize)
	if err != nil {
		utils.ErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SuccessResponseWithData(w, authors)
}

func (h *BookHandler) GetAuthor(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, "invalid author ID")
		return
	}

	author, err := h.bookService.GetAuthor(id)
	if err != nil {
		utils.ErrorResponse(w, http.StatusNotFound, err.Error())
		return
	}

	utils.SuccessResponseWithData(w, author)
}

func (h *BookHandler) UpdateAuthor(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, "invalid author ID")
		return
	}

	var req domain.AuthorCreate
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if err := validator.Validate(&req); err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	author, err := h.bookService.UpdateAuthor(id, &req)
	if err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	utils.SuccessResponseWithData(w, author)
}

func (h *BookHandler) DeleteAuthor(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, "invalid author ID")
		return
	}

	err = h.bookService.DeleteAuthor(id)
	if err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	utils.SuccessResponseWithMessage(w, "Author deleted successfully")
}
