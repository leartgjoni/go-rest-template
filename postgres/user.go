package postgres

import (
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	app "github.com/leartgjoni/go-rest-template"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"strconv"
	"strings"
	"time"

	_ "github.com/lib/pq"
)

// Ensure service implements interface.
var _ app.UserService = &UserService{}

// UserService represents a service to manage users.
type UserService struct {
	db        *DB
	apiSecret string
}

// NewUserService returns a new instance of UserService.
func NewUserService(db *DB, apiSecret string) *UserService {
	return &UserService{
		db:        db,
		apiSecret: apiSecret,
	}
}

func (s *UserService) CreateToken(userId uint32) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userId": userId,
		"exp":    time.Now().Add(time.Hour * 24).Unix(),
	})

	return token.SignedString([]byte(s.apiSecret))
}

func (s *UserService) ExtractAuthenticationToken(r *http.Request) (uint32, error) {
	tokenString := extractToken(r)
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(s.apiSecret), nil
	})
	if err != nil {
		return 0, err
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if ok && token.Valid {
		uid, err := strconv.ParseUint(fmt.Sprintf("%.0f", claims["userId"]), 10, 32)
		if err != nil {
			return 0, err
		}
		return uint32(uid), nil
	}
	return 0, nil
}

func (s *UserService) Save(user *app.User) error {
	count := 0
	err := s.db.QueryRow("SELECT COUNT(id) FROM users WHERE email = $1", user.Email).Scan(&count)
	if err != nil {
		return err
	}

	if count > 0 {
		return app.ErrEmailAlreadyUsed
	}

	hashedPassword, err := hash(user.Password)
	if err != nil {
		return app.ErrWrongPasswordFormat
	}

	user.Password = string(hashedPassword)

	row := s.db.QueryRow("INSERT INTO users (username, email, password, created_at, updated_at) VALUES ($1, $2, $3, $4, $5) RETURNING id", user.Username, user.Email, user.Password, user.CreatedAt, user.UpdatedAt)

	if err := row.Scan(&user.ID); err != nil {
		return nil
	}

	if user.ID == 0 {
		return errors.New("unable to save")
	}

	return nil
}

func (s *UserService) GetById(userId uint32) (*app.User, error) {
	var user app.User
	err := s.db.QueryRow("SELECT * FROM users WHERE id = $1", userId).Scan(&user.ID, &user.Username, &user.Email, &user.Password, &user.CreatedAt, &user.UpdatedAt)

	if err != nil || user.ID == 0 {
		return &app.User{}, app.ErrUserNotFound
	}

	return &user, nil
}

func (s *UserService) Login(u *app.User) (string, error) {
	var row struct {
		id        uint32
		username  string
		password  string
		createdAt time.Time
		updatedAt time.Time
	}

	err := s.db.QueryRow("SELECT id, username, password, created_at, updated_at FROM users WHERE email = $1 LIMIT 1", u.Email).Scan(&row.id, &row.username, &row.password, &row.createdAt, &row.updatedAt)

	if err != nil || row.id == 0 {
		return "", app.ErrWrongCredentials
	}
	err = verifyPassword(row.password, u.Password)
	if err != nil {
		return "", app.ErrWrongCredentials
	}

	u.ID = row.id
	u.Username = row.username
	u.Password = row.password
	u.CreatedAt = row.createdAt
	u.UpdatedAt = row.updatedAt

	return s.CreateToken(u.ID)
}

func extractToken(r *http.Request) string {
	bearerToken := r.Header.Get("Authorization")
	if len(strings.Split(bearerToken, " ")) == 2 {
		return strings.Split(bearerToken, " ")[1]
	}
	return ""
}

func hash(password string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
}

func verifyPassword(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}
