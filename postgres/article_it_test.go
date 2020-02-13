package postgres

import (
	app "github.com/leartgjoni/go-rest-template"
	"regexp"
	"testing"
	"time"
)

func TestArticleServiceIntegration_GetAll(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	db := Suite.GetDb(t)
	Suite.CleanDb(t)

	userId := createUser(db, t)

	timeNow := time.Now().Truncate(time.Millisecond)
	article1 := app.Article{
		ID:        0,
		Slug:      "slug-1",
		Title:     "title one",
		Body:      "body one",
		UserId:    userId,
		CreatedAt: timeNow,
		UpdatedAt: timeNow,
	}

	article2 := app.Article{
		ID:        0,
		Slug:      "slug-2",
		Title:     "title two",
		Body:      "body two",
		UserId:    userId,
		CreatedAt: time.Time{},
		UpdatedAt: time.Time{},
	}

	articles := []app.Article{article1, article2}

	for _, a := range articles {
		if _, err := db.Exec("INSERT INTO articles (slug, title, body, user_id, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id", a.Slug, a.Title, a.Body, a.UserId, a.CreatedAt, a.UpdatedAt); err != nil {
			t.Fatal("cannot insert article", err)
		}
	}

	as := NewArticleService(db)
	dbArticles, err := as.GetAll()

	if err != nil {
		t.Fatal("err on GetAll", err)
	}

	if len(dbArticles) != 2 {
		t.Fatalf("wrong length. %v instead of 2", len(dbArticles))
	}

	for i, a := range articles {
		dbArticle := dbArticles[i]
		if a.Slug != dbArticle.Slug || a.Title != dbArticle.Title || a.Body != dbArticle.Body || a.UserId != dbArticle.UserId || !a.CreatedAt.Equal(dbArticle.CreatedAt) || !a.UpdatedAt.Equal(dbArticle.UpdatedAt) {
			t.Errorf("Expected %v but got %v", a, dbArticle)
		}
	}

}

func TestArticleServiceIntegration_GetBySlug(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	db := Suite.GetDb(t)
	Suite.CleanDb(t)

	userId := createUser(db, t)

	timeNow := time.Now().Truncate(time.Millisecond)
	// expected article
	eArticle := app.Article{
		ID:        0,
		Slug:      "random-slug",
		Title:     "random title",
		Body:      "random body",
		UserId:    userId,
		CreatedAt: timeNow,
		UpdatedAt: timeNow,
	}

	row := db.QueryRow("INSERT INTO articles (slug, title, body, user_id, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id", eArticle.Slug, eArticle.Title, eArticle.Body, eArticle.UserId, eArticle.CreatedAt, eArticle.UpdatedAt)
	if err := row.Scan(&eArticle.ID); err != nil {
		t.Fatal("error while inserting article", err)
	}

	as := NewArticleService(db)
	// actual article
	aArticle, err := as.GetBySlug(eArticle.Slug)
	if err != nil {
		t.Fatal("error getting the article", err)
	}

	if aArticle.Slug != eArticle.Slug || aArticle.Title != eArticle.Title || aArticle.Body != eArticle.Body || aArticle.UserId != eArticle.UserId || !aArticle.CreatedAt.Equal(eArticle.CreatedAt) || !aArticle.UpdatedAt.Equal(eArticle.UpdatedAt) {
		t.Errorf("Expected %v but got %v", eArticle, aArticle)
	}
}

func TestArticleServiceIntegration_Save(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	db := Suite.GetDb(t)
	Suite.CleanDb(t)

	userId := createUser(db, t)

	timeNow := time.Now().Truncate(time.Millisecond)
	// expected article
	article := app.Article{
		ID:        0,
		Slug:      "",
		Title:     "random title",
		Body:      "random body",
		UserId:    userId,
		CreatedAt: timeNow,
		UpdatedAt: timeNow,
	}

	as := NewArticleService(db)
	if err := as.Save(&article); err != nil {
		t.Fatal("cannot save article", err)
	}

	if article.ID == 0 {
		t.Fatal("article id still zero")
	}

	var dbArticle app.Article
	if err := db.QueryRow("SELECT * FROM articles WHERE id = $1", article.ID).Scan(&dbArticle.ID, &dbArticle.Slug, &dbArticle.Title, &dbArticle.Body, &dbArticle.UserId, &dbArticle.CreatedAt, &dbArticle.UpdatedAt); err != nil {
		t.Fatal("cannot read article from db", err)
	}

	if !regexp.MustCompile(`random-title-[a-zA-Z0-9]{12}`).Match([]byte(dbArticle.Slug)) {
		t.Fatal("slug format is wrong", dbArticle.Slug)
	}

	if article.ID != dbArticle.ID || article.Title != dbArticle.Title || article.Body != dbArticle.Body || article.UserId != dbArticle.UserId || !article.CreatedAt.Equal(dbArticle.CreatedAt) || !article.UpdatedAt.Equal(dbArticle.UpdatedAt) {
		t.Fatalf("Expected %v but got %v", article, dbArticle)
	}
}

func TestArticleServiceIntegration_Update(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	db := Suite.GetDb(t)
	Suite.CleanDb(t)

	userId := createUser(db, t)

	timeNow := time.Now().Truncate(time.Millisecond)
	// expected article
	article := app.Article{
		ID:        0,
		Slug:      "",
		Title:     "random title",
		Body:      "random body",
		UserId:    userId,
		CreatedAt: timeNow,
		UpdatedAt: timeNow,
	}

	row := db.QueryRow("INSERT INTO articles (slug, title, body, user_id, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id", article.Slug, article.Title, article.Body, article.UserId, article.CreatedAt, article.UpdatedAt)
	if err := row.Scan(&article.ID); err != nil {
		t.Fatal("error while inserting article", err)
	}

	as := NewArticleService(db)

	article.Title = "random title updated"
	article.Body = "random body updated"

	if err := as.Update(&article); err != nil {
		t.Fatal("cannot update article", err)
	}

	var dbArticle app.Article
	if err := db.QueryRow("SELECT * FROM articles WHERE id = $1", article.ID).Scan(&dbArticle.ID, &dbArticle.Slug, &dbArticle.Title, &dbArticle.Body, &dbArticle.UserId, &dbArticle.CreatedAt, &dbArticle.UpdatedAt); err != nil {
		t.Fatal("cannot read article from db", err)
	}

	if article.ID != dbArticle.ID || dbArticle.Slug[:len(dbArticle.Slug)-13] != "random-title-updated" || article.Title != dbArticle.Title || article.Body != dbArticle.Body || article.UserId != dbArticle.UserId || !article.CreatedAt.Equal(dbArticle.CreatedAt) || !article.UpdatedAt.Equal(dbArticle.UpdatedAt) {
		t.Fatalf("Expected %v but got %v", article, dbArticle)
	}
}

func TestArticleServiceIntegration_Delete(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	db := Suite.GetDb(t)
	Suite.CleanDb(t)

	userId := createUser(db, t)

	timeNow := time.Now()
	// expected article
	article := app.Article{
		ID:        0,
		Slug:      "random-title-slug",
		Title:     "random title",
		Body:      "random body",
		UserId:    userId,
		CreatedAt: timeNow,
		UpdatedAt: timeNow,
	}

	row := db.QueryRow("INSERT INTO articles (slug, title, body, user_id, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id", article.Slug, article.Title, article.Body, article.UserId, article.CreatedAt, article.UpdatedAt)
	if err := row.Scan(&article.ID); err != nil {
		t.Fatal("error while inserting article", err)
	}

	as := NewArticleService(db)

	if err := as.Delete(article.Slug); err != nil {
		t.Fatal("cannot delete article", err)
	}

	count := 0
	err := db.QueryRow("SELECT COUNT(id) FROM articles WHERE slug = $1", article.Slug).Scan(&count)
	if err != nil {
		t.Fatal("could not count articles", err)
	}

	if count != 0 {
		t.Fatal("article not deleted", count)
	}
}

func createUser(db *DB, t *testing.T) uint32 {
	var userId uint32 = 0
	row := db.QueryRow("INSERT INTO users (username, email, password, created_at, updated_at) VALUES ($1, $2, $3, $4, $5) RETURNING id", "random username", "test@test.com", "random-password", time.Now(), time.Now())
	if err := row.Scan(&userId); err != nil {
		t.Fatal("error while inserting user", err)
	}

	return userId
}
