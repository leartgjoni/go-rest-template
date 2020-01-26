package http

import (
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"net/http"
)

func (s *Server) router() http.Handler {
	r := chi.NewRouter()

	// Attach router middleware.
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// Create API routes.
	r.Route("/", func(r chi.Router) {
		r.Get("/health", s.handlePing)

		r.Route("/auth", func (r chi.Router) {
			r.Post("/signup", s.authHandler.HandleSignup)
			r.Post("/login", s.authHandler.HandleLogin)
			r.With(s.authHandler.Authentication).Get("/me", s.authHandler.HandleMe)
		})

		r.Route("/articles", func (r chi.Router) {
			r.Get("/", s.articleHandler.HandleList)
			r.With(s.articleHandler.ArticleCtx).Get("/{articleSlug}", s.articleHandler.HandleGet)
			r.Route("/", func (r chi.Router) {
				r.Use(s.authHandler.Authentication)
				r.Post("/", s.articleHandler.HandleCreate)
				r.Route("/{articleSlug}", func (r chi.Router) {
					r.Use(s.articleHandler.ArticleCtx, s.articleHandler.ArticleOwner)

					r.Patch("/", s.articleHandler.HandleUpdate)
					r.Delete("/", s.articleHandler.HandleDelete)
				})
			})
		})
	})

	return r
}
