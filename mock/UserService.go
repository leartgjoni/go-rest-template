package mock

import (
	app "github.com/leartgjoni/go-rest-template"
	"net/http"
)

// UserService represents a mock implementation of app.UserService.
type UserService struct {
	CreateTokenFn func(user *app.User) (string, error)
	CreateTokenInvoked bool

	ExtractAuthenticationTokenFn func (r *http.Request) (uint32, error)
	ExtractAuthenticationTokenInvoked bool

	SaveFn func (user *app.User) error
	SaveInvoked bool

	GetByIdFn func (userId uint32) (*app.User, error)
	GetByIdInvoked bool

	LoginFn func (u *app.User) (string, error)
	LoginInvoked bool
}

// CreateToken invokes the mock implementation and marks the function as invoked.
func (s *UserService) CreateToken(user *app.User) (string, error) {
	s.CreateTokenInvoked = true
	return s.CreateTokenFn(user)
}

// ExtractAuthenticationToken invokes the mock implementation and marks the function as invoked.
func (s *UserService) ExtractAuthenticationToken(r *http.Request) (uint32, error) {
	s.ExtractAuthenticationTokenInvoked = true
	return s.ExtractAuthenticationToken(r)
}

// Save invokes the mock implementation and marks the function as invoked.
func (s *UserService) Save(user *app.User) error {
	s.SaveInvoked = true
	return s.SaveFn(user)
}

// GetById invokes the mock implementation and marks the function as invoked.
func (s *UserService) GetById(userId uint32) (*app.User, error) {
	s.GetByIdInvoked = true
	return s.GetByIdFn(userId)
}

// Login invokes the mock implementation and marks the function as invoked.
func (s *UserService) Login(u *app.User) (string, error) {
	s.LoginInvoked = true
	return s.LoginFn(u)
}