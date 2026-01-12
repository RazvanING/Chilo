# Library App

A full-stack library management application built with Go, PostgreSQL, and Vue.js. Features include user authentication with 2FA, book management, reading lists, favorites, and commenting system.

## Features

### For Users
- 📚 Browse and search books
- 🔐 Secure authentication (Email/Password + Google OAuth)
- 🔒 Two-Factor Authentication (2FA) with TOTP
- ❤️ Favorite books
- 📖 Reading lists (Want to Read, Currently Reading, Read)
- 💬 Comment on books
- 👤 User profile management

### For Admins
- ➕ Add/Edit/Delete books
- ✍️ Manage authors
- 📸 Upload book covers
- 👥 Promote users to admin
- 📊 Admin dashboard

## Technology Stack

### Backend
- **Go 1.21+** - Backend language
- **Gorilla Mux** - HTTP router
- **PostgreSQL** - Database
- **JWT** - Authentication
- **bcrypt** - Password hashing
- **TOTP** - Two-factor authentication

### Frontend
- **Vue 3** - Frontend framework
- **Vite** - Build tool
- **Pinia** - State management
- **Vue Router** - Routing
- **TailwindCSS** - Styling
- **Axios** - HTTP client

### Infrastructure
- **Docker** - Containerization
- **Docker Compose** - Multi-container orchestration

## Architecture

The application follows **Clean Architecture** principles:

```
backend/
├── cmd/api/              # Application entry point
├── internal/
│   ├── domain/           # Business entities
│   ├── handlers/         # HTTP handlers (controllers)
│   ├── middleware/       # Middleware (auth, CORS, logging)
│   ├── repository/       # Data access layer
│   ├── service/          # Business logic
│   └── utils/            # Utility functions
├── pkg/
│   ├── auth/             # Authentication utilities
│   ├── database/         # Database connection
│   └── validator/        # Input validation
└── migrations/           # SQL migrations

frontend/
├── src/
│   ├── components/       # Reusable Vue components
│   ├── views/            # Page components
│   ├── router/           # Vue Router configuration
│   ├── stores/           # Pinia stores
│   ├── services/         # API service layer
│   └── assets/           # Static assets
```

## Getting Started

### Prerequisites
- Docker and Docker Compose
- Go 1.21+ (for local development)
- Node.js 20+ (for local development)

### Quick Start with Docker

1. **Clone the repository**
   ```bash
   git clone <repository-url>
   cd Chilo
   ```

2. **Create environment files**
   ```bash
   cp .env.example backend/.env
   cp frontend/.env.example frontend/.env
   ```

3. **Update environment variables**
   Edit `backend/.env` and set your configuration:
   ```env
   DB_HOST=postgres
   DB_PORT=5432
   DB_USER=libraryuser
   DB_PASSWORD=librarypass
   DB_NAME=librarydb
   JWT_SECRET=your-super-secret-jwt-key-change-in-production
   PORT=8080
   FRONTEND_URL=http://localhost:3000
   GOOGLE_CLIENT_ID=your-google-client-id
   GOOGLE_CLIENT_SECRET=your-google-client-secret
   ```

4. **Start the application**
   ```bash
   docker-compose up -d
   ```

5. **Access the application**
   - Frontend: http://localhost:3000
   - Backend API: http://localhost:8080
   - Health Check: http://localhost:8080/api/health

6. **Default Admin Account**
   - Email: `admin@library.com`
   - Password: `admin123`
   - ⚠️ **Change this password immediately in production!**

### Local Development Setup

#### Backend

1. **Install dependencies**
   ```bash
   cd backend
   go mod download
   ```

2. **Set up PostgreSQL**
   ```bash
   docker run -d \
     --name library_postgres \
     -e POSTGRES_USER=libraryuser \
     -e POSTGRES_PASSWORD=librarypass \
     -e POSTGRES_DB=librarydb \
     -p 5432:5432 \
     postgres:15-alpine
   ```

3. **Run migrations**
   ```bash
   psql -h localhost -U libraryuser -d librarydb -f migrations/001_init_schema.sql
   ```

4. **Start the backend**
   ```bash
   go run cmd/api/main.go
   ```

#### Frontend

1. **Install dependencies**
   ```bash
   cd frontend
   npm install
   ```

2. **Start the development server**
   ```bash
   npm run dev
   ```

## API Documentation

### Authentication Endpoints

#### Register
```http
POST /api/auth/register
Content-Type: application/json

{
  "email": "user@example.com",
  "username": "johndoe",
  "password": "password123"
}
```

#### Login
```http
POST /api/auth/login
Content-Type: application/json

{
  "email": "user@example.com",
  "password": "password123"
}
```

#### Setup 2FA
```http
POST /api/auth/2fa/setup
Authorization: Bearer <token>
```

#### Verify 2FA
```http
POST /api/auth/2fa/verify
Authorization: Bearer <token>
Content-Type: application/json

{
  "code": "123456"
}
```

### Book Endpoints

#### Get All Books
```http
GET /api/books?page=1&page_size=20
```

#### Search Books
```http
GET /api/books/search?q=searchterm&page=1&page_size=20
```

#### Get Book Details
```http
GET /api/books/:id
```

#### Create Book (Admin)
```http
POST /api/books
Authorization: Bearer <token>
Content-Type: application/json

{
  "title": "Book Title",
  "description": "Book description",
  "isbn": "1234567890123",
  "published_at": "2024-01-01",
  "author_ids": [1, 2]
}
```

#### Upload Book Cover (Admin)
```http
POST /api/books/:id/cover
Authorization: Bearer <token>
Content-Type: multipart/form-data

cover: <file>
```

### User Book Endpoints

#### Add to Reading List
```http
POST /api/user/books/:id/reading-list
Authorization: Bearer <token>
Content-Type: application/json

{
  "status": "want_to_read"  // or "reading" or "read"
}
```

#### Add to Favorites
```http
POST /api/user/books/:id/favorites
Authorization: Bearer <token>
```

#### Add Comment
```http
POST /api/books/:id/comments
Authorization: Bearer <token>
Content-Type: application/json

{
  "content": "Great book!"
}
```

## Security Features

### Authentication & Authorization
- **JWT tokens** with 24-hour expiration
- **Refresh tokens** with 7-day expiration
- **bcrypt** password hashing with cost factor 12
- **TOTP-based 2FA** for enhanced security
- **Role-based access control** (User/Admin)

### Security Headers
- X-Content-Type-Options: nosniff
- X-Frame-Options: DENY
- X-XSS-Protection: 1; mode=block
- Strict-Transport-Security
- Content-Security-Policy
- Referrer-Policy

### Input Validation
- Request body validation
- SQL injection prevention (parameterized queries)
- File upload validation (type and size)
- XSS protection

### CORS Configuration
- Configurable allowed origins
- Credentials support
- Preflight request handling

## Database Schema

```sql
Users (id, email, username, password_hash, is_admin, google_id,
       two_factor_secret, two_factor_enabled, email_verified)

Books (id, title, description, cover_url, isbn, published_at)

Authors (id, name, bio)

Book_Authors (book_id, author_id)  -- Many-to-many relationship

User_Books (id, user_id, book_id, status)  -- Reading lists

Favorites (id, user_id, book_id)

Comments (id, user_id, book_id, content)
```

## Environment Variables

### Backend (.env)
```env
# Database
DB_HOST=localhost
DB_PORT=5432
DB_USER=libraryuser
DB_PASSWORD=librarypass
DB_NAME=librarydb

# JWT
JWT_SECRET=your-secret-key
JWT_EXPIRATION=24h

# Server
PORT=8080
FRONTEND_URL=http://localhost:3000

# Google OAuth (Optional)
GOOGLE_CLIENT_ID=your-client-id
GOOGLE_CLIENT_SECRET=your-client-secret
GOOGLE_REDIRECT_URL=http://localhost:8080/api/auth/google/callback

# File Upload
MAX_UPLOAD_SIZE=10485760
UPLOAD_PATH=./uploads

# 2FA
APP_NAME=LibraryApp

# CORS
ALLOWED_ORIGINS=http://localhost:3000
```

### Frontend (.env)
```env
VITE_API_URL=http://localhost:8080/api
```

## Testing

### Manual Testing
1. Register a new user
2. Enable 2FA in profile settings
3. Browse and search books
4. Add books to reading lists and favorites
5. Comment on books
6. Login as admin to manage books and authors

### Admin Features Testing
1. Login with admin credentials
2. Create new books and authors
3. Upload book covers
4. Edit and delete books
5. Promote users to admin

## Deployment

### Production Considerations

1. **Security**
   - Change default admin password
   - Use strong JWT secret
   - Enable HTTPS
   - Configure proper CORS origins
   - Set up rate limiting

2. **Database**
   - Use managed PostgreSQL service
   - Enable backups
   - Set up replication

3. **File Storage**
   - Use object storage (S3, GCS) for book covers
   - Configure CDN for static assets

4. **Monitoring**
   - Set up logging
   - Configure error tracking
   - Monitor performance metrics

## Troubleshooting

### Database connection fails
- Ensure PostgreSQL is running
- Check database credentials
- Verify network connectivity

### Frontend cannot connect to backend
- Check CORS configuration
- Verify API URL in frontend .env
- Check backend is running on correct port

### 2FA not working
- Ensure time is synchronized on server
- Verify TOTP secret is properly stored
- Check authenticator app time sync

## Contributing

1. Fork the repository
2. Create a feature branch
3. Commit your changes
4. Push to the branch
5. Create a Pull Request

## License

MIT License - feel free to use this project for learning or production.

## Support

For issues and questions, please open an issue on GitHub.
