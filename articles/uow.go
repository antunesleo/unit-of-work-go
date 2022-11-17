package articles

import (
	"context"

	"github.com/go-pg/pg/v10"
)

type UowStore interface {
	GetArticleRepository() ArticleRepository
	GetCategoryRepository() CategoryRepository
}

type TransactionFunc func(UowStore) error

type Uow interface {
	WithinTx(context.Context, TransactionFunc) error
}

type goPgUowStore struct {
	articleRepository  *goPgArticleRepository
	categoryRepository *goPgCategoryRepository
}

func NewGoPgUowStore(
	articleRepository *goPgArticleRepository,
	categoryRepository *goPgCategoryRepository,
) *goPgUowStore {
	return &goPgUowStore{
		articleRepository:  articleRepository,
		categoryRepository: categoryRepository,
	}
}

func (a goPgUowStore) GetArticleRepository() ArticleRepository {
	return a.articleRepository
}

func (a goPgUowStore) GetCategoryRepository() CategoryRepository {
	return a.categoryRepository
}

type goPgUow struct {
	db *pg.DB
}

func NewGoPgUow(db *pg.DB) *goPgUow {
	return &goPgUow{db: db}
}

func (u *goPgUow) WithinTx(ctx context.Context, fn TransactionFunc) error {
	err := u.db.RunInTransaction(ctx, func(tx *pg.Tx) error {
		uowStore := NewGoPgUowStore(
			NewGoPgArticleRepository(tx),
			NewGoPgCategoryRepository(tx),
		)
		return fn(uowStore)
	})
	return err
}
