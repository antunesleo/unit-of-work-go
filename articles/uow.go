package articles

import (
	"context"

	"github.com/go-pg/pg/v10"
)

type UowStore interface {
	GetArticleRepository() ArticleRepository
}

type TransactionFunc func(UowStore) error

type Uow interface {
	WithinTx(context.Context, TransactionFunc) error
}

type goPgUowStore struct {
	articleRepository *goPgArticleRepository
}

func NewGoPgUowStore(articleRepository *goPgArticleRepository) *goPgUowStore {
	return &goPgUowStore{articleRepository: articleRepository}
}

func (a goPgUowStore) GetArticleRepository() ArticleRepository {
	return a.articleRepository
}

type goPgUow struct {
	db *pg.DB
}

func NewGoPgUow(db *pg.DB) *goPgUow {
	return &goPgUow{db: db}
}

func (u *goPgUow) WithinTx(ctx context.Context, fn TransactionFunc) error {
	err := u.db.RunInTransaction(ctx, func(tx *pg.Tx) error {
		articleRepository := goPgArticleRepository{tx}
		uowStore := &goPgUowStore{
			articleRepository: &articleRepository,
		}
		return fn(uowStore)
	})
	return err
}
