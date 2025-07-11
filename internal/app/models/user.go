package models

import (
	"time"
)

// User representa um usuário do sistema
type User struct {
	ID           int       `json:"id" db:"id"`
	Email        string    `json:"email" db:"email" validate:"required,email"`
	PasswordHash string    `json:"-" db:"password_hash"` // Campo oculto no JSON
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time `json:"updated_at" db:"updated_at"`
}

// UserSignupRequest representa os dados de registro de usuário
type UserSignupRequest struct {
	Email    string `json:"email" validate:"required,email" example:"usuario@exemplo.com"`
	Password string `json:"password" validate:"required,min=6" example:"senha123"`
}

// UserLoginRequest representa os dados de login
type UserLoginRequest struct {
	Email    string `json:"email" validate:"required,email" example:"usuario@exemplo.com"`
	Password string `json:"password" validate:"required" example:"senha123"`
}

// AuthResponse representa a resposta de autenticação
type AuthResponse struct {
	Token string `json:"token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
	User  User   `json:"user"`
}

// ErrorResponse representa uma resposta de erro
type ErrorResponse struct {
	Error   string `json:"error" example:"Email já está em uso"`
	Message string `json:"message,omitempty" example:"Detalhes adicionais do erro"`
}
