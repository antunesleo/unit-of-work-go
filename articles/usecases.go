package articles

import (
	"context"
)

type ArticleUseCases interface {
	GetArticles(ctx context.Context) ([]*Article, error)
	CreateArticle(ctx context.Context, title, content, desc, categoryName string) (*Article, error)
}

type articleUseCases struct {
	uow Uow
}

func NewArticleUseCases(uow Uow) *articleUseCases {
	return &articleUseCases{uow: uow}
}

func (a articleUseCases) CreateArticle(
	ctx context.Context,
	title, content, desc, categoryName string,
) (*Article, error) {
	article := Article{
		Title:   title,
		Content: content,
		Desc:    desc,
	}

	withinTx := func(uows UowStore) error {
		categoryRepository := uows.GetCategoryRepository()
		category, err := categoryRepository.FindByName(categoryName)

		if err != nil {
			if err != NotFoundError {
				return err
			}

			category = Category{
				Name: categoryName,
			}
			addErr := categoryRepository.Add(&category)
			if addErr != nil {
				return addErr
			}
		}

		article.Category = category
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
