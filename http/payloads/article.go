package payloads

import (
	"errors"
	"github.com/go-chi/render"
	app "github.com/leartgjoni/go-rest-template"
	"net/http"
	"strings"
	"time"
)

type ArticleRequest struct {
	*app.Article

	Action string
}

func (a *ArticleRequest) Bind(r *http.Request) error {
	// a.Article is nil if no Article fields are sent in the request. Return an
	// error to avoid a nil pointer dereference.
	if a.Article == nil {
		return errors.New("missing required Article fields")
	}

	//post-process after a decode
	if a.Action == "create" {
		a.Prepare()
	} else if a.Action == "update" {
		ctxArticle := r.Context().Value("article").(*app.Article)
		a.Slug = ctxArticle.Slug
		a.CreatedAt = ctxArticle.CreatedAt
		a.ID = ctxArticle.ID
		a.UpdatedAt = time.Now()
	}
	a.UserId = r.Context().Value("userId").(uint32)
	return a.Validate(a.Action)
}

// response
type ArticleResponse struct {
	*app.Article
}

func (rd *ArticleResponse) Render(http.ResponseWriter, *http.Request) error {
	return nil
}

func NewArticleResponse(article *app.Article) *ArticleResponse {
	return &ArticleResponse{Article: article}
}

func NewArticleListResponse(articles []*app.Article) []render.Renderer {
	var list []render.Renderer
	for _, article := range articles {
		list = append(list, NewArticleResponse(article))
	}
	return list
}

func (a *ArticleRequest) Prepare() {
	//a.ID = 0
	a.CreatedAt = time.Now()
	a.UpdatedAt = time.Now()
}

func (a *ArticleRequest) Validate(action string) error {
	switch strings.ToLower(action) {
	case "create":
		if a.Title == "" {
			return errors.New("required title")
		}
		if a.Body == "" {
			return errors.New("required body")
		}
		return nil
	case "update":
		if a.Title == "" {
			return errors.New("required title")
		}
		if a.Body == "" {
			return errors.New("required body")
		}
		return nil
	default:
		return nil
	}
}
