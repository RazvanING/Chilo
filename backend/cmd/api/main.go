package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/razvan/library-app/internal/handlers"
	"github.com/razvan/library-app/internal/middleware"
	"github.com/razvan/library-app/internal/repository"
	"github.com/razvan/library-app/internal/service"
	"github.com/razvan/library-app/pkg/database"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	// Database configuration
	dbConfig := database.Config{
		Host:     getEnv("DB_HOST", "localhost"),
		Port:     getEnv("DB_PORT", "5432"),
		User:     getEnv("DB_USER", "libraryuser"),
		Password: getEnv("DB_PASSWORD", "librarypass"),
		DBName:   getEnv("DB_NAME", "librarydb"),
		SSLMode:  getEnv("DB_SSLMODE", "disable"),
	}

	// Connect to database
	db, err := database.NewPostgresDB(dbConfig)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	log.Println("Successfully connected to database")

	// Initialize repositories
	userRepo := repository.NewUserRepository(db)
	bookRepo := repository.NewBookRepository(db)
	userBookRepo := repository.NewUserBookRepository(db)

	// Initialize services
	jwtSecret := getEnv("JWT_SECRET", "your-secret-key-change-in-production")
	appName := getEnv("APP_NAME", "LibraryApp")
	authService := service.NewAuthService(userRepo, jwtSecret, appName)
	bookService := service.NewBookService(bookRepo)
	userBookService := service.NewUserBookService(userBookRepo, bookRepo)

	// Initialize handlers
	uploadDir := getEnv("UPLOAD_PATH", "./uploads")
	authHandler := handlers.NewAuthHandler(authService)
	bookHandler := handlers.NewBookHandler(bookService, uploadDir)
	userBookHandler := handlers.NewUserBookHandler(userBookService)

	// Initialize middleware
	authMiddleware := middleware.NewAuthMiddleware(jwtSecret)
	corsMiddleware := middleware.NewCORSMiddleware(getEnv("ALLOWED_ORIGINS", "http://localhost:3000"))

	// Setup router
	r := mux.NewRouter()

	// Apply global middleware
	r.Use(middleware.Logger)
	r.Use(corsMiddleware.Handler)
	r.Use(middleware.SecurityHeaders)

	// API routes
	api := r.PathPrefix("/api").Subrouter()

	// Public routes
	api.HandleFunc("/health", healthCheck).Methods("GET")

	// Auth routes
	auth := api.PathPrefix("/auth").Subrouter()
	auth.HandleFunc("/register", authHandler.Register).Methods("POST")
	auth.HandleFunc("/login", authHandler.Login).Methods("POST")
	auth.HandleFunc("/2fa/verify-login", authHandler.VerifyTwoFactorLogin).Methods("POST")

	// Protected auth routes
	authProtected := auth.PathPrefix("").Subrouter()
	authProtected.Use(authMiddleware.Authenticate)
	authProtected.HandleFunc("/me", authHandler.GetMe).Methods("GET")
	authProtected.HandleFunc("/2fa/setup", authHandler.SetupTwoFactor).Methods("POST")
	authProtected.HandleFunc("/2fa/verify", authHandler.VerifyTwoFactorSetup).Methods("POST")
	authProtected.HandleFunc("/2fa/disable", authHandler.DisableTwoFactor).Methods("POST")

	// Admin only routes
	authAdmin := auth.PathPrefix("").Subrouter()
	authAdmin.Use(authMiddleware.Authenticate)
	authAdmin.Use(authMiddleware.RequireAdmin)
	authAdmin.HandleFunc("/users/{id}/make-admin", authHandler.MakeAdmin).Methods("POST")

	// Book routes (public read, admin write)
	books := api.PathPrefix("/books").Subrouter()
	books.HandleFunc("", bookHandler.GetAllBooks).Methods("GET")
	books.HandleFunc("/search", bookHandler.SearchBooks).Methods("GET")
	books.HandleFunc("/{id}", bookHandler.GetBook).Methods("GET")

	// Protected book routes (admin only)
	booksAdmin := books.PathPrefix("").Subrouter()
	booksAdmin.Use(authMiddleware.Authenticate)
	booksAdmin.Use(authMiddleware.RequireAdmin)
	booksAdmin.HandleFunc("", bookHandler.CreateBook).Methods("POST")
	booksAdmin.HandleFunc("/{id}", bookHandler.UpdateBook).Methods("PUT")
	booksAdmin.HandleFunc("/{id}", bookHandler.DeleteBook).Methods("DELETE")
	booksAdmin.HandleFunc("/{id}/cover", bookHandler.UploadCover).Methods("POST")

	// Author routes
	authors := api.PathPrefix("/authors").Subrouter()
	authors.HandleFunc("", bookHandler.GetAllAuthors).Methods("GET")
	authors.HandleFunc("/{id}", bookHandler.GetAuthor).Methods("GET")

	// Protected author routes (admin only)
	authorsAdmin := authors.PathPrefix("").Subrouter()
	authorsAdmin.Use(authMiddleware.Authenticate)
	authorsAdmin.Use(authMiddleware.RequireAdmin)
	authorsAdmin.HandleFunc("", bookHandler.CreateAuthor).Methods("POST")
	authorsAdmin.HandleFunc("/{id}", bookHandler.UpdateAuthor).Methods("PUT")
	authorsAdmin.HandleFunc("/{id}", bookHandler.DeleteAuthor).Methods("DELETE")

	// User book routes (protected)
	userBooks := api.PathPrefix("/user").Subrouter()
	userBooks.Use(authMiddleware.Authenticate)

	// Reading lists
	userBooks.HandleFunc("/reading-list", userBookHandler.GetReadingList).Methods("GET")
	userBooks.HandleFunc("/books/{id}/reading-list", userBookHandler.AddToReadingList).Methods("POST")
	userBooks.HandleFunc("/books/{id}/reading-list", userBookHandler.RemoveFromReadingList).Methods("DELETE")

	// Favorites
	userBooks.HandleFunc("/favorites", userBookHandler.GetFavorites).Methods("GET")
	userBooks.HandleFunc("/books/{id}/favorites", userBookHandler.AddToFavorites).Methods("POST")
	userBooks.HandleFunc("/books/{id}/favorites", userBookHandler.RemoveFromFavorites).Methods("DELETE")

	// Comments
	comments := api.PathPrefix("/books/{id}/comments").Subrouter()
	comments.HandleFunc("", userBookHandler.GetBookComments).Methods("GET")

	commentsProtected := comments.PathPrefix("").Subrouter()
	commentsProtected.Use(authMiddleware.Authenticate)
	commentsProtected.HandleFunc("", userBookHandler.CreateComment).Methods("POST")

	api.HandleFunc("/comments/{id}", userBookHandler.UpdateComment).
		Methods("PUT").
		Handler(authMiddleware.Authenticate(http.HandlerFunc(userBookHandler.UpdateComment)))
	api.HandleFunc("/comments/{id}", userBookHandler.DeleteComment).
		Methods("DELETE").
		Handler(authMiddleware.Authenticate(http.HandlerFunc(userBookHandler.DeleteComment)))

	// Serve uploaded files
	r.PathPrefix("/uploads/").Handler(
		http.StripPrefix("/uploads/", http.FileServer(http.Dir(uploadDir))),
	)

	// Start server
	port := getEnv("PORT", "8080")
	srv := &http.Server{
		Addr:         ":" + port,
		Handler:      r,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	log.Printf("Server starting on port %s", port)
	if err := srv.ListenAndServe(); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}

func healthCheck(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"ok","message":"Library API is running"}`))
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
