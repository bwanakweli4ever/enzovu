# Enzovu - Elegant Go Web Framework

**Enzovu** is a robust and elegant Go framework inspired by the strength and wisdom of elephants. Built with Laravel's architectural philosophy in mind, it brings familiar MVC patterns to Go while maintaining the language's performance and simplicity.

---

## ✨ Features

- **🏗️ MVC Architecture** - Clean separation of concerns with Models, Views, and Controllers
- **🔥 Hot Reload** - Built-in file watching for instant development feedback
- **🚀 Smart Routing** - Intelligent request handling with parameter extraction
- **🛡️ Memory-Safe** - Built with Go's strong type system and safety guarantees
- **📦 Scalable by Design** - From prototypes to enterprise applications
- **🎨 CLI Code Generation** - Laravel-inspired `go-craft` tool for rapid development
- **🔌 Middleware Support** - Flexible request/response processing
- **💾 Database Ready** - Built-in support for MySQL, PostgreSQL, and SQLite
- **⚙️ Environment Configuration** - Flexible configuration management

---

## 🚀 Quick Start

### Prerequisites

- **Go 1.21+** - [Download Go](https://golang.org/dl/)
- **Database** (optional) - MySQL, PostgreSQL, or SQLite

### Installation

```bash
# Clone the repository
git clone https://github.com/yourusername/enzovu.git
cd enzovu

# Initialize Go module
go mod init enzovu

# Install dependencies
go mod tidy

# Copy environment configuration
cp .env.example .env

# Start the development server with hot reload
go run main.go
```

Visit `http://localhost:8000` and you'll see Enzovu running! 🎉

---

## 🔥 Hot Reload Development

Enzovu includes built-in hot reload - no external tools needed!

```bash
# Start development server (automatically watches for changes)
go run main.go

# Edit any .go file and save
# Changes are instantly reflected - no restart needed!
```

**Features:**
- ✅ Instant route reloading
- ✅ No server restarts
- ✅ Same port always (no conflicts)
- ✅ Only enabled in development mode

---

## 🛠️ CLI Tool - go-craft

Generate code quickly with the built-in `go-craft` CLI tool:

```bash
# Create a new model
go run cmd/go-craft.go create model Post

# Create a new controller
go run cmd/go-craft.go create controller Post

# Create middleware
go run cmd/go-craft.go create middleware Auth

# Create database migration
go run cmd/go-craft.go create migration create_posts_table
```

---

## 📁 Project Structure

```
enzovu/
├── app/
│   ├── Http/
│   │   ├── Controllers/     # Request handlers
│   │   └── Middleware/      # HTTP middleware
│   ├── Models/              # Data models
│   └── commands/            # CLI commands
├── bootstrap/               # App initialization
├── config/                  # Configuration files
├── database/
│   ├── migrations/          # Database migrations
│   └── seeds/               # Database seeders
├── public/                  # Static assets (CSS, JS, images)
├── resources/views/         # Templates
├── routes/                  # Route definitions
├── cmd/                     # CLI tools
├── .env                     # Environment variables
└── main.go                  # Application entry point
```

---

## 🎯 Building Your First Feature

### 1. Create a Model
```bash
go run cmd/go-craft.go create model Article
```

This creates `app/Models/article.go` with:
```go
type Article struct {
    ID        int       `json:"id" db:"id"`
    Title     string    `json:"title" db:"title"`
    Content   string    `json:"content" db:"content"`
    CreatedAt time.Time `json:"created_at" db:"created_at"`
    UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}
```

### 2. Create a Controller
```bash
go run cmd/go-craft.go create controller Article
```

This creates `app/Http/Controllers/article_controller.go` with full CRUD operations:
- `Index()` - List all articles
- `Show()` - Show specific article
- `Create()` - Create new article
- `Update()` - Update article
- `Delete()` - Delete article

### 3. Add Routes
Edit `routes/web.go`:
```go
func SetupRoutes() http.Handler {
    router := NewRouter()
    
    // Article routes
    articleController := &controllers.ArticleController{}
    router.GET("/api/articles", articleController.Index)
    router.POST("/api/articles", articleController.Create)
    router.GET("/api/articles/{id}", articleController.Show)
    router.PUT("/api/articles/{id}", articleController.Update)
    router.DELETE("/api/articles/{id}", articleController.Delete)
    
    return router
}
```

### 4. Test Your API
```bash
# Get all articles
curl http://localhost:8000/api/articles

# Create an article
curl -X POST http://localhost:8000/api/articles \
  -H "Content-Type: application/json" \
  -d '{"title":"Hello World","content":"My first article"}'
```

---

## ⚙️ Configuration

### Environment Variables

Create `.env` file:
```env
# Application
APP_NAME="My Enzovu App"
APP_ENV=development
APP_PORT=8000
APP_DEBUG=true

# Database
DB_DRIVER=mysql
DB_HOST=localhost
DB_PORT=3306
DB_USER=root
DB_PASSWORD=your_password
DB_DATABASE=enzovu_db
```

### Available Environments
- `development` - Enables hot reload, debug logging
- `production` - Optimized for performance, no hot reload

---

## 🗄️ Database Integration

### Supported Databases
- **MySQL** - `go get github.com/go-sql-driver/mysql`
- **PostgreSQL** - `go get github.com/lib/pq`
- **SQLite** - `go get github.com/mattn/go-sqlite3`

### Database Configuration
```env
# MySQL
DB_DRIVER=mysql
DB_HOST=localhost
DB_PORT=3306
DB_USER=root
DB_PASSWORD=password
DB_DATABASE=enzovu_db

# PostgreSQL
DB_DRIVER=postgres
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=password
DB_DATABASE=enzovu_db

# SQLite
DB_DRIVER=sqlite3
DB_DATABASE=database.db
```

### Migrations

Create and run migrations:
```bash
# Create migration
go run cmd/go-craft.go create migration create_users_table

# Migrations run automatically in development
go run main.go
```

---

## 🚦 Routing

### Basic Routes
```go
router.GET("/", homeHandler)
router.POST("/api/users", createUser)
router.PUT("/api/users/{id}", updateUser)
router.DELETE("/api/users/{id}", deleteUser)
```

### Route Parameters
```go
func getUserHandler(w http.ResponseWriter, r *http.Request) {
    id := routes.GetParam(r, "id")
    // Use the ID parameter
}
```

### Route Groups
```go
api := router.Group("/api")
{
    api.GET("/users", getUsersHandler)
    api.POST("/users", createUserHandler)
}
```

### Middleware
```go
// Global middleware
router.Use(middleware.LoggingMiddleware)

// Route-specific middleware
router.GET("/admin", adminHandler, middleware.AuthMiddleware)

// Group middleware
protected := router.Group("/admin", middleware.AuthMiddleware)
```

---

## 🔌 Middleware

### Built-in Middleware
- `LoggingMiddleware` - Request logging with colors in development
- `AuthMiddleware` - Authentication example

### Creating Custom Middleware
```bash
go run cmd/go-craft.go create middleware CORS
```

Example middleware:
```go
func CORSMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Access-Control-Allow-Origin", "*")
        next.ServeHTTP(w, r)
    })
}
```

---

## 📊 API Examples

### JSON API Response
```go
func getUsersAPI(w http.ResponseWriter, r *http.Request) {
    users := models.GetAllUsers()
    
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(map[string]interface{}{
        "data": users,
        "count": len(users),
        "status": "success",
    })
}
```

### Error Handling
```go
func createUser(w http.ResponseWriter, r *http.Request) {
    var user models.User
    if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
        http.Error(w, "Invalid JSON", http.StatusBadRequest)
        return
    }
    
    if errors := user.Validate(); len(errors) > 0 {
        w.Header().Set("Content-Type", "application/json")
        w.WriteHeader(http.StatusBadRequest)
        json.NewEncoder(w).Encode(map[string]interface{}{
            "errors": errors,
        })
        return
    }
    
    // Save user...
}
```

---

## 🏗️ Production Deployment

### Build for Production
```bash
# Build binary
go build -o enzovu main.go

# Run in production mode
APP_ENV=production ./enzovu
```

### Environment Setup
```bash
# Set production environment variables
export APP_ENV=production
export APP_PORT=8080
export DB_HOST=your-db-host
export DB_PASSWORD=your-secure-password

# Run the application
./enzovu
```

### Docker Support
```dockerfile
FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY . .
RUN go mod download
RUN go build -o enzovu main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/enzovu .
COPY --from=builder /app/public ./public
CMD ["./enzovu"]
```

---

## 🤝 Contributing

We welcome contributions! Please see [CONTRIBUTING.md](CONTRIBUTING.md) for guidelines.

### Development Setup
```bash
# Fork and clone the repo
git clone https://github.com/yourusername/enzovu.git
cd enzovu

# Create a feature branch
git checkout -b feature/amazing-feature

# Make your changes and test
go run main.go

# Submit a pull request
```

---

## 📚 Documentation

- **[API Reference](docs/api.md)** - Complete API documentation
- **[Middleware Guide](docs/middleware.md)** - Creating custom middleware
- **[Database Guide](docs/database.md)** - Working with databases
- **[Deployment Guide](docs/deployment.md)** - Production deployment

---

## 🎯 Roadmap

- [ ] **Authentication System** - Built-in user authentication
- [ ] **WebSocket Support** - Real-time communication
- [ ] **Template Engine** - Server-side rendering
- [ ] **Cache Layer** - Redis/Memory caching
- [ ] **Queue System** - Background job processing
- [ ] **Testing Framework** - Built-in testing utilities
- [ ] **Plugin System** - Extensible architecture
- [ ] **Admin Panel** - Web-based administration

---

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

---

## 🙏 Acknowledgments

- Inspired by [Laravel](https://laravel.com/) for its elegant architecture
- Built with [Go](https://golang.org/) for performance and reliability
- Logo elephant represents strength, wisdom, and memory

---

## 💬 Support

- **GitHub Issues** - [Report bugs or request features](https://github.com/bwanakweli4ever/enzovu/issues)
- **Discussions** - [Ask questions and share ideas](https://github.com/bwanakweli4ever/enzovu/discussions)
- **Email** - your-email@example.com

---

<p align="center">
  <strong>Made with ❤️ and 🐘 by the Enzovu Team</strong>
</p>