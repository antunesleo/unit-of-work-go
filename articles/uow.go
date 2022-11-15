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

type PgGoUowStore struct {
	articleRepository *PgGoArticleRepository
}

func NewPgGoUowStore(articleRepository *PgGoArticleRepository) *PgGoUowStore {
	return &PgGoUowStore{articleRepository: articleRepository}
}

func (a PgGoUowStore) GetArticleRepository() ArticleRepository {
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

type PgGoUow struct {
	db *pg.DB
}

func NewPgGoUow(db *pg.DB) *PgGoUow {
	return &PgGoUow{db: db}
}

func (u *PgGoUow) WithinTx(ctx context.Context, fn TransactionFunc) error {
	tx, err := u.db.Begin()
	fmt.Println("started")
	if err != nil {
		fmt.Println("err1", err)
		return err
	}
	defer tx.Close()

	articleRepository := PgGoArticleRepository{tx}
	uowStore := &PgGoUowStore{
		articleRepository: &articleRepository,
	}

	if err := fn(uowStore); err != nil {
		fmt.Println("err", err)
		fmt.Println("rolled back")
		_ = tx.Rollback()
		return err
	}

	if err := tx.Commit(); err != nil {
		fmt.Println("err", err)
		return err
	}
	fmt.Println("commited")

	return nil
}
