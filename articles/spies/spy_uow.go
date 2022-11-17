package spy_articles

import (
	"context"

	"github.com/antunesleo/rest-api-go/articles"
	mock_repositories "github.com/antunesleo/rest-api-go/articles/mocks"
)

type SpyUowStore struct {
	articleRepository  *mock_repositories.MockArticleRepository
	categoryRepository *mock_repositories.MockCategoryRepository
}

func (a SpyUowStore) GetArticleRepository() articles.ArticleRepository {
	return a.articleRepository
}

func (a SpyUowStore) GetCategoryRepository() articles.CategoryRepository {
	return a.categoryRepository
}

type WithinTxCall struct {
	Ctx context.Context
	Fn  articles.TransactionFunc
}

type SpyUow struct {
	uowStore      *SpyUowStore
	WithinTxCalls []WithinTxCall
}

func NewSpyUow(
	articleRepository *mock_repositories.MockArticleRepository,
	categoryRepository *mock_repositories.MockCategoryRepository,
) *SpyUow {
	uowStore := &SpyUowStore{articleRepository, categoryRepository}
	return &SpyUow{uowStore: uowStore}
}

func (u *SpyUow) WithinTx(ctx context.Context, fn articles.TransactionFunc) error {
	u.WithinTxCalls = append(u.WithinTxCalls, WithinTxCall{ctx, fn})
	return fn(u.uowStore)
}
