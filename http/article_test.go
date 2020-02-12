package http

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"github.com/go-chi/chi"
	app "github.com/leartgjoni/go-rest-template"
	"github.com/leartgjoni/go-rest-template/mock"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestArticleHandler_HandleCreate(t *testing.T) {
	// mock time
	now := time.Unix(0, 0)
	nowString := now.Format(time.RFC3339)

	var tests = []struct {
		name             string
		SaveFn           func(a *app.Article) error
		SaveInvoked      bool
		body             []byte
		expectedResponse string
	}{
		{
			name: "success",
			SaveFn: func(a *app.Article) error {
				a.ID = uint32(1)
				a.Slug = "random-title-123456789012"
				a.CreatedAt = now
				a.UpdatedAt = now
				return nil
			},
			SaveInvoked:      true,
			body:             []byte(`{"title":"random title","body":"random body"}`),
			expectedResponse: fmt.Sprintf(`{"id":1,"slug":"random-title-123456789012","title":"random title","body":"random body","user_id":1,"created_at":"%s","updated_at":"%s"}`, nowString, nowString),
		},
		{
			name: "Save() error",
			SaveFn: func(a *app.Article) error {
				return errors.New("save fn error")
			},
			SaveInvoked:      true,
			body:             []byte(`{"title":"random title","body":"random body"}`),
			expectedResponse: `{"message":"Server Error","error":"save fn error"}`,
		},
		{
			name:             "Invalid request",
			SaveFn:           nil,
			SaveInvoked:      false,
			body:             nil,
			expectedResponse: `{"message":"Invalid request.","error":"EOF"}`,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// Inject our mock into our handler.
			var as mock.ArticleService
			h := NewArticleHandler(&as)

			// Mock our Save() call.
			as.SaveFn = test.SaveFn

			// request body
			var body = test.body

			// Invoke the handler.
			w := httptest.NewRecorder()
			r, _ := http.NewRequest("POST", "/article", bytes.NewBuffer(body))
			r.Header.Set("Content-Type", "application/json")
			ctx := context.WithValue(r.Context(), "userId", uint32(1))

			httpHandler := http.HandlerFunc(h.HandleCreate)
			httpHandler.ServeHTTP(w, r.WithContext(ctx))

			// Validate mock.
			if as.SaveInvoked != test.SaveInvoked {
				t.Fatalf("expected SaveInvoked to be %v", test.SaveInvoked)
			}

			expected := test.expectedResponse
			received := strings.TrimSpace(w.Body.String())

			if received != expected {
				t.Fatalf("expected %s but received %s", expected, received)
			}
		})
	}
}

func TestArticleHandler_HandleList(t *testing.T) {
	// mock time
	now := time.Unix(0, 0)
	nowString := now.Format(time.RFC3339)

	var tests = []struct {
		name             string
		GetAllFn         func() ([]*app.Article, error)
		GetAllInvoked    bool
		expectedResponse string
	}{
		{
			name: "success",
			GetAllFn: func() ([]*app.Article, error) {
				return []*app.Article{
					{
						ID:        1,
						Slug:      "title-one-123456789012",
						Title:     "title one",
						Body:      "body one",
						UserId:    1,
						CreatedAt: now,
						UpdatedAt: now,
					},
					{
						ID:        2,
						Slug:      "title-two-123456789012",
						Title:     "title two",
						Body:      "body two",
						UserId:    2,
						CreatedAt: now,
						UpdatedAt: now,
					},
				}, nil
			},
			GetAllInvoked:    true,
			expectedResponse: fmt.Sprintf(`[{"id":1,"slug":"title-one-123456789012","title":"title one","body":"body one","user_id":1,"created_at":"%s","updated_at":"%s"},{"id":2,"slug":"title-two-123456789012","title":"title two","body":"body two","user_id":2,"created_at":"%s","updated_at":"%s"}]`, nowString, nowString, nowString, nowString),
		},
		{
			name: "GetAll() error",
			GetAllFn: func() ([]*app.Article, error) {
				return []*app.Article{}, errors.New("getAll fn error")
			},
			GetAllInvoked:    true,
			expectedResponse: `{"message":"Server Error","error":"getAll fn error"}`,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// Inject our mock into our handler.
			var as mock.ArticleService
			h := NewArticleHandler(&as)

			// Mock our Save() call.
			as.GetAllFn = test.GetAllFn

			// Invoke the handler.
			w := httptest.NewRecorder()
			r, _ := http.NewRequest("GET", "/article", nil)
			r.Header.Set("Content-Type", "application/json")
			httpHandler := http.HandlerFunc(h.HandleList)
			httpHandler.ServeHTTP(w, r)

			// Validate mock.
			if as.GetAllInvoked != test.GetAllInvoked {
				t.Fatalf("expected GetAllInvoked to be %v", test.GetAllInvoked)
			}

			expected := test.expectedResponse
			received := strings.TrimSpace(w.Body.String())

			if received != expected {
				t.Fatalf("expected %s but received %s", expected, received)
			}
		})
	}
}

func TestArticleHandler_HandleGet(t *testing.T) {
	// mock time
	now := time.Unix(0, 0)
	nowString := now.Format(time.RFC3339)

	var tests = []struct {
		name             string
		article          *app.Article
		expectedResponse string
	}{
		{
			name: "success",
			article: &app.Article{
				ID:        1,
				Slug:      "title-one-123456789012",
				Title:     "title one",
				Body:      "body one",
				UserId:    1,
				CreatedAt: now,
				UpdatedAt: now,
			},
			expectedResponse: fmt.Sprintf(`{"id":1,"slug":"title-one-123456789012","title":"title one","body":"body one","user_id":1,"created_at":"%s","updated_at":"%s"}`, nowString, nowString),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// Inject our mock into our handler.
			var as mock.ArticleService
			h := NewArticleHandler(&as)

			// Invoke the handler.
			w := httptest.NewRecorder()
			r, _ := http.NewRequest("GET", "/article/slug", nil)
			r.Header.Set("Content-Type", "application/json")
			ctx := context.WithValue(r.Context(), "article", test.article)

			httpHandler := http.HandlerFunc(h.HandleGet)
			httpHandler.ServeHTTP(w, r.WithContext(ctx))

			expected := test.expectedResponse
			received := strings.TrimSpace(w.Body.String())

			if received != expected {
				t.Fatalf("expected %s but received %s", expected, received)
			}
		})
	}
}

func TestArticleHandler_HandleUpdate(t *testing.T) {
	// mock time
	now := time.Unix(0, 0)
	nowString := now.Format(time.RFC3339)

	var tests = []struct {
		name             string
		UpdateFn         func(a *app.Article) error
		UpdateInvoked    bool
		body             []byte
		expectedResponse string
	}{
		{
			name: "success",
			UpdateFn: func(a *app.Article) error {
				a.Slug = "random-title-updated-123456789012"
				a.UpdatedAt = now
				a.CreatedAt = now
				return nil
			},
			UpdateInvoked:    true,
			body:             []byte(`{"title":"random title updated","body":"random body updated"}`),
			expectedResponse: fmt.Sprintf(`{"id":1,"slug":"random-title-updated-123456789012","title":"random title updated","body":"random body updated","user_id":1,"created_at":"%s","updated_at":"%s"}`, nowString, nowString),
		},
		{
			name: "Update() error",
			UpdateFn: func(a *app.Article) error {
				return errors.New("update fn error")
			},
			UpdateInvoked:    true,
			body:             []byte(`{"title":"random title updated","body":"random body updated"}`),
			expectedResponse: `{"message":"Server Error","error":"update fn error"}`,
		},
		{
			name:             "Invalid request",
			UpdateFn:         nil,
			UpdateInvoked:    false,
			body:             nil,
			expectedResponse: `{"message":"Invalid request.","error":"EOF"}`,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// Inject our mock into our handler.
			var as mock.ArticleService
			h := NewArticleHandler(&as)

			// Mock our Save() call.
			as.UpdateFn = test.UpdateFn

			// request body
			var body = test.body

			// Invoke the handler.
			w := httptest.NewRecorder()
			r, _ := http.NewRequest("PATCH", "/article/slug", bytes.NewBuffer(body))
			r.Header.Set("Content-Type", "application/json")
			ctx := context.WithValue(r.Context(), "article", &app.Article{ID: uint32(1), Title: "random title", Body: "random body", Slug: "random-title-123456789012", CreatedAt: now})
			ctx = context.WithValue(ctx, "userId", uint32(1))

			httpHandler := http.HandlerFunc(h.HandleUpdate)
			httpHandler.ServeHTTP(w, r.WithContext(ctx))

			// Validate mock.
			if as.UpdateInvoked != test.UpdateInvoked {
				t.Fatalf("expected UpdateInvoked to be %v", test.UpdateInvoked)
			}

			expected := test.expectedResponse
			received := strings.TrimSpace(w.Body.String())

			if received != expected {
				t.Fatalf("expected %s but received %s", expected, received)
			}
		})
	}
}

func TestArticleHandler_HandleDelete(t *testing.T) {
	var tests = []struct {
		name             string
		DeleteFn         func(slug string) error
		DeleteInvoked    bool
		expectedResponse string
	}{
		{
			name: "success",
			DeleteFn: func(slug string) error {
				return nil
			},
			DeleteInvoked:    true,
			expectedResponse: "",
		},
		{
			name: "Delete() error",
			DeleteFn: func(slug string) error {
				return errors.New("delete fn error")
			},
			DeleteInvoked:    true,
			expectedResponse: `{"message":"Server Error","error":"delete fn error"}`,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// Inject our mock into our handler.
			var as mock.ArticleService
			h := NewArticleHandler(&as)

			// Mock our Delete() call.
			as.DeleteFn = test.DeleteFn

			// Invoke the handler.
			w := httptest.NewRecorder()
			r, _ := http.NewRequest("DELETE", "/article/slug", nil)
			r.Header.Set("Content-Type", "application/json")
			ctx := context.WithValue(r.Context(), "article", &app.Article{Slug: "random-slug"})

			httpHandler := http.HandlerFunc(h.HandleDelete)
			httpHandler.ServeHTTP(w, r.WithContext(ctx))

			// Validate mock.
			if as.DeleteInvoked != test.DeleteInvoked {
				t.Fatalf("expected DeleteInvoked to be %v", test.DeleteInvoked)
			}

			expected := test.expectedResponse
			received := strings.TrimSpace(w.Body.String())

			if received != expected {
				t.Fatalf("expected %s but received %s", expected, received)
			}
		})
	}
}

func TestArticleHandler_ArticleCtx(t *testing.T) {
	// mock time
	now := time.Unix(0, 0)

	var tests = []struct {
		name             string
		GetBySlugFn      func(slug string) (*app.Article, error)
		GetBySlugInvoked bool
		expectedArticle  app.Article
		expectedErr      string
	}{
		{
			name: "success",
			GetBySlugFn: func(slug string) (*app.Article, error) {
				return &app.Article{
					ID:        1,
					Slug:      "slug",
					Title:     "title",
					Body:      "body",
					UserId:    1,
					CreatedAt: now,
					UpdatedAt: now,
				}, nil
			},
			GetBySlugInvoked: true,
			expectedArticle: app.Article{
				ID:        1,
				Slug:      "slug",
				Title:     "title",
				Body:      "body",
				UserId:    1,
				CreatedAt: now,
				UpdatedAt: now,
			},
			expectedErr: "",
		},
		{
			name: "get by slug error",
			GetBySlugFn: func(slug string) (*app.Article, error) {
				return &app.Article{}, app.ErrArticleNotFound
			},
			GetBySlugInvoked: true,
			expectedArticle:  app.Article{},
			expectedErr:      `{"message":"Resource not found."}`,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// Inject our mock into our handler.
			var as mock.ArticleService
			h := NewArticleHandler(&as)

			// Mock our GetBySlug() call.
			as.GetBySlugFn = test.GetBySlugFn

			// Invoke the handler.
			w := httptest.NewRecorder()
			r, _ := http.NewRequest("GET", "/{articleSlug}", nil)
			r.Header.Set("Content-Type", "application/json")
			chiCtx := chi.NewRouteContext()
			chiCtx.URLParams.Add("key", "value")

			r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, chiCtx))

			nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				article := r.Context().Value("article").(*app.Article)

				if article.ID != test.expectedArticle.ID ||
					article.Slug != test.expectedArticle.Slug ||
					article.Title != test.expectedArticle.Title ||
					article.Body != test.expectedArticle.Body ||
					article.UserId != test.expectedArticle.UserId ||
					!article.CreatedAt.Equal(test.expectedArticle.CreatedAt) ||
					!article.UpdatedAt.Equal(test.expectedArticle.UpdatedAt) {
					t.Fatalf("expected %v but received %v", test.expectedArticle, article)
				}
			})
			h.ArticleCtx(nextHandler).ServeHTTP(w, r)

			// Validate mock.
			if as.GetBySlugInvoked != test.GetBySlugInvoked {
				t.Fatalf("expected GetBySlugInvoked to be %v", test.GetBySlugInvoked)
			}

			// check for error
			if test.expectedArticle.ID == 0 && strings.TrimSpace(w.Body.String()) != test.expectedErr {
				t.Fatalf("wrong error. expected %v but received %v", test.expectedErr, w.Body.String())
			}
		})
	}
}

func TestArticleHandler_ArticleOwner(t *testing.T) {
	var tests = []struct {
		name             string
		userId           uint32
		article          *app.Article
		expectedErr      string
		expectedResponse string
	}{
		{
			name:             "is owner",
			userId:           1,
			article:          &app.Article{UserId: 1},
			expectedResponse: "",
		},
		{
			name:             "not owner",
			userId:           1,
			article:          &app.Article{UserId: 2},
			expectedResponse: `{"message":"Unauthorized"}`,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// Inject our mock into our handler.
			var as mock.ArticleService
			h := NewArticleHandler(&as)

			// Invoke the handler.
			w := httptest.NewRecorder()
			r, _ := http.NewRequest("GET", "/{articleSlug}", nil)
			r.Header.Set("Content-Type", "application/json")
			ctx := context.WithValue(r.Context(), "userId", test.userId)
			ctx = context.WithValue(ctx, "article", test.article)

			nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
			h.ArticleOwner(nextHandler).ServeHTTP(w, r.WithContext(ctx))

			expected := test.expectedResponse
			received := strings.TrimSpace(w.Body.String())

			if received != expected {
				t.Fatalf("expected %s but received %s", expected, received)
			}
		})
	}
}
