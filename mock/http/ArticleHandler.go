package mock

import (
	"net/http"
)

// ArticleHandler represents a mock implementation of http.IArticleHandler.
type ArticleHandler struct {
	Invoked *[]string
}

func NewMockArticleHandler(invoked *[]string) *ArticleHandler {
	return &ArticleHandler{invoked}
}

func (h *ArticleHandler) HandleCreate(w http.ResponseWriter, r *http.Request) {
	*h.Invoked = append(*h.Invoked, "ArticleHandler.HandleCreate")
}
func (h *ArticleHandler) HandleList(w http.ResponseWriter, r *http.Request) {
	*h.Invoked = append(*h.Invoked, "ArticleHandler.HandleList")
}
func (h *ArticleHandler) HandleGet(w http.ResponseWriter, r *http.Request) {
	*h.Invoked = append(*h.Invoked, "ArticleHandler.HandleGet")
}
func (h *ArticleHandler) HandleUpdate(w http.ResponseWriter, r *http.Request) {
	*h.Invoked = append(*h.Invoked, "ArticleHandler.HandleUpdate")
}
func (h *ArticleHandler) HandleDelete(w http.ResponseWriter, r *http.Request) {
	*h.Invoked = append(*h.Invoked, "ArticleHandler.HandleDelete")
}
func (h *ArticleHandler) ArticleCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		*h.Invoked = append(*h.Invoked, "ArticleHandler.ArticleCtx")
		next.ServeHTTP(w, r)
	})
}
func (h *ArticleHandler) ArticleOwner(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		*h.Invoked = append(*h.Invoked, "ArticleHandler.ArticleOwner")
		next.ServeHTTP(w, r)
	})
}
