package articles

import (
	"context"
)

type ArticleUseCases interface {
	GetArticles(ctx context.Context) ([]*Article, error)
	CreateArticle(ctx context.Context, id, title, content, desc string) (*Article, error)
}

type articleUseCases struct {
	uow Uow
}

func NewArticleUseCases(uow Uow) *articleUseCases {
	return &articleUseCases{uow: uow}
}

func (a articleUseCases) CreateArticle(ctx context.Context, id, title, content, desc string) (*Article, error) {
	var article Article
	article.Id = id
	article.Title = title
	article.Content = content
	article.Desc = desc

	withinTx := func(uows UowStore) error {
		return uows.GetArticleRepository().Add(&article)
	}
	err := a.uow.WithinTx(ctx, withinTx)

	return &article, err
}

func (a articleUseCases) GetArticles(ctx context.Context) ([]*Article, error) {
	var (
		articles []*Article
		err      error
	)
	withinTx := func(uows UowStore) error {
		articles, err = uows.GetArticleRepository().FindAll()
		return err
	}
	err = a.uow.WithinTx(ctx, withinTx)
	return articles, err
}
