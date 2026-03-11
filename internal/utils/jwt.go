package utils

import (
    "errors"
    "time"
    
    "github.com/golang-jwt/jwt/v5"
    "github.com/RajendraMahadana/perpustakaan-clean/internal/infrastructure/config"
)

type JWTUtil struct {
    secretKey []byte
    expired   int
}

var (
    jwtSecret     []byte
    jwtExpired    int
)

type JWTClaims struct {
    UserID   uint   `json:"user_id"`
    Email    string `json:"email"`
    Role     string `json:"role"`
    jwt.RegisteredClaims
}

// InitJWT inisialisasi konfigurasi JWT (tanpa config)
func InitJWT() {
    // Ini akan dipanggil jika tidak menggunakan config
    // Bisa dikosongkan atau diisi dengan default
}

// InitJWTWithConfig inisialisasi JWT dengan config
func InitJWTWithConfig(cfg *config.Config) {
    jwtSecret = []byte(cfg.JWTSecret)
    jwtExpired = cfg.JWTExpired
}

// NewJWTUtil membuat instance baru JWTUtil
func NewJWTUtil() *JWTUtil {
    return &JWTUtil{
        secretKey: jwtSecret,
        expired:   jwtExpired,
    }
}

// GenerateToken membuat token JWT baru (method)
func (j *JWTUtil) GenerateToken(userID uint, email, role string) (string, error) {
    claims := JWTClaims{
        UserID: userID,
        Email:  email,
        Role:   role,
        RegisteredClaims: jwt.RegisteredClaims{
            ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(j.expired) * time.Hour)),
            IssuedAt:  jwt.NewNumericDate(time.Now()),
            NotBefore: jwt.NewNumericDate(time.Now()),
            Issuer:    "perpustakaan-go",
        },
    }

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString(j.secretKey)
}

// ValidateToken memvalidasi token JWT (method)
func (j *JWTUtil) ValidateToken(tokenString string) (*JWTClaims, error) {
    token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
        if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
            return nil, errors.New("unexpected signing method")
        }
        return j.secretKey, nil
    })

    if err != nil {
        return nil, err
    }

    if claims, ok := token.Claims.(*JWTClaims); ok && token.Valid {
        return claims, nil
    }

    return nil, errors.New("invalid token")
}

// RefreshToken memperbarui token JWT (method)
func (j *JWTUtil) RefreshToken(tokenString string) (string, error) {
    claims, err := j.ValidateToken(tokenString)
    if err != nil {
        return "", err
    }

    return j.GenerateToken(claims.UserID, claims.Email, claims.Role)
}

// ==================== FUNGSI GLOBAL ====================

func GenerateToken(userID uint, email, role string) (string, error) {
    util := NewJWTUtil()
    return util.GenerateToken(userID, email, role)
}

func ValidateToken(tokenString string) (*JWTClaims, error) {
    util := NewJWTUtil()
    return util.ValidateToken(tokenString)
}

func RefreshToken(tokenString string) (string, error) {
    util := NewJWTUtil()
    return util.RefreshToken(tokenString)
}