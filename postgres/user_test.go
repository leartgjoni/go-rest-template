package postgres

import (
	"errors"
	"fmt"
	"github.com/DATA-DOG/go-sqlmock"
	app "github.com/leartgjoni/go-rest-template"
	"net/http"
	"testing"
	"time"
)

func TestUserService_ExtractAuthenticationToken(t *testing.T) {
	db, _, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	tests := []struct {
		name    string
		userId  uint32
		invalid string
		error   error
	}{
		{
			name:    "extracts correctly",
			userId:  1,
			invalid: "",
			error:   nil,
		},
		{
			name:    "wrong token encoding",
			userId:  0,
			invalid: "wrong format",
			error:   errors.New("token contains an invalid number of segments"),
		},
		{
			name:    "wrong header format",
			userId:  0,
			invalid: "random",
			error:   errors.New("token contains an invalid number of segments"),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			us := NewUserService(&DB{db}, "random")

			token, err := us.CreateToken(test.userId)
			if err != nil {
				t.Fatal("cannot create token", token)
			}

			r, _ := http.NewRequest("", "", nil)
			if test.invalid != "" {
				r.Header.Set("Authorization", test.invalid)
			} else {
				r.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
			}

			userId, err := us.ExtractAuthenticationToken(r)

			if test.invalid != "" && err.Error() != test.error.Error() {
				t.Fatalf("wrong error. expected %s but got %s", test.error, err)
			}

			if userId != test.userId {
				t.Fatalf("wrong user id. Expected %v but got %v", test.userId, userId)
			}
		})
	}
}

func TestUserService_Save(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	tests := []struct {
		name         string
		countResult  *sqlmock.Rows
		insertResult *sqlmock.Rows
		expected     error
	}{
		{
			name: "Email already used",
			countResult: sqlmock.NewRows([]string{"count"}).
				AddRow(1),
			insertResult: nil,
			expected:     app.ErrEmailAlreadyUsed,
		},
		{
			name: "User without ID after saving",
			countResult: sqlmock.NewRows([]string{"count"}).
				AddRow(0),
			insertResult: sqlmock.NewRows([]string{"id"}).
				AddRow(0),
			expected: errors.New("unable to save"),
		},
		{
			name: "Success",
			countResult: sqlmock.NewRows([]string{"count"}).
				AddRow(0),
			insertResult: sqlmock.NewRows([]string{"id"}).
				AddRow(1),
			expected: nil,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			mock.ExpectQuery("^SELECT (.+) FROM users*").WillReturnRows(test.countResult)
			if test.insertResult != nil {
				mock.ExpectQuery("^INSERT INTO users *").WillReturnRows(test.insertResult)
			}

			us := NewUserService(&DB{db}, "random")

			err = us.Save(&app.User{})

			// we make sure that all expectations were met
			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}

			if err != nil {
				if err.Error() != test.expected.Error() {
					t.Fatalf("wrong error. expected %s but got %s", test.expected, err)
				}
			} else {
				if err != test.expected {
					t.Fatalf("wrong error. expected %s but got %s", test.expected, err)
				}
			}
		})
	}
}

func TestUserService_GetById(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	dbUser := app.User{
		ID:        1,
		Username:  "test",
		Email:     "test@test.com",
		Password:  "password-hashed",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	tests := []struct {
		name      string
		sqlResult *sqlmock.Rows
		error     error
		user      app.User
	}{
		{
			name:      "Found by id",
			sqlResult: sqlmock.NewRows([]string{"id", "username", "email", "password", "created_at", "updated_at"}).AddRow(dbUser.ID, dbUser.Username, dbUser.Email, dbUser.Password, dbUser.CreatedAt, dbUser.UpdatedAt),
			error:     nil,
			user:      dbUser,
		},
		{
			name:      "Not found by id",
			sqlResult: sqlmock.NewRows([]string{"id", "username", "email", "password", "created_at", "updated_at"}),
			error:     app.ErrUserNotFound,
			user:      app.User{},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			mock.ExpectQuery("^SELECT (.+) FROM users WHERE id*").WillReturnRows(test.sqlResult)

			us := NewUserService(&DB{db}, "random")

			user, err := us.GetById(1)

			// we make sure that all expectations were met
			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}

			if test.error != nil && err != test.error {
				t.Fatalf("wrong error. expected %s but got %s", test.error, err)
			}

			if err == nil {
				if user.ID != test.user.ID || user.Username != test.user.Username || user.Email != test.user.Email || user.Password != test.user.Password || !user.CreatedAt.Equal(test.user.CreatedAt) || !user.UpdatedAt.Equal(test.user.UpdatedAt) {
					t.Fatalf("wrong user. expected %v but got %v", test.user, user)
				}
			}
		})
	}
}

func TestUserService_Login(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	hashedPassword, err := hash("password")
	if err != nil {
		t.Fatal("error while hashing password")
	}

	dbUser := app.User{
		ID:        1,
		Username:  "test",
		Email:     "test@test.com",
		Password:  string(hashedPassword),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	tests := []struct {
		name      string
		sqlResult *sqlmock.Rows
		error     error
	}{
		{
			name:      "correct login",
			sqlResult: sqlmock.NewRows([]string{"id", "username", "password", "created_at", "updated_at"}).AddRow(dbUser.ID, dbUser.Username, dbUser.Password, dbUser.CreatedAt, dbUser.UpdatedAt),
			error:     nil,
		},
		{
			name:      "wrong email",
			sqlResult: sqlmock.NewRows([]string{"id", "username", "password", "created_at", "updated_at"}),
			error:     app.ErrWrongCredentials,
		},
		{
			name:      "wrong password",
			sqlResult: sqlmock.NewRows([]string{"id", "username", "password", "created_at", "updated_at"}).AddRow(dbUser.ID, dbUser.Username, "password", dbUser.CreatedAt, dbUser.UpdatedAt),
			error:     app.ErrWrongCredentials,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			mock.ExpectQuery("^SELECT (.+) FROM users WHERE email*").WillReturnRows(test.sqlResult)

			us := NewUserService(&DB{db}, "random")

			token, err := us.Login(&app.User{Email: "test@test.com", Password: "password"})

			// we make sure that all expectations were met
			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}

			if err != test.error {
				t.Fatalf("wrong error. expected %s but got %s", test.error, err)
			}

			if err == nil && token == "" {
				t.Fatal("token was not expected to be empty")
			}
		})
	}
}
