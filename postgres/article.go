package postgres

import (
	"errors"
	"fmt"
	app "github.com/leartgjoni/go-rest-template"
	"math/rand"
	"strings"
	"time"

	_ "github.com/lib/pq"
)

// Ensure service implements interface.
var _ app.ArticleService = &ArticleService{}

// ArticleService represents a service to manage users.
type ArticleService struct {
	db *DB
}

// NewArticleService returns a new instance of ArticleService.
func NewArticleService(db *DB) *ArticleService {
	return &ArticleService{
		db: db,
	}
}

func (s *ArticleService) GetAll() ([]*app.Article, error) {
	rows, err := s.db.Query("SELECT * FROM articles")
	if err != nil {
		return []*app.Article{}, err
	}
	defer func() {
		if dErr := rows.Close(); dErr != nil && err == nil {
			err = dErr
		}
	}()

	var articles []*app.Article
	for rows.Next() {
		var article app.Article
		err := rows.Scan(&article.ID, &article.Slug, &article.Title, &article.Body, &article.UserId, &article.CreatedAt, &article.UpdatedAt)
		if err != nil {
			return []*app.Article{}, err
		}
		articles = append(articles, &article)
	}

	return articles, nil
}

func (s *ArticleService) GetBySlug(slug string) (*app.Article, error) {
	var article app.Article
	err := s.db.QueryRow("SELECT * FROM articles WHERE slug = $1", slug).Scan(&article.ID, &article.Slug, &article.Title, &article.Body, &article.UserId, &article.CreatedAt, &article.UpdatedAt)

	if err != nil {
		return &app.Article{}, err
	}

	if article.ID == 0 {
		return &app.Article{}, errors.New("not found")
	}

	return &article, nil
}
func (s *ArticleService) Save(a *app.Article) error {

	a.Slug = getSlug(a.Title, 12)
	row := s.db.QueryRow("INSERT INTO articles (slug, title, body, user_id, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id", a.Slug, a.Title, a.Body, a.UserId, a.CreatedAt, a.UpdatedAt)

	if err := row.Scan(&a.ID); err != nil {
		return err
	}

	if a.ID == 0 {
		return errors.New("unable to save")
	}

	return nil
}
func (s *ArticleService) Update(a *app.Article) error {
	err := s.db.QueryRow("UPDATE articles SET slug = $1, title = $2, body = $3, updated_at = $4 WHERE slug = $5 RETURNING id, slug, created_at", getSlug(a.Title, 12), a.Title, a.Body, a.UpdatedAt, a.Slug).Scan(&a.ID, &a.Slug, &a.CreatedAt)
	return err
}
func (s *ArticleService) Delete(slug string) error {
	_, err := s.db.Query("DELETE FROM articles WHERE slug LIKE $1", slug)
	return err
}

func getSlug(title string, length int) string {
	rand.Seed(time.Now().UnixNano())
	chars := []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789")
	var b strings.Builder
	for i := 0; i < length; i++ {
		b.WriteRune(chars[rand.Intn(len(chars))])
	}
	randomString := b.String()

	return fmt.Sprintf("%s-%s", strings.Replace(strings.ToLower(title), " ", "-", -1), randomString)
}
