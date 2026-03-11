package usecase

import (
    "errors"
    "time"
    
    "golang.org/x/crypto/bcrypt"
    "github.com/RajendraMahadana/perpustakaan-clean/internal/domain"
    "github.com/RajendraMahadana/perpustakaan-clean/internal/dto"
    "github.com/RajendraMahadana/perpustakaan-clean/internal/repository"
    "github.com/RajendraMahadana/perpustakaan-clean/internal/utils"  
)

type AuthUsecase interface {
    Register(req dto.RegisterRequest) (*dto.AuthResponse, error)
    Login(req dto.LoginRequest) (*dto.AuthResponse, error)
    RefreshToken(userID uint) (*dto.AuthResponse, error)
    GetUserByID(userID uint) (*domain.User, error)
}

type authUsecase struct {
    userRepo repository.UserRepository
}

func NewAuthUsecase(userRepo repository.UserRepository) AuthUsecase {
    return &authUsecase{
        userRepo: userRepo,
    }
}

func (u *authUsecase) Register(req dto.RegisterRequest) (*dto.AuthResponse, error) {
    // Cek apakah email sudah terdaftar
    existingUser, _ := u.userRepo.FindByEmail(req.Email)
    if existingUser != nil {
        return nil, errors.New("email already registered")
    }

    // Hash password
    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
    if err != nil {
        return nil, err
    }

    // Set default role jika tidak disertakan
    if req.Role == "" {
        req.Role = "user"
    }

    // Buat user baru
    user := &domain.User{
        Name:     req.Name,
        Email:    req.Email,
        Password: string(hashedPassword),
        Role:     req.Role,
    }

    if err := u.userRepo.Create(user); err != nil {
        return nil, err
    }

    // Generate JWT token menggunakan fungsi global
    token, err := utils.GenerateToken(user.ID, user.Email, user.Role)
    if err != nil {
        return nil, err
    }

    // Dapatkan expiry dari konfigurasi JWT
    expiry := time.Now().Add(time.Duration(24) * time.Hour).Unix() // default 24 jam

    response := &dto.AuthResponse{
        Token:     token,
        UserID:    user.ID,
        Name:      user.Name,
        Email:     user.Email,
        Role:      user.Role,
        ExpiresIn: expiry,
    }

    return response, nil
}

func (u *authUsecase) Login(req dto.LoginRequest) (*dto.AuthResponse, error) {
    // Cari user berdasarkan email
    user, err := u.userRepo.FindByEmail(req.Email)
    if err != nil {
        return nil, errors.New("invalid email or password")
    }

    // Verifikasi password
    err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
    if err != nil {
        return nil, errors.New("invalid email or password")
    }

    // Generate JWT token menggunakan fungsi global
    token, err := utils.GenerateToken(user.ID, user.Email, user.Role)
    if err != nil {
        return nil, err
    }

    // Dapatkan expiry dari konfigurasi JWT
    expiry := time.Now().Add(time.Duration(24) * time.Hour).Unix() // default 24 jam

    response := &dto.AuthResponse{
        Token:     token,
        UserID:    user.ID,
        Name:      user.Name,
        Email:     user.Email,
        Role:      user.Role,
        ExpiresIn: expiry,
    }

    return response, nil
}

func (u *authUsecase) RefreshToken(userID uint) (*dto.AuthResponse, error) {
    // Cari user berdasarkan ID
    user, err := u.userRepo.FindByID(userID)
    if err != nil {
        return nil, errors.New("user not found")
    }

    // Generate new token menggunakan fungsi global
    token, err := utils.GenerateToken(user.ID, user.Email, user.Role)
    if err != nil {
        return nil, err
    }

    // Dapatkan expiry dari konfigurasi JWT
    expiry := time.Now().Add(time.Duration(24) * time.Hour).Unix() // default 24 jam

    response := &dto.AuthResponse{
        Token:     token,
        UserID:    user.ID,
        Name:      user.Name,
        Email:     user.Email,
        Role:      user.Role,
        ExpiresIn: expiry,
    }

    return response, nil
}

func (u *authUsecase) GetUserByID(userID uint) (*domain.User, error) {
    user, err := u.userRepo.FindByID(userID)
    if err != nil {
        return nil, errors.New("user not found")
    }
    
    return user, nil
}