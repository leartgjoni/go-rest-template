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
type AuthHandler struct {
	// The server's base URL.
	baseUrl url.URL

	// Services
	UserService app.UserService
}

func NewAuthHandler(us app.UserService) *AuthHandler {
	return &AuthHandler{UserService: us}
}

func (h *AuthHandler) HandleSignup(w http.ResponseWriter, r *http.Request) {
	data := &payloads.UserRequest{Action: "signup"}
	if err := render.Bind(r, data); err != nil {
		utils.Render(w, r, payloads.ErrInvalidRequest(err))
		return
	}

	user := data.User

	err := h.UserService.Save(user)
	if err != nil {
		utils.Render(w, r, payloads.ErrServer(err))
		return
	}

	jwtToken, err := h.UserService.CreateToken(user)
	if err != nil {
		utils.Render(w, r, payloads.ErrServer(err))
		return
	}

	render.Status(r, http.StatusCreated)
	utils.Render(w, r, payloads.NewUserResponse(user, jwtToken))
}


func (h *AuthHandler) HandleLogin(w http.ResponseWriter, r *http.Request) {
	data := &payloads.UserRequest{Action: "login"}
	if err := render.Bind(r, data); err != nil {
		utils.Render(w, r, payloads.ErrInvalidRequest(err))
		return
	}

	user := data.User

	jwtToken, err := h.UserService.Login(user)
	if err != nil {
		utils.Render(w, r, payloads.ErrServer(err))
		return
	}

	utils.Render(w, r, payloads.NewUserResponse(user, jwtToken))
}

func (h *AuthHandler) HandleMe(w http.ResponseWriter, r *http.Request) {
	userId := r.Context().Value("userId").(uint32)


	user, err := h.UserService.GetById(userId)

	if err != nil {
		utils.Render(w, r, payloads.ErrServer(err))
		return
	}

	utils.Render(w, r, payloads.NewUserResponse(user, ""))
}

func (h *AuthHandler) Authentication(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userId, err := h.UserService.ExtractAuthenticationToken(r)
		if err != nil {
			utils.Render(w, r, payloads.ErrUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), "userId", userId)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
