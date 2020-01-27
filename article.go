package app

import "time"

type Article struct {
	ID        uint32    `json:"id"`
	Slug      string    `json:"slug"`
	Title     string    `json:"title"`
	Body      string    `json:"body"`
	UserId    uint32    `json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type ArticleService interface {
	GetAll() ([]*Article, error)
	GetBySlug(slug string) (*Article, error)
	Save(a *Article) error
	Update(a *Article) error
	Delete(slug string) error
}
