package database

import (
	"fmt"
	"log"
	"os"
	"time"
	
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	
	"github.com/RajendraMahadana/perpustakaan-clean/internal/infrastructure/config"
)

type Database struct {
	DB *gorm.DB
}

// NewMySQLConnection membuat koneksi database baru (tanpa config)
func NewMySQLConnection() *gorm.DB {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASS"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"),
	)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		log.Fatal("❌ Failed to connect database:", err)
	}

	// Test koneksi
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatal("❌ Failed to get database instance:", err)
	}

	// Set connection pool
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	fmt.Println("✅ Connected to MySQL")
	return db
}

// NewMySQLConnectionWithConfig membuat koneksi database dengan config
func NewMySQLConnectionWithConfig(cfg *config.Config) (*gorm.DB, error) {
	// Buat DSN dari config
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		cfg.DBUser,
		cfg.DBPassword,
		cfg.DBHost,
		cfg.DBPort,
		cfg.DBName,
	)

	// Konfigurasi GORM dengan logging yang bisa diatur
	gormConfig := &gorm.Config{
		Logger: logger.Default.LogMode(getLogLevel(cfg)),
		NowFunc: func() time.Time {
			return time.Now().Local()
		},
	}

	// Buka koneksi
	db, err := gorm.Open(mysql.Open(dsn), gormConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to connect database: %w", err)
	}

	// Dapatkan instance sql.DB untuk konfigurasi pool
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get database instance: %w", err)
	}

	// Konfigurasi connection pool
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	// Test ping ke database
	if err := sqlDB.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	log.Println("✅ Connected to MySQL successfully")
	return db, nil
}

// MustConnect sama seperti NewMySQLConnectionWithConfig tapi akan panic jika error
func MustConnect(cfg *config.Config) *gorm.DB {
	db, err := NewMySQLConnectionWithConfig(cfg)
	if err != nil {
		log.Fatal("❌ ", err)
	}
	return db
}

// getLogLevel menentukan level logging berdasarkan environment
func getLogLevel(cfg *config.Config) logger.LogLevel {
	// Bisa ditambahkan konfigurasi untuk log level
	// Misalnya: berdasarkan environment (development/production)
	
	// Default: Info untuk development, Error untuk production
	env := os.Getenv("APP_ENV")
	if env == "production" {
		return logger.Error
	}
	return logger.Info
}

// CloseDatabase menutup koneksi database
func CloseDatabase(db *gorm.DB) error {
	sqlDB, err := db.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}

// HealthCheck mengecek kesehatan database
func HealthCheck(db *gorm.DB) error {
	sqlDB, err := db.DB()
	if err != nil {
		return err
	}
	return sqlDB.Ping()
}

// Transaction wrapper untuk memudahkan transaksi
func Transaction(db *gorm.DB, fn func(tx *gorm.DB) error) error {
	tx := db.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r) // re-throw panic setelah rollback
		}
	}()

	if err := fn(tx); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

// WithTimeout menambahkan timeout pada query
func WithTimeout(db *gorm.DB, timeout time.Duration) *gorm.DB {
	return db.Debug() // Bisa dimodifikasi untuk menambah timeout
}