package postgres

import (
	"errors"
	"github.com/DATA-DOG/go-sqlmock"
	app "github.com/leartgjoni/go-rest-template"
	"regexp"
	"testing"
	"time"
)

func TestArticleService_GetAll(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	articles := []*app.Article{
		{
			ID:        1,
			Slug:      "slug-1",
			Title:     "title 1",
			Body:      "body 1",
			UserId:    1,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			ID:        2,
			Slug:      "slug-2",
			Title:     "title 2",
			Body:      "body 2",
			UserId:    2,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}

	tests := []struct {
		name      string
		sqlResult *sqlmock.Rows
		error     error
		result    []*app.Article
	}{
		{
			name: "normal case",
			sqlResult: sqlmock.NewRows([]string{"id", "slug", "title", "body", "user_id", "created_at", "updated_at"}).
				AddRow(articles[0].ID, articles[0].Slug, articles[0].Title, articles[0].Body, articles[0].UserId, articles[0].CreatedAt, articles[0].UpdatedAt).
				AddRow(articles[1].ID, articles[1].Slug, articles[1].Title, articles[1].Body, articles[1].UserId, articles[1].CreatedAt, articles[1].UpdatedAt),
			error:  nil,
			result: articles,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			mock.ExpectQuery("^SELECT (.+) FROM articles").WillReturnRows(test.sqlResult)

			as := NewArticleService(&DB{db})

			articles, err = as.GetAll()

			// we make sure that all expectations were met
			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}

			if err != nil && err.Error() != test.error.Error() {
				t.Fatalf("wrong error. expected %s but got %s", test.error, err)
			}

			for i, a := range articles {
				if a.ID != test.result[i].ID ||
					a.Slug != test.result[i].Slug ||
					a.Title != test.result[i].Title ||
					a.Body != test.result[i].Body ||
					a.UserId != test.result[i].UserId ||
					!a.CreatedAt.Equal(test.result[i].CreatedAt) ||
					!a.UpdatedAt.Equal(test.result[i].UpdatedAt) {
					t.Fatalf("wrong article. expected %v but got %v", test.result[i], a)
				}
			}
		})
	}
}

func TestArticleService_GetBySlug(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	article := app.Article{
		ID:        1,
		Slug:      "slug-1",
		Title:     "title 1",
		Body:      "body 1",
		UserId:    1,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	tests := []struct {
		name         string
		sqlResult    *sqlmock.Rows
		insertResult *sqlmock.Rows
		error        error
		result       app.Article
	}{
		{
			name: "normal case",
			sqlResult: sqlmock.NewRows([]string{"id", "slug", "title", "body", "user_id", "created_at", "updated_at"}).
				AddRow(article.ID, article.Slug, article.Title, article.Body, article.UserId, article.CreatedAt, article.UpdatedAt),
			error:  nil,
			result: article,
		},
		{
			name:      "article not found",
			sqlResult: sqlmock.NewRows([]string{"id", "slug", "title", "body", "user_id", "created_at", "updated_at"}),
			error:     app.ErrArticleNotFound,
			result:    app.Article{},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			mock.ExpectQuery("^SELECT (.+) FROM articles WHERE slug=*").WillReturnRows(test.sqlResult)

			as := NewArticleService(&DB{db})

			result, err := as.GetBySlug("random-slug")

			// we make sure that all expectations were met
			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}

			if err != nil && err.Error() != test.error.Error() {
				t.Fatalf("wrong error. expected %s but got %s", test.error, err)
			}

			if (result == nil && err == nil) || result.ID != test.result.ID ||
				result.Slug != test.result.Slug ||
				result.Title != test.result.Title ||
				result.Body != test.result.Body ||
				result.UserId != test.result.UserId ||
				!result.CreatedAt.Equal(test.result.CreatedAt) ||
				!result.UpdatedAt.Equal(test.result.UpdatedAt) {
				t.Fatalf("wrong article. expected %v but got %v", test.result, result)
			}
		})
	}
}

func TestArticleService_Save(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	article := app.Article{
		Title:     "title 1",
		Body:      "body 1",
		UserId:    1,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	tests := []struct {
		name         string
		sqlResult    *sqlmock.Rows
		insertResult *sqlmock.Rows
		error        error
		article      app.Article
		slugRegex    string
	}{
		{
			name: "normal case",
			sqlResult: sqlmock.NewRows([]string{"id"}).
				AddRow(1),
			error:     nil,
			article:   article,
			slugRegex: `title-1-[a-zA-Z0-9]{12}`,
		},
		{
			name: "save failed",
			sqlResult: sqlmock.NewRows([]string{"id"}).
				AddRow(0),
			error: errors.New("unable to save"),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			mock.ExpectQuery("^INSERT INTO (.+) VALUES (.+) RETURNING id").WillReturnRows(test.sqlResult)

			as := NewArticleService(&DB{db})

			err := as.Save(&test.article)

			// we make sure that all expectations were met
			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}

			if err != nil && err.Error() != test.error.Error() {
				t.Fatalf("wrong error. expected %s but got %s", test.error, err)
			}

			if err == nil && (test.article.ID != 1 || !regexp.MustCompile(test.slugRegex).Match([]byte(test.article.Slug))) {
				t.Fatal("save error. expected article id and slug to be set")
			}
		})
	}
}
