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

	// 🔴 BARU: Borrowing imports
	borrowingHttp "github.com/RajendraMahadana/perpustakaan-clean/internal/delivery/http"
	borrowingRepoImpl "github.com/RajendraMahadana/perpustakaan-clean/internal/infrastructure/repository"
	borrowingRepo "github.com/RajendraMahadana/perpustakaan-clean/internal/repository"
	borrowingUsecase "github.com/RajendraMahadana/perpustakaan-clean/internal/usecase"

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
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal("Error loading config: ", err)
	}

	// Connect Database
	db, err := database.NewMySQLConnectionWithConfig(cfg)
	if err != nil {
		log.Fatal("Failed to connect database:", err)
	}

	// Auto Migrate
	database.AutoMigrate(db)

	// Inisialisasi JWT dengan config
	utils.InitJWTWithConfig(cfg)

	// ==================== INITIALIZE REPOSITORIES ====================
	// Book Repository
	var bookRepo bookRepo.BookRepository = bookRepoImpl.NewBookRepository(db)

	// User Repository (untuk auth)
	var userRepo authRepo.UserRepository = authRepoImpl.NewUserRepository(db)

	// 🔴 BARU: Borrowing Repository
	var borrowingRepo borrowingRepo.BorrowingRepository = borrowingRepoImpl.NewBorrowingRepository(db)

	// ==================== INITIALIZE USECASES ====================
	// Book Usecase
	bookUsecase := bookUsecase.NewBookUsecase(bookRepo)

	// Auth Usecase
	authUsecase := authUsecase.NewAuthUsecase(userRepo)

	// 🔴 BARU: Borrowing Usecase (TANPA TRANSACTION MANAGER)
	borrowingUsecase := borrowingUsecase.NewBorrowingUsecase(
		borrowingRepo,
		bookRepo,
		userRepo,
		// TIDAK PERLU TX MANAGER
	)

	// ==================== INITIALIZE HANDLERS ====================
	// Book Handler
	bookHandler := bookHttp.NewBookHandler(bookUsecase)

	// Auth Handler
	authHandler := authHttp.NewAuthHandler(authUsecase)

	// 🔴 BARU: Borrowing Handler
	borrowingHandler := borrowingHttp.NewBorrowingHandler(borrowingUsecase)

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

		// 🔴 BARU: Borrowing routes (semua user yang sudah login)
		borrowings := protected.Group("/borrowings")
		{
			borrowings.POST("/", borrowingHandler.PinjamBuku)
			borrowings.POST("/return", borrowingHandler.KembalikanBuku)
			borrowings.GET("/my", borrowingHandler.GetMyBorrowings)
			borrowings.GET("/active", borrowingHandler.GetActiveBorrowings)
			borrowings.GET("/:id", borrowingHandler.GetBorrowingByID)
		}

		admin := protected.Group("/admin")
		admin.Use(authMiddleware.Authorize("admin"))
		{
			admin.POST("/books", bookHandler.Create)
			admin.PUT("/books/:id", bookHandler.Update)
			admin.DELETE("/books/:id", bookHandler.Delete)

			admin.GET("/borrowings", borrowingHandler.GetAllBorrowings)
		}
	}

	// ==================== RUN SERVER ====================
	port := cfg.AppPort
	if port == "" {
		port = "8080"
	}

	log.Printf("✅ Server running on port %s", port)
	log.Printf("✅ Fitur peminjaman buku siap digunakan!")
	r.Run(":" + port)
}
