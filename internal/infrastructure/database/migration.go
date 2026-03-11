package database

import (
	"log"
	"time"
	
	"gorm.io/gorm"
	"golang.org/x/crypto/bcrypt"
)

// Book Model
type BookModel struct {
	ID        uint   `gorm:"primaryKey"`
	Title     string `gorm:"size:255;not null;index"`
	Author    string `gorm:"size:255;not null;index"`
	ISBN      string `gorm:"size:100;uniqueIndex"`
	Stock     int    `gorm:"default:0"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

// TableName untuk BookModel
func (BookModel) TableName() string {
	return "books"
}

// User Model
type UserModel struct {
	ID        uint      `gorm:"primaryKey"`
	Name      string    `gorm:"size:100;not null"`
	Email     string    `gorm:"size:100;uniqueIndex;not null"`
	Password  string    `gorm:"size:255;not null"`
	Role      string    `gorm:"size:20;default:user;index"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

// TableName untuk UserModel
func (UserModel) TableName() string {
	return "users"
}

// BeforeCreate hook untuk hash password sebelum insert
func (u *UserModel) BeforeCreate(tx *gorm.DB) error {
	if u.Password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
		if err != nil {
			return err
		}
		u.Password = string(hashedPassword)
	}
	return nil
}

// BeforeUpdate hook untuk hash password sebelum update (jika password diubah)
func (u *UserModel) BeforeUpdate(tx *gorm.DB) error {
	if u.Password != "" {
		// Cek apakah password sudah di-hash atau belum
		if len(u.Password) < 60 { // bcrypt hash biasanya 60 karakter
			hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
			if err != nil {
				return err
			}
			u.Password = string(hashedPassword)
		}
	}
	return nil
}

// Response DTO untuk User (tanpa password)
type UserResponse struct {
	ID        uint      `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Role      string    `json:"role"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// ToResponse mengkonversi UserModel ke UserResponse
func (u *UserModel) ToResponse() *UserResponse {
	return &UserResponse{
		ID:        u.ID,
		Name:      u.Name,
		Email:     u.Email,
		Role:      u.Role,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}
}

// CheckPassword memverifikasi password
func (u *UserModel) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	return err == nil
}

// AutoMigrate menjalankan migrasi database
func AutoMigrate(db *gorm.DB) {
	// Daftar model yang akan dimigrasi
	err := db.AutoMigrate(
		&BookModel{},
		&UserModel{},
	)
	
	if err != nil {
		log.Fatal("❌ Failed to migrate database:", err)
	}

	// Buat admin default jika tabel users kosong
	createDefaultAdmin(db)

	log.Println("✅ Database migrated successfully")
}

// createDefaultAdmin membuat admin default jika belum ada
func createDefaultAdmin(db *gorm.DB) {
	var count int64
	db.Model(&UserModel{}).Count(&count)

	if count == 0 {
		// Buat admin default
		admin := &UserModel{
			Name:     "Administrator",
			Email:    "admin@perpustakaan.com",
			Password: "admin123", // Akan di-hash otomatis oleh BeforeCreate hook
			Role:     "admin",
		}

		if err := db.Create(admin).Error; err != nil {
			log.Println("⚠️  Warning: Failed to create default admin:", err)
		} else {
			log.Println("✅ Default admin created (email: admin@perpustakaan.com, password: admin123)")
		}
	}
}

// Seeder function (optional - untuk data dummy)
func SeedDatabase(db *gorm.DB) {
	// Seed books jika diperlukan
	seedBooks(db)
}

func seedBooks(db *gorm.DB) {
	var count int64
	db.Model(&BookModel{}).Count(&count)

	if count == 0 {
		books := []BookModel{
			{
				Title:  "Clean Code",
				Author: "Robert C. Martin",
				ISBN:   "978-0132350884",
				Stock:  5,
			},
			{
				Title:  "The Pragmatic Programmer",
				Author: "David Thomas",
				ISBN:   "978-0201616224",
				Stock:  3,
			},
			{
				Title:  "Design Patterns",
				Author: "Erich Gamma",
				ISBN:   "978-0201633610",
				Stock:  2,
			},
		}

		for _, book := range books {
			db.Create(&book)
		}
		log.Println("✅ Sample books seeded")
	}
}