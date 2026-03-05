package main

import (
	"log"
	"os"

	"github.com/RajendraMahadana/perpustakaan-clean/internal/infrastructure/database"
	// "github.com/RajendraMahadana/perpustakaan-clean/internal/usecase"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	bookHandler "github.com/RajendraMahadana/perpustakaan-clean/internal/delivery/http"
	bookRepoImpl "github.com/RajendraMahadana/perpustakaan-clean/internal/infrastructure/repository"
	bookRepo "github.com/RajendraMahadana/perpustakaan-clean/internal/repository"
	bookUsecase "github.com/RajendraMahadana/perpustakaan-clean/internal/usecase"
)

func main() {
	// Load env
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Connect Database
	db := database.NewMySQLConnection()
	database.AutoMigrate(db)

	var repo bookRepo.BookRepository = bookRepoImpl.NewBookRepository(db)
	usecase := bookUsecase.NewBookUsecase(repo)
	handler := bookHandler.NewBookHandler(usecase)

	_ = db //inject repository

	// Init Gin
	r := gin.Default()

	r.POST("/books", handler.Create)
	r.GET("/books", handler.GetAllBooks)
	r.GET("/books/limit", handler.GetBook)
	r.GET("/books/:id", handler.GetByID)
	r.PUT("/books/:id", handler.Update)
	r.DELETE("/books/:id", handler.Delete)

	port := os.Getenv("APP_PORT")
	r.Run(":" + port)
}
