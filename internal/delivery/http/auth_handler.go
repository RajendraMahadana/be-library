package http

import (
    "net/http"
    
    "github.com/gin-gonic/gin"
    "github.com/RajendraMahadana/perpustakaan-clean/internal/dto"
    "github.com/RajendraMahadana/perpustakaan-clean/internal/usecase"
    "github.com/RajendraMahadana/perpustakaan-clean/pkg/middleware"
)

type AuthHandler struct {
    authUsecase usecase.AuthUsecase  // Sekarang sudah terdefinisi
}

func NewAuthHandler(authUsecase usecase.AuthUsecase) *AuthHandler {
    return &AuthHandler{
        authUsecase: authUsecase,
    }
}

// Register handler
func (h *AuthHandler) Register(c *gin.Context) {
    var req dto.RegisterRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{
            "error": err.Error(),
        })
        return
    }

    response, err := h.authUsecase.Register(req)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{
            "error": err.Error(),
        })
        return
    }

    c.JSON(http.StatusCreated, response)
}

// Login handler
func (h *AuthHandler) Login(c *gin.Context) {
    var req dto.LoginRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{
            "error": err.Error(),
        })
        return
    }

    response, err := h.authUsecase.Login(req)
    if err != nil {
        c.JSON(http.StatusUnauthorized, gin.H{
            "error": "Invalid email or password",
        })
        return
    }

    c.JSON(http.StatusOK, response)
}

// RefreshToken handler
func (h *AuthHandler) RefreshToken(c *gin.Context) {
    userID := middleware.GetUserID(c)
    
    response, err := h.authUsecase.RefreshToken(userID)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{
            "error": err.Error(),
        })
        return
    }

    c.JSON(http.StatusOK, response)
}

// GetProfile handler
func (h *AuthHandler) GetProfile(c *gin.Context) {
    userID := middleware.GetUserID(c)
    
    user, err := h.authUsecase.GetUserByID(userID)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{
            "error": err.Error(),
        })
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "user_id": user.ID,
        "name":    user.Name,
        "email":   user.Email,
        "role":    user.Role,
    })
}

// Logout handler
func (h *AuthHandler) Logout(c *gin.Context) {
    // Implementasi logout (bisa menggunakan blacklist token jika diperlukan)
    c.JSON(http.StatusOK, gin.H{
        "message": "Successfully logged out",
    })
}