package middleware

import (
    "net/http"
    "strings"
    
    "github.com/gin-gonic/gin"
    "github.com/RajendraMahadana/perpustakaan-clean/internal/utils"
)

type AuthMiddleware struct {
    jwtUtil *utils.JWTUtil
}

func NewAuthMiddleware() *AuthMiddleware {
    return &AuthMiddleware{
        jwtUtil: utils.NewJWTUtil(),
    }
}

// Authenticate middleware untuk autentikasi
func (m *AuthMiddleware) Authenticate() gin.HandlerFunc {
    return func(c *gin.Context) {
        // Ambil token dari header Authorization
        authHeader := c.GetHeader("Authorization")
        if authHeader == "" {
            c.JSON(http.StatusUnauthorized, gin.H{
                "error": "Authorization header is required",
            })
            c.Abort()
            return
        }

        // Format: Bearer {token}
        parts := strings.Split(authHeader, " ")
        if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
            c.JSON(http.StatusUnauthorized, gin.H{
                "error": "Invalid authorization header format",
            })
            c.Abort()
            return
        }

        // Validasi token menggunakan method dari jwtUtil
        claims, err := m.jwtUtil.ValidateToken(parts[1])
        if err != nil {
            c.JSON(http.StatusUnauthorized, gin.H{
                "error": "Invalid or expired token",
            })
            c.Abort()
            return
        }

        // Set user info ke context
        c.Set("user_id", claims.UserID)
        c.Set("email", claims.Email)
        c.Set("role", claims.Role)

        c.Next()
    }
}

// Authorize middleware untuk otorisasi berdasarkan role
func (m *AuthMiddleware) Authorize(allowedRoles ...string) gin.HandlerFunc {
    return func(c *gin.Context) {
        userRole, exists := c.Get("role")
        if !exists {
            c.JSON(http.StatusForbidden, gin.H{
                "error": "Role not found",
            })
            c.Abort()
            return
        }

        // Cek apakah role user termasuk dalam allowedRoles
        for _, role := range allowedRoles {
            if role == userRole {
                c.Next()
                return
            }
        }

        c.JSON(http.StatusForbidden, gin.H{
            "error": "Insufficient permissions",
        })
        c.Abort()
    }
}

// GetUserID helper untuk mengambil user_id dari context
func GetUserID(c *gin.Context) uint {
    userID, exists := c.Get("user_id")
    if !exists {
        return 0
    }
    return userID.(uint)
}

// GetUserRole helper untuk mengambil role dari context
func GetUserRole(c *gin.Context) string {
    role, exists := c.Get("role")
    if !exists {
        return ""
    }
    return role.(string)
}

// GetUserEmail helper untuk mengambil email dari context
func GetUserEmail(c *gin.Context) string {
    email, exists := c.Get("email")
    if !exists {
        return ""
    }
    return email.(string)
}