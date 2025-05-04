package auth

import (
	"database/sql"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"github.com/nktauserum/web-calculation/shared"
	"golang.org/x/crypto/bcrypt"
)

type UserStorage struct {
	db *sql.DB
}

func NewUserStorage(dbPath string) (*UserStorage, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}

	// Создаем таблицу пользователей, если она не существует
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS users (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			username TEXT UNIQUE NOT NULL,
			email TEXT UNIQUE NOT NULL,
			password TEXT NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)
	`)
	if err != nil {
		return nil, err
	}

	return &UserStorage{db: db}, nil
}

func (s *UserStorage) CreateUser(user *shared.RegisterRequest) (*shared.User, error) {
	// Хешируем пароль
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	result, err := s.db.Exec(
		"INSERT INTO users (username, email, password) VALUES (?, ?, ?)",
		user.Username, user.Email, string(hashedPassword),
	)
	if err != nil {
		return nil, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	return &shared.User{
		ID:        id,
		Username:  user.Username,
		Email:     user.Email,
		CreatedAt: time.Now(),
	}, nil
}

func (s *UserStorage) GetUserByUsername(username string) (*shared.User, error) {
	var user shared.User
	var passwordHash string

	err := s.db.QueryRow(
		"SELECT id, username, email, password, created_at FROM users WHERE username = ?",
		username,
	).Scan(&user.ID, &user.Username, &user.Email, &passwordHash, &user.CreatedAt)
	if err != nil {
		return nil, err
	}

	user.Password = passwordHash
	return &user, nil
}

func (s *UserStorage) GetUserByID(id int64) (*shared.User, error) {
	var user shared.User

	err := s.db.QueryRow(
		"SELECT id, username, email, created_at FROM users WHERE id = ?",
		id,
	).Scan(&user.ID, &user.Username, &user.Email, &user.CreatedAt)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (s *UserStorage) Close() error {
	return s.db.Close()
}
