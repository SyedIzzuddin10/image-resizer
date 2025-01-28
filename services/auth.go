package services

import (
	"database/sql"
	"errors"

	"resizer/models"

	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	db *sql.DB
}

type LoginResponse struct {
	User struct {
		ID       int64  `json:"id"`
		Username string `json:"username"`
	} `json:"user"`
}

func NewAuthService(db *sql.DB) *AuthService {
	return &AuthService{db: db}
}

func (a *AuthService) CreateUser(username, password string) error {
	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	_, err = a.db.Exec("INSERT INTO users (username, password) VALUES (?, ?)", username, hashedPassword)
	return err
}

func (a *AuthService) Login(username, password string) (*LoginResponse, error) {
	var user models.User
	err := a.db.QueryRow("SELECT id, username, password FROM users WHERE username = ?", username).Scan(
		&user.ID,
		&user.Username,
		&user.Password,
	)
	if err == sql.ErrNoRows {
		return nil, errors.New("invalid credentials")
	}
	if err != nil {
		return nil, err
	}

	// Compare password
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	response := &LoginResponse{}
	response.User.ID = user.ID
	response.User.Username = user.Username

	return response, nil
}
