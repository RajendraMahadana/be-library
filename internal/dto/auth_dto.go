package dto

// LoginRequest DTO untuk request login
type LoginRequest struct {
    Email    string `json:"email" binding:"required,email"`
    Password string `json:"password" binding:"required,min=6"`
}

// RegisterRequest DTO untuk request register
type RegisterRequest struct {
    Name     string `json:"name" binding:"required,min=3"`
    Email    string `json:"email" binding:"required,email"`
    Password string `json:"password" binding:"required,min=6"`
    Role     string `json:"role" binding:"omitempty,oneof=admin user"`
}

// AuthResponse DTO untuk response autentikasi
type AuthResponse struct {
    Token     string `json:"token"`
    UserID    uint   `json:"user_id"`
    Name      string `json:"name"`
    Email     string `json:"email"`
    Role      string `json:"role"`
    ExpiresIn int64  `json:"expires_in"` // dalam detik
}