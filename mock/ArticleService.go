package mock

import app "github.com/leartgjoni/go-rest-template"

type ArticleService struct {
	GetAllFn      func() ([]*app.Article, error)
	GetAllInvoked bool

	GetBySlugFn      func(slug string) (*app.Article, error)
	GetBySlugInvoked bool

	SaveFn      func(a *app.Article) error
	SaveInvoked bool

	UpdateFn      func(a *app.Article) error
	UpdateInvoked bool

	DeleteFn      func(slug string) error
	DeleteInvoked bool
}

func (s *ArticleService) GetAll() ([]*app.Article, error) {
	s.GetAllInvoked = true
	return s.GetAllFn()
}

func (s *ArticleService) GetBySlug(slug string) (*app.Article, error) {
	s.GetBySlugInvoked = true
	return s.GetBySlugFn(slug)
}

func (s *ArticleService) Save(a *app.Article) error {
	s.SaveInvoked = true
	return s.SaveFn(a)
}

func (s *ArticleService) Update(a *app.Article) error {
	s.UpdateInvoked = true
	return s.UpdateFn(a)
}

func (s *ArticleService) Delete(slug string) error {
	s.DeleteInvoked = true
	return s.DeleteFn(slug)
}
