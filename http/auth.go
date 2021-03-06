package http

import (
	"context"
	"github.com/go-chi/render"
	app "github.com/leartgjoni/go-rest-template"
	"github.com/leartgjoni/go-rest-template/http/payloads"
	"github.com/leartgjoni/go-rest-template/http/utils"
	"net/http"
	"net/url"
)

// AuthHandler represents an HTTP handler for managing authentication.
type AuthHandler interface {
	HandleSignup(w http.ResponseWriter, r *http.Request)
	HandleLogin(w http.ResponseWriter, r *http.Request)
	HandleMe(w http.ResponseWriter, r *http.Request)
	Authentication(next http.Handler) http.Handler
}

// struct that implements interface
type authHandler struct {
	// The server's base URL.
	baseUrl url.URL

	// Services
	UserService app.UserService
}

func NewAuthHandler(us app.UserService) *authHandler {
	return &authHandler{UserService: us}
}

func (h *authHandler) HandleSignup(w http.ResponseWriter, r *http.Request) {
	data := &payloads.UserRequest{Action: "signup"}
	if err := render.Bind(r, data); err != nil {
		utils.Render(w, r, payloads.ErrInvalidRequest(err))
		return
	}

	user := data.User

	err := h.UserService.Save(user)
	if err != nil {
		utils.Render(w, r, authHttpError(err))
		return
	}

	jwtToken, err := h.UserService.CreateToken(user.ID)
	if err != nil {
		utils.Render(w, r, authHttpError(err))
		return
	}

	render.Status(r, http.StatusCreated)
	utils.Render(w, r, payloads.NewUserResponse(user, jwtToken))
}

func (h *authHandler) HandleLogin(w http.ResponseWriter, r *http.Request) {
	data := &payloads.UserRequest{Action: "login"}
	if err := render.Bind(r, data); err != nil {
		utils.Render(w, r, payloads.ErrInvalidRequest(err))
		return
	}

	user := data.User

	jwtToken, err := h.UserService.Login(user)
	if err != nil {
		utils.Render(w, r, authHttpError(err))
		return
	}

	utils.Render(w, r, payloads.NewUserResponse(user, jwtToken))
}

func (h *authHandler) HandleMe(w http.ResponseWriter, r *http.Request) {
	userId := r.Context().Value("userId").(uint32)

	user, err := h.UserService.GetById(userId)

	if err != nil {
		utils.Render(w, r, authHttpError(err))
		return
	}

	utils.Render(w, r, payloads.NewUserResponse(user, ""))
}

func (h *authHandler) Authentication(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userId, err := h.UserService.ExtractAuthenticationToken(r)
		if err != nil {
			utils.Render(w, r, authHttpError(err))
			return
		}

		ctx := context.WithValue(r.Context(), "userId", userId)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// app error to http error
func authHttpError(err error) render.Renderer {
	switch err {
	case app.ErrEmailAlreadyUsed,
		app.ErrWrongPasswordFormat:
		return payloads.ErrInvalidRequest(err)
	case app.ErrWrongCredentials:
		return payloads.ErrUnauthorized
	case app.ErrUserNotFound:
		return payloads.ErrNotFound
	default:
		return payloads.ErrServer(err)
	}
}
