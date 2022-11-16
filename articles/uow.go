package articles

import (
	"context"
	"fmt"

	"github.com/go-pg/pg/v10"
)

type UowStore interface {
	GetArticleRepository() ArticleRepository
}

type InMemoryUowStore struct {
	articleRepository *InMemoryArticleRepository
}

func NewInMemoryUowStore(articleRepository *InMemoryArticleRepository) *InMemoryUowStore {
	return &InMemoryUowStore{articleRepository: articleRepository}
}

func (a InMemoryUowStore) GetArticleRepository() ArticleRepository {
	return a.articleRepository
}

type GoPgUowStore struct {
	articleRepository *PgGoArticleRepository
}

func NewGoPgUowStore(articleRepository *PgGoArticleRepository) *GoPgUowStore {
	return &GoPgUowStore{articleRepository: articleRepository}
}

func (a GoPgUowStore) GetArticleRepository() ArticleRepository {
	return a.articleRepository
}

type TransactionFunc func(UowStore) error

type Uow interface {
	WithinTx(context.Context, TransactionFunc) error
}

type InMemoryUow struct {
	articleRepository *InMemoryArticleRepository
}

func NewInMemoryUow(articleRepository *InMemoryArticleRepository) *InMemoryUow {
	return &InMemoryUow{articleRepository: articleRepository}
}

func (u *InMemoryUow) WithinTx(ctx context.Context, fn TransactionFunc) error {
	uowStore := &InMemoryUowStore{
		articleRepository: u.articleRepository,
	}
	fmt.Println("Started Tx")
	err := fn(uowStore)

	if err != nil {
		fmt.Println("Rolled back!")
	} else {
		fmt.Println("Commited!")
	}

	return err
}

type GoPgUow struct {
	db *pg.DB
}

func NewGoPgUow(db *pg.DB) *GoPgUow {
	return &GoPgUow{db: db}
}

func (u *GoPgUow) WithinTx(ctx context.Context, fn TransactionFunc) error {
	err := u.db.RunInTransaction(ctx, func(tx *pg.Tx) error {
		articleRepository := PgGoArticleRepository{tx}
		uowStore := &GoPgUowStore{
			articleRepository: &articleRepository,
		}
		return fn(uowStore)
	})
	return err
}
