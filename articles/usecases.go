package articles

import (
	"context"
	"fmt"
)

type ArticleUseCasesExecutor interface {
	GetArticle(ctx context.Context, id string) (Article, error)
	GetArticles(ctx context.Context) ([]*Article, error)
	CreateArticle(ctx context.Context, id, title, content, desc string) (*Article, error)
	UpdateArticle(ctx context.Context, id, newTitle, newContent, newDesc string) (*Article, error)
	DeleteArticle(ctx context.Context, id string) error
}

type ArticleUseCases struct {
	articles []*Article
	uow      Uow
}

func NewArticleUseCases(articles []*Article, uow Uow) *ArticleUseCases {
	return &ArticleUseCases{articles: articles, uow: uow}
}

func (a ArticleUseCases) GetArticle(ctx context.Context, id string) (Article, error) {
	var (
		article Article
		err     error
	)

	withinTx := func(uows UowStore) error {
		article, err = uows.GetArticleRepository().FindById(id)
		return err
	}

	err = a.uow.WithinTx(ctx, withinTx)
	return article, err
}

func (a ArticleUseCases) CreateArticle(ctx context.Context, id, title, content, desc string) (*Article, error) {
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

func (a ArticleUseCases) GetArticles(ctx context.Context) ([]*Article, error) {
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

func (a ArticleUseCases) UpdateArticle(ctx context.Context, id, newTitle, newContent, newDesc string) (*Article, error) {
	var (
		article Article
		err     error
	)

	withinTx := func(uows UowStore) error {
		articleRepository := uows.GetArticleRepository()
		article, err = articleRepository.FindById(id)
		if err != nil {
			return err
		}
		fmt.Println("newTitle", newTitle)
		fmt.Println("article.Title", article.Title)
		article.Title = newTitle

		article.Content = newContent
		article.Desc = newDesc
		return articleRepository.Update(id, article)
	}

	err = a.uow.WithinTx(ctx, withinTx)
	return &article, err
}

func (a ArticleUseCases) DeleteArticle(ctx context.Context, id string) error {
	withinTx := func(uows UowStore) error {
		return uows.GetArticleRepository().Remove(id)
	}
	return a.uow.WithinTx(ctx, withinTx)
}
