package http_test

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	. "github.com/leartgjoni/go-rest-template/http"
	"strings"

	app "github.com/leartgjoni/go-rest-template"
	"github.com/leartgjoni/go-rest-template/mock"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestAuthHandler_HandleSignup(t *testing.T) {
	// mock time
	now := time.Unix(0, 0)
	nowString := now.Format(time.RFC3339)

	var tests = []struct {
		name               string
		SaveFn             func(user *app.User) error
		SaveInvoked        bool
		CreateTokenFn      func(userId uint32) (string, error)
		CreateTokenInvoked bool
		body               []byte
		expectedResponse   string
	}{
		{
			name: "success",
			SaveFn: func(user *app.User) error {
				user.ID = 1
				user.CreatedAt = now
				user.UpdatedAt = now
				return nil
			},
			SaveInvoked: true,
			CreateTokenFn: func(userId uint32) (string, error) {
				return "random-token", nil
			},
			CreateTokenInvoked: true,
			body:               []byte(`{"username":"test","email":"test@test.com","password":"random"}`),
			expectedResponse:   fmt.Sprintf(`{"id":1,"username":"test","email":"test@test.com","created_at":"%s","updated_at":"%s","token":"random-token"}`, nowString, nowString),
		},
		{
			name: "Save() error",
			SaveFn: func(user *app.User) error {
				return errors.New("save fn error")
			},
			SaveInvoked:        true,
			CreateTokenFn:      nil,
			CreateTokenInvoked: false,
			body:               []byte(`{"username":"test","email":"test@test.com","password":"random"}`),
			expectedResponse:   `{"message":"Server Error","error":"save fn error"}`,
		},
		{
			name: "CreateToken() error",
			SaveFn: func(user *app.User) error {
				return nil
			},
			SaveInvoked: true,
			CreateTokenFn: func(userId uint32) (string, error) {
				return "", errors.New("create token fn error")
			},
			CreateTokenInvoked: true,
			body:               []byte(`{"username":"test","email":"test@test.com","password":"random"}`),
			expectedResponse:   `{"message":"Server Error","error":"create token fn error"}`,
		},
		{
			name:               "Invalid request",
			SaveFn:             nil,
			SaveInvoked:        false,
			CreateTokenFn:      nil,
			CreateTokenInvoked: false,
			body:               nil,
			expectedResponse:   `{"message":"Invalid request.","error":"EOF"}`,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// Inject our mock into our handler.
			var us mock.UserService
			h := NewAuthHandler(&us)

			// Mock our Save() call.
			us.SaveFn = test.SaveFn

			// Mock our CreateToken() call.
			us.CreateTokenFn = test.CreateTokenFn

			// request body
			var body = test.body

			// Invoke the handler.
			w := httptest.NewRecorder()
			r, _ := http.NewRequest("POST", "/signup", bytes.NewBuffer(body))
			r.Header.Set("Content-Type", "application/json")
			httpHandler := http.HandlerFunc(h.HandleSignup)
			httpHandler.ServeHTTP(w, r)

			// Validate mock.
			if us.SaveInvoked != test.SaveInvoked {
				t.Fatalf("expected SaveInvoked to be %v", test.SaveInvoked)
			}

			if us.CreateTokenInvoked != test.CreateTokenInvoked {
				t.Fatalf("expected SaveInvoked to be %v", test.CreateTokenInvoked)
			}

			expected := test.expectedResponse
			received := strings.TrimSpace(w.Body.String())

			if received != expected {
				t.Fatalf("expected %s but received %s", expected, received)
			}
		})
	}
}

func TestAuthHandler_HandleLogin(t *testing.T) {
	// mock time
	now := time.Unix(0, 0)
	nowString := now.Format(time.RFC3339)

	var tests = []struct {
		name             string
		LoginFn          func(u *app.User) (string, error)
		LoginInvoked     bool
		body             []byte
		expectedResponse string
	}{
		{
			name: "success",
			LoginFn: func(u *app.User) (string, error) {
				u.ID = 1
				u.Username = "test"
				u.Password = "hashed-password"
				u.CreatedAt = now
				u.UpdatedAt = now
				return "random-token", nil
			},
			LoginInvoked:     true,
			body:             []byte(`{"email":"test@test.com","password":"random"}`),
			expectedResponse: fmt.Sprintf(`{"id":1,"username":"test","email":"test@test.com","created_at":"%s","updated_at":"%s","token":"random-token"}`, nowString, nowString),
		},
		{
			name: "wrong credentials",
			LoginFn: func(u *app.User) (string, error) {
				return "", app.ErrWrongCredentials
			},
			LoginInvoked:     true,
			body:             []byte(`{"email":"test@test.com","password":"random"}`),
			expectedResponse: `{"message":"Unauthorized"}`,
		},
		{
			name: "invalid request",
			LoginFn: func(u *app.User) (string, error) {
				return "", app.ErrWrongCredentials
			},
			LoginInvoked:     false,
			body:             nil,
			expectedResponse: `{"message":"Invalid request.","error":"EOF"}`,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// Inject our mock into our handler.
			var us mock.UserService
			h := NewAuthHandler(&us)

			// Mock our Login() call.
			us.LoginFn = test.LoginFn

			// request body
			var body = test.body

			// Invoke the handler.
			w := httptest.NewRecorder()
			r, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(body))
			r.Header.Set("Content-Type", "application/json")
			httpHandler := http.HandlerFunc(h.HandleLogin)
			httpHandler.ServeHTTP(w, r)

			// Validate mock.
			if us.LoginInvoked != test.LoginInvoked {
				t.Fatalf("expected LoginInvoked to be %v", test.LoginInvoked)
			}

			expected := test.expectedResponse
			received := strings.TrimSpace(w.Body.String())

			if received != expected {
				t.Fatalf("expected %s but received %s", expected, received)
			}
		})
	}
}

func TestAuthHandler_HandleMe(t *testing.T) {
	// mock time
	now := time.Unix(0, 0)
	nowString := now.Format(time.RFC3339)

	var tests = []struct {
		name             string
		GetByIdFn        func(userId uint32) (*app.User, error)
		GetByIdInvoked   bool
		expectedResponse string
	}{
		{
			name: "success",
			GetByIdFn: func(userId uint32) (*app.User, error) {
				return &app.User{
					ID:        1,
					Username:  "test",
					Email:     "test@test.com",
					Password:  "hashed-password",
					CreatedAt: now,
					UpdatedAt: now,
				}, nil
			},
			GetByIdInvoked:   true,
			expectedResponse: fmt.Sprintf(`{"id":1,"username":"test","email":"test@test.com","created_at":"%s","updated_at":"%s"}`, nowString, nowString),
		},
		{
			name: "wrong token",
			GetByIdFn: func(userId uint32) (*app.User, error) {
				return &app.User{}, app.ErrUserNotFound
			},
			GetByIdInvoked:   true,
			expectedResponse: `{"message":"Resource not found."}`,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// Inject our mock into our handler.
			var us mock.UserService
			h := NewAuthHandler(&us)

			// Mock our Login() call.
			us.GetByIdFn = test.GetByIdFn

			// Invoke the handler.
			w := httptest.NewRecorder()
			r, _ := http.NewRequest("GET", "/me", nil)
			r.Header.Set("Content-Type", "application/json")
			ctx := context.WithValue(r.Context(), "userId", uint32(1))

			httpHandler := http.HandlerFunc(h.HandleMe)
			httpHandler.ServeHTTP(w, r.WithContext(ctx))

			// Validate mock.
			if us.GetByIdInvoked != test.GetByIdInvoked {
				t.Fatalf("expected GetByIdInvoked to be %v", test.GetByIdInvoked)
			}

			expected := test.expectedResponse
			received := strings.TrimSpace(w.Body.String())

			if received != expected {
				t.Fatalf("expected %s but received %s", expected, received)
			}
		})
	}
}

func TestAuthHandler_Authentication(t *testing.T) {
	var tests = []struct {
		name                              string
		ExtractAuthenticationTokenFn      func(r *http.Request) (uint32, error)
		ExtractAuthenticationTokenInvoked bool
		expectedId                        uint32
		expectedErr                       string
	}{
		{
			name: "authenticated",
			ExtractAuthenticationTokenFn: func(r *http.Request) (uint32, error) {
				return 1, nil
			},
			ExtractAuthenticationTokenInvoked: true,
			expectedId:                        1,
			expectedErr:                       "",
		},
		{
			name: "wrong token",
			ExtractAuthenticationTokenFn: func(r *http.Request) (uint32, error) {
				return 1, errors.New("random internal error")
			},
			ExtractAuthenticationTokenInvoked: true,
			expectedId:                        0,
			expectedErr:                       `{"message":"Server Error","error":"random internal error"}`,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// Inject our mock into our handler.
			var us mock.UserService
			h := NewAuthHandler(&us)

			// Mock our Login() call.
			us.ExtractAuthenticationTokenFn = test.ExtractAuthenticationTokenFn

			// Invoke the handler.
			w := httptest.NewRecorder()
			r, _ := http.NewRequest("GET", "/test", nil)
			r.Header.Set("Content-Type", "application/json")

			nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				userId := r.Context().Value("userId").(uint32)

				if userId != test.expectedId {
					t.Fatalf("expected %v but received %v", test.expectedId, userId)
				}
			})
			h.Authentication(nextHandler).ServeHTTP(w, r)

			// Validate mock.
			if us.ExtractAuthenticationTokenInvoked != test.ExtractAuthenticationTokenInvoked {
				t.Fatalf("expected GetByIdInvoked to be %v", test.ExtractAuthenticationTokenInvoked)
			}

			// check for error
			if test.expectedId == 0 && strings.TrimSpace(w.Body.String()) != test.expectedErr {
				t.Fatalf("wrong error. expected %v but received %v", test.expectedErr, w.Body.String())
			}
		})
	}
}
