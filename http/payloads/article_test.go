package payloads

import (
	"context"
	"errors"
	app "github.com/leartgjoni/go-rest-template"
	"net/http"
	"testing"
	"time"
)

func TestArticleRequest_Bind(t *testing.T) {
	t.Run("create", testArticleCreate)
	t.Run("update", testArticleUpdate)
}

func testArticleCreate(t *testing.T) {
	tests := []struct {
		name        string
		article     *app.Article
		expectedErr error
		userId      uint32
	}{
		{
			name:        "missing required User fields",
			article:     nil,
			expectedErr: errors.New("missing required Article fields"),
		},
		{
			name:        "no title",
			article:     &app.Article{Body: "random body"},
			expectedErr: errors.New("required title"),
		},
		{
			name:        "no body",
			article:     &app.Article{Title: "random title"},
			expectedErr: errors.New("required body"),
		},
		{
			name:        "correct",
			article:     &app.Article{Title: "random title", Body: "random body"},
			expectedErr: nil,
			userId:      1,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			a := ArticleRequest{Action: "create", Article: test.article}
			rq, _ := http.NewRequest("GET", "/", nil)
			ctx := context.WithValue(rq.Context(), "userId", test.userId)

			err := a.Bind(rq.WithContext(ctx))

			if test.expectedErr != nil {
				if err == nil || err.Error() != test.expectedErr.Error() {
					t.Fatalf("wrong error. expected %s but got %s", test.expectedErr, err)
				}
			} else {
				if err != nil {
					t.Fatalf("wrong error. expected %s but got %s", test.expectedErr, err)
				}

				if a.UserId != test.userId {
					t.Fatal("userId was not extracted properly in context")
				}
			}

		})
	}

	// test prepare()
	t.Run("prepare", func(t *testing.T) {
		r := ArticleRequest{Action: "create", Article: &app.Article{Title: "random title", Body: "random body"}}
		r.prepare()

		if r.Article.CreatedAt.IsZero() || r.Article.UpdatedAt.IsZero() {
			t.Fatalf("incorrect prepare: %v", r.Article)
		}
	})
}

func testArticleUpdate(t *testing.T) {
	contextArticle := &app.Article{
		ID:        1,
		Slug:      "slug",
		Title:     "title",
		Body:      "body",
		CreatedAt: time.Now(),
	}
	tests := []struct {
		name        string
		article     *app.Article
		expectedErr error
		userId      uint32
		ctxArticle  *app.Article
	}{
		{
			name:        "missing required User fields",
			article:     nil,
			expectedErr: errors.New("missing required Article fields"),
			ctxArticle:  contextArticle,
		},
		{
			name:        "no title",
			article:     &app.Article{Body: "random body"},
			expectedErr: errors.New("required title"),
			ctxArticle:  contextArticle,
		},
		{
			name:        "no body",
			article:     &app.Article{Title: "random title"},
			expectedErr: errors.New("required body"),
			ctxArticle:  contextArticle,
		},
		{
			name:        "correct",
			article:     &app.Article{Title: "random title", Body: "random body"},
			expectedErr: nil,
			userId:      1,
			ctxArticle:  contextArticle,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			a := ArticleRequest{Action: "update", Article: test.article}
			rq, _ := http.NewRequest("GET", "/", nil)
			ctx := context.WithValue(rq.Context(), "userId", test.userId)
			ctx = context.WithValue(ctx, "article", test.ctxArticle)

			err := a.Bind(rq.WithContext(ctx))

			if test.expectedErr != nil {
				if err == nil || err.Error() != test.expectedErr.Error() {
					t.Fatalf("wrong error. expected %s but got %s", test.expectedErr, err)
				}
			} else {
				if err != nil {
					t.Fatalf("wrong error. expected %s but got %s", test.expectedErr, err)
				}

				if a.UserId != test.userId {
					t.Fatal("userId was not extracted properly in context")
				}

				if a.Slug != test.ctxArticle.Slug || a.ID != test.ctxArticle.ID || a.CreatedAt != test.ctxArticle.CreatedAt || a.UpdatedAt.Before(test.ctxArticle.CreatedAt) {
					t.Fatal("ctxArticle was not extracted properly in context")
				}
			}

		})
	}
}
