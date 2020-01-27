package http_test

import (
	"bytes"
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
		CreateTokenFn      func(user *app.User) (string, error)
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
			CreateTokenFn: func(user *app.User) (string, error) {
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
			CreateTokenFn: func(user *app.User) (string, error) {
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
