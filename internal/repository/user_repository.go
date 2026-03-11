package repository

import (
    "github.com/RajendraMahadana/perpustakaan-clean/internal/domain"
)

type UserRepository interface {
    Create(user *domain.User) error
    FindByEmail(email string) (*domain.User, error)
    FindByID(id uint) (*domain.User, error)
    Update(user *domain.User) error
    Delete(id uint) error
}

type userRepository struct {
    // Tambahkan dependency database di sini
}

func NewUserRepository() UserRepository {
    return &userRepository{}
}

func (r *userRepository) Create(user *domain.User) error {
    // Implementasi create user
    return nil
}

func (r *userRepository) FindByEmail(email string) (*domain.User, error) {
    // Implementasi find by email
    return nil, nil
}

func (r *userRepository) FindByID(id uint) (*domain.User, error) {
    // Implementasi find by ID
    return nil, nil
}

func (r *userRepository) Update(user *domain.User) error {
    // Implementasi update user
    return nil
}

func (r *userRepository) Delete(id uint) error {
    // Implementasi delete user
    return nil
}