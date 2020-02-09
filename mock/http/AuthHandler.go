package mock

import "net/http"

type AuthHandler struct {
	Invoked *[]string
}

func NewMockAuthHandler(invoked *[]string) *AuthHandler {
	return &AuthHandler{invoked}
}

func (h *AuthHandler) HandleSignup(w http.ResponseWriter, r *http.Request) {
	*h.Invoked = append(*h.Invoked, "AuthHandler.HandleSignup")
}
func (h *AuthHandler) HandleLogin(w http.ResponseWriter, r *http.Request) {
	*h.Invoked = append(*h.Invoked, "AuthHandler.HandleLogin")
}
func (h *AuthHandler) HandleMe(w http.ResponseWriter, r *http.Request) {
	*h.Invoked = append(*h.Invoked, "AuthHandler.HandleMe")
}
func (h *AuthHandler) Authentication(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		*h.Invoked = append(*h.Invoked, "AuthHandler.Authentication")
		next.ServeHTTP(w, r)
	})
}
