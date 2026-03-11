package main

import (
	"log"

	"github.com/RajendraMahadana/perpustakaan-clean/internal/infrastructure/config"
	"github.com/RajendraMahadana/perpustakaan-clean/internal/infrastructure/database"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	// Auth imports
	authHttp "github.com/RajendraMahadana/perpustakaan-clean/internal/delivery/http"
	authRepoImpl "github.com/RajendraMahadana/perpustakaan-clean/internal/infrastructure/repository"
	authRepo "github.com/RajendraMahadana/perpustakaan-clean/internal/repository"
	authUsecase "github.com/RajendraMahadana/perpustakaan-clean/internal/usecase"

	// Book imports
	bookHttp "github.com/RajendraMahadana/perpustakaan-clean/internal/delivery/http"
	bookRepoImpl "github.com/RajendraMahadana/perpustakaan-clean/internal/infrastructure/repository"
	bookRepo "github.com/RajendraMahadana/perpustakaan-clean/internal/repository"
	bookUsecase "github.com/RajendraMahadana/perpustakaan-clean/internal/usecase"

	// Middleware
	"github.com/RajendraMahadana/perpustakaan-clean/pkg/middleware"

	// Utils
	"github.com/RajendraMahadana/perpustakaan-clean/internal/utils"
)

func main() {
	// Load env
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Load Config
	cfg, err := config.LoadConfig() // Ubah nama variabel jadi cfg
	if err != nil {
		log.Fatal("Error loading config: ", err)
	}

	// Inisialisasi JWT dengan config
	db, err := database.NewMySQLConnectionWithConfig(cfg)
	if err != nil {
		log.Fatal("Failed to connect database:", err)
	}

	// Connect Database
	database.AutoMigrate(db)

	// Gunakan config untuk inisialisasi JWT
	utils.InitJWTWithConfig(cfg)

	// ==================== INITIALIZE REPOSITORIES ====================
	// Book Repository
	var bookRepo bookRepo.BookRepository = bookRepoImpl.NewBookRepository(db)

	// User Repository (untuk auth)
	var userRepo authRepo.UserRepository = authRepoImpl.NewUserRepository(db)

	// ==================== INITIALIZE USECASES ====================
	// Book Usecase
	bookUsecase := bookUsecase.NewBookUsecase(bookRepo)

	// Auth Usecase
	authUsecase := authUsecase.NewAuthUsecase(userRepo)

	// ==================== INITIALIZE HANDLERS ====================
	// Book Handler
	bookHandler := bookHttp.NewBookHandler(bookUsecase)

	// Auth Handler
	authHandler := authHttp.NewAuthHandler(authUsecase)

	// ==================== INITIALIZE MIDDLEWARE ====================
	authMiddleware := middleware.NewAuthMiddleware()

	// ==================== INIT GIN ====================
	r := gin.Default()

	// ==================== PUBLIC ROUTES ====================
	auth := r.Group("/api/v1/auth")
	{
		auth.POST("/register", authHandler.Register)
		auth.POST("/login", authHandler.Login)
	}

	// ==================== PROTECTED ROUTES ====================
	protected := r.Group("/api/v1")
	protected.Use(authMiddleware.Authenticate())
	{
		// Auth protected routes
		authProtected := protected.Group("/auth")
		{
			authProtected.POST("/refresh", authHandler.RefreshToken)
			authProtected.POST("/logout", authHandler.Logout)
			authProtected.GET("/profile", authHandler.GetProfile)
		}

		// Book routes (semua user yang sudah login)
		books := protected.Group("/books")
		{
			books.GET("/", bookHandler.GetAllBooks)
			books.GET("/limit", bookHandler.GetBook)
			books.GET("/:id", bookHandler.GetByID)
		}

		// Admin only routes
		admin := protected.Group("/admin")
		admin.Use(authMiddleware.Authorize("admin"))
		{
			admin.POST("/books", bookHandler.Create)
			admin.PUT("/books/:id", bookHandler.Update)
			admin.DELETE("/books/:id", bookHandler.Delete)
		}
	}

	// ==================== RUN SERVER ====================
	port := cfg.AppPort // Gunakan port dari config
	if port == "" {
		port = "8080"
	}

	log.Printf("Server running on port %s", port)
	r.Run(":" + port)
}
