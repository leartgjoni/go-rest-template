package mock

import "net/http"

type AuthHandler struct {
	Invoked []string
}

func (h *AuthHandler) HandleSignup(w http.ResponseWriter, r *http.Request) {
	h.Invoked = append(h.Invoked, "HandleSignup")
}
func (h *AuthHandler) HandleLogin(w http.ResponseWriter, r *http.Request) {
	h.Invoked = append(h.Invoked, "HandleLogin")
}
func (h *AuthHandler) HandleMe(w http.ResponseWriter, r *http.Request) {
	h.Invoked = append(h.Invoked, "HandleMe")
}
func (h *AuthHandler) Authentication(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		h.Invoked = append(h.Invoked, "Authentication")
		next.ServeHTTP(w, r)
	})
}