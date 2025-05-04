package auth

import (
	"errors"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/nktauserum/web-calculation/shared"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	userStorage *UserStorage
	jwtSecret   []byte
	tokenExpiry time.Duration
}

func NewAuthService(userStorage *UserStorage, jwtSecret string, tokenExpiry time.Duration) *AuthService {
	return &AuthService{
		userStorage: userStorage,
		jwtSecret:   []byte(jwtSecret),
		tokenExpiry: tokenExpiry,
	}
}

func (s *AuthService) Register(req *shared.RegisterRequest) (*shared.AuthResponse, error) {
	user, err := s.userStorage.CreateUser(req)
	if err != nil {
		return nil, err
	}

	token, err := s.generateToken(user.ID)
	if err != nil {
		return nil, err
	}

	return &shared.AuthResponse{
		Token: token,
		User:  *user,
	}, nil
}

func (s *AuthService) Login(req *shared.LoginRequest) (*shared.AuthResponse, error) {
	user, err := s.userStorage.GetUserByUsername(req.Username)
	if err != nil {
		return nil, errors.New("неверное имя пользователя или пароль")
	}

	// Проверяем пароль
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		return nil, errors.New("неверное имя пользователя или пароль")
	}

	token, err := s.generateToken(user.ID)
	if err != nil {
		return nil, err
	}

	// Очищаем пароль перед отправкой
	user.Password = ""

	return &shared.AuthResponse{
		Token: token,
		User:  *user,
	}, nil
}

func (s *AuthService) ValidateToken(tokenString string) (int64, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("неожиданный метод подписи")
		}
		return s.jwtSecret, nil
	})

	if err != nil {
		return 0, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		userID := int64(claims["user_id"].(float64))
		return userID, nil
	}

	return 0, errors.New("недействительный токен")
}

func (s *AuthService) generateToken(userID int64) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(s.tokenExpiry).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(s.jwtSecret)
}
