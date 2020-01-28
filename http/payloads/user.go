package payloads

import (
	"errors"
	"github.com/badoux/checkmail"
	app "github.com/leartgjoni/go-rest-template"
	"html"
	"net/http"
	"strings"
	"time"
)

type UserRequest struct {
	*app.User

	Action string // application-level action, helps in controlling logic flow
}

func (u *UserRequest) Bind(*http.Request) error {
	if u.User == nil {
		return errors.New("missing required User fields")
	}

	//post-process after a decode
	if u.Action == "signup" {
		u.prepare()
	}
	return u.validate(u.Action)
}

// response
type UserResponse struct {
	*app.User

	Password string `json:"password,omitempty"` // remove password from response

	Token string `json:"token,omitempty"` // add token to response
}

func NewUserResponse(user *app.User, token string) *UserResponse {
	return &UserResponse{User: user, Token: token}
}

func (rd *UserResponse) Render(http.ResponseWriter, *http.Request) error {
	return nil
}

func (u *UserRequest) prepare() {
	u.Username = html.EscapeString(strings.TrimSpace(u.Username))
	u.Email = html.EscapeString(strings.TrimSpace(u.Email))
	u.CreatedAt = time.Now()
	u.UpdatedAt = time.Now()
}

func (u *UserRequest) validate(action string) error {
	switch strings.ToLower(action) {
	case "signup":
		if u.Username == "" {
			return errors.New("required username")
		}
		if u.Password == "" {
			return errors.New("required password")
		}
		if u.Email == "" {
			return errors.New("required email")
		}
		if err := checkmail.ValidateFormat(u.Email); err != nil {
			return errors.New("invalid email")
		}
		return nil
	case "login":
		if u.Password == "" {
			return errors.New("required password")
		}
		if u.Email == "" {
			return errors.New("required email")
		}
		if err := checkmail.ValidateFormat(u.Email); err != nil {
			return errors.New("invalid email")
		}
		return nil
	default:
		return nil
	}
}
