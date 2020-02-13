package payloads

import (
	"errors"
	app "github.com/leartgjoni/go-rest-template"
	"testing"
)

func TestUserRequest_Bind(t *testing.T) {
	t.Run("login", testUserLogin)
	t.Run("signup", testUserSignup)
}

func testUserLogin(t *testing.T) {
	tests := []struct {
		name        string
		user        *app.User
		expectedErr error
	}{
		{
			name:        "missing required User fields",
			user:        nil,
			expectedErr: errors.New("missing required User fields"),
		},
		{
			name:        "required password",
			user:        &app.User{Email: "test@test.com"},
			expectedErr: errors.New("required password"),
		},
		{
			name:        "no email",
			user:        &app.User{Password: "random-password"},
			expectedErr: errors.New("required email"),
		},
		{
			name:        "invalid email",
			user:        &app.User{Password: "random-password", Email: "test-random"},
			expectedErr: errors.New("invalid email"),
		},
		{
			name:        "correct",
			user:        &app.User{Password: "random-password", Email: "test@random.com"},
			expectedErr: nil,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			r := UserRequest{Action: "login", User: test.user}
			err := r.Bind(nil)

			if test.expectedErr != nil {
				if err == nil || err.Error() != test.expectedErr.Error() {
					t.Fatalf("wrong error. expected %s but got %s", test.expectedErr, err)
				}
			} else {
				if err != nil {
					t.Fatalf("wrong error. expected %s but got %s", test.expectedErr, err)
				}
			}

		})
	}
}

func testUserSignup(t *testing.T) {
	tests := []struct {
		name        string
		user        *app.User
		expectedErr error
	}{
		{
			name:        "missing required User fields",
			user:        nil,
			expectedErr: errors.New("missing required User fields"),
		},
		{
			name:        "required password",
			user:        &app.User{Email: "test@test.com", Username: "test"},
			expectedErr: errors.New("required password"),
		},
		{
			name:        "no email",
			user:        &app.User{Password: "random-password", Username: "test"},
			expectedErr: errors.New("required email"),
		},
		{
			name:        "invalid email",
			user:        &app.User{Password: "random-password", Email: "test-random", Username: "test"},
			expectedErr: errors.New("invalid email"),
		},
		{
			name:        "required username",
			user:        &app.User{Password: "random-password", Email: "test@random.com"},
			expectedErr: errors.New("required username"),
		},
		{
			name:        "correct",
			user:        &app.User{Password: "random-password", Email: "test@random.com", Username: "test"},
			expectedErr: nil,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			r := UserRequest{Action: "signup", User: test.user}
			err := r.Bind(nil)

			if test.expectedErr != nil {
				if err == nil || err.Error() != test.expectedErr.Error() {
					t.Fatalf("wrong error. expected %s but got %s", test.expectedErr, err)
				}
			} else {
				if err != nil {
					t.Fatalf("wrong error. expected %s but got %s", test.expectedErr, err)
				}
			}

		})
	}

	// test prepare()
	t.Run("prepare", func(t *testing.T) {
		r := UserRequest{Action: "signup", User: &app.User{Email: " test>@test.com ", Username: " te>st "}}
		r.prepare()

		if r.User.Email != "test&gt;@test.com" || r.User.Username != "te&gt;st" || r.User.CreatedAt.IsZero() || r.User.UpdatedAt.IsZero() {
			t.Fatalf("incorrect prepare: %v", r.User)
		}
	})
}
