package mock

import (
	"net/http"
)

// ArticleHandler represents a mock implementation of http.IArticleHandler.
type ArticleHandler struct {
	Invoked []string
}

func (h *ArticleHandler) HandleCreate(w http.ResponseWriter, r *http.Request) {
	h.Invoked = append(h.Invoked, "HandleCreate")
}
func (h *ArticleHandler) HandleList(w http.ResponseWriter, r *http.Request) {
	h.Invoked = append(h.Invoked, "HandleList")
}
func (h *ArticleHandler) HandleGet(w http.ResponseWriter, r *http.Request) {
	h.Invoked = append(h.Invoked, "HandleGet")
}
func (h *ArticleHandler) HandleUpdate(w http.ResponseWriter, r *http.Request) {
	h.Invoked = append(h.Invoked, "HandleUpdate")
}
func (h *ArticleHandler) HandleDelete(w http.ResponseWriter, r *http.Request) {
	h.Invoked = append(h.Invoked, "HandleDelete")
}
func (h *ArticleHandler) ArticleCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		h.Invoked = append(h.Invoked, "ArticleCtx")
		next.ServeHTTP(w, r)
	})
}
func (h *ArticleHandler) ArticleOwner(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		h.Invoked = append(h.Invoked, "ArticleOwner")
		next.ServeHTTP(w, r)
	})
}