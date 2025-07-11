package services

import (
	"database/sql"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"

	"github.com/luizdequeiroz/rpg-backend/internal/app/models"
	"github.com/luizdequeiroz/rpg-backend/pkg/db"
)

var (
	ErrEmailAlreadyExists = errors.New("email já está em uso")
	ErrInvalidCredentials = errors.New("credenciais inválidas")
	ErrUserNotFound       = errors.New("usuário não encontrado")
)

// AuthService gerencia operações de autenticação
type AuthService struct {
	db        *db.DB
	jwtSecret []byte
}

// NewAuthService cria uma nova instância do serviço de autenticação
func NewAuthService(database *db.DB) *AuthService {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = "default-secret-key-change-in-production"
	}

	return &AuthService{
		db:        database,
		jwtSecret: []byte(secret),
	}
}

// Signup registra um novo usuário
func (s *AuthService) Signup(req models.UserSignupRequest) (*models.AuthResponse, error) {
	// Verificar se email já existe
	var count int
	err := s.db.Get(&count, "SELECT COUNT(*) FROM users WHERE email = ?", req.Email)
	if err != nil {
		return nil, fmt.Errorf("erro ao verificar email: %w", err)
	}
	if count > 0 {
		return nil, ErrEmailAlreadyExists
	}

	// Hash da senha
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("erro ao gerar hash da senha: %w", err)
	}

	// Inserir usuário no banco
	query := `
		INSERT INTO users (email, password_hash, created_at, updated_at) 
		VALUES (?, ?, ?, ?)
	`
	now := time.Now()
	result, err := s.db.Exec(query, req.Email, string(hashedPassword), now, now)
	if err != nil {
		return nil, fmt.Errorf("erro ao criar usuário: %w", err)
	}

	// Obter ID do usuário criado
	userID, err := result.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("erro ao obter ID do usuário: %w", err)
	}

	// Buscar usuário completo
	user := models.User{}
	err = s.db.Get(&user, "SELECT id, email, created_at, updated_at FROM users WHERE id = ?", userID)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar usuário criado: %w", err)
	}

	// Gerar JWT
	token, err := s.generateJWT(user.ID, user.Email)
	if err != nil {
		return nil, fmt.Errorf("erro ao gerar token: %w", err)
	}

	return &models.AuthResponse{
		Token: token,
		User:  user,
	}, nil
}

// Login autentica um usuário
func (s *AuthService) Login(req models.UserLoginRequest) (*models.AuthResponse, error) {
	// Buscar usuário por email
	user := models.User{}
	query := "SELECT id, email, password_hash, created_at, updated_at FROM users WHERE email = ?"
	err := s.db.Get(&user, query, req.Email)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrInvalidCredentials
		}
		return nil, fmt.Errorf("erro ao buscar usuário: %w", err)
	}

	// Verificar senha
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password))
	if err != nil {
		return nil, ErrInvalidCredentials
	}

	// Gerar JWT
	token, err := s.generateJWT(user.ID, user.Email)
	if err != nil {
		return nil, fmt.Errorf("erro ao gerar token: %w", err)
	}

	// Limpar senha do retorno
	user.PasswordHash = ""

	return &models.AuthResponse{
		Token: token,
		User:  user,
	}, nil
}

// generateJWT gera um token JWT para o usuário
func (s *AuthService) generateJWT(userID int, email string) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"email":   email,
		"exp":     time.Now().Add(time.Hour * 24).Unix(), // Expira em 24 horas
		"iat":     time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(s.jwtSecret)
}

// ValidateJWT valida um token JWT e retorna os claims
func (s *AuthService) ValidateJWT(tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("método de assinatura inesperado: %v", token.Header["alg"])
		}
		return s.jwtSecret, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("token inválido")
}

// GetUserByID busca um usuário por ID
func (s *AuthService) GetUserByID(userID int) (*models.User, error) {
	user := models.User{}
	err := s.db.Get(&user, "SELECT id, email, created_at, updated_at FROM users WHERE id = ?", userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrUserNotFound
		}
		return nil, fmt.Errorf("erro ao buscar usuário: %w", err)
	}
	return &user, nil
}
