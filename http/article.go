package http

import (
	"context"
	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	app "github.com/leartgjoni/go-rest-template"
	"github.com/leartgjoni/go-rest-template/http/payloads"
	"github.com/leartgjoni/go-rest-template/http/utils"
	"net/http"
	"net/url"
)

// AuthHandler represents an HTTP handler for managing authentication.
type ArticleHandler struct {
	// The server's base URL.
	baseUrl url.URL

	// Services
	ArticleService app.ArticleService
}

func NewArticleHandler(as app.ArticleService) *ArticleHandler {
	return &ArticleHandler{ArticleService: as}
}

func (h *ArticleHandler) HandleCreate(w http.ResponseWriter, r *http.Request) {
	data := &payloads.ArticleRequest{Action: "create"}
	if err := render.Bind(r, data); err != nil {
		utils.Render(w, r, payloads.ErrInvalidRequest(err))
		return
	}

	article := data.Article

	err := h.ArticleService.Save(article)
	if err != nil {
		utils.Render(w, r, payloads.ErrServer(err))
	}

	render.Status(r, http.StatusCreated)
	utils.Render(w, r, payloads.NewArticleResponse(article))
}

func (h *ArticleHandler) HandleList(w http.ResponseWriter, r *http.Request) {
	articles, err := h.ArticleService.GetAll()
	if err != nil {
		utils.Render(w, r, payloads.ErrServer(err))
		return
	}

	utils.RenderList(w, r, payloads.NewArticleListResponse(articles))
}

func (h *ArticleHandler) HandleGet(w http.ResponseWriter, r *http.Request) {
	article := r.Context().Value("article").(*app.Article)

	utils.Render(w, r, payloads.NewArticleResponse(article))
}

func (h *ArticleHandler) HandleUpdate(w http.ResponseWriter, r *http.Request) {
	data := &payloads.ArticleRequest{Action: "update"}
	if err := render.Bind(r, data); err != nil {
		utils.Render(w, r, payloads.ErrInvalidRequest(err))
		return
	}

	article := data.Article

	err := h.ArticleService.Update(article)
	if err != nil {
		utils.Render(w, r, payloads.ErrServer(err))
		return
	}

	utils.Render(w, r, payloads.NewArticleResponse(article))
}

func (h *ArticleHandler) HandleDelete(w http.ResponseWriter, r *http.Request) {
	article := r.Context().Value("article").(*app.Article)

	err := h.ArticleService.Delete(article.Slug)

	if err != nil {
		utils.Render(w, r, payloads.ErrServer(err))
		return
	}
}

// middlewares
func (h *ArticleHandler) ArticleCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		articleSlug := chi.URLParam(r, "articleSlug")
		article, err := h.ArticleService.GetBySlug(articleSlug)
		if err != nil {
			utils.Render(w, r, payloads.ErrNotFound)
			return
		}

		ctx := context.WithValue(r.Context(), "article", article)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// check that the requester is the owner of the article
func (h *ArticleHandler) ArticleOwner(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		article := r.Context().Value("article").(*app.Article)
		userId := r.Context().Value("userId").(uint32)

		if article.UserId != userId {
			utils.Render(w, r, payloads.ErrUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}
