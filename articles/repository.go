package articles

import (
	"errors"
	"fmt"

	"github.com/go-pg/pg/v10"
)

type ArticleRepository interface {
	FindById(id string) (Article, error)
	FindAll() ([]*Article, error)
	Update(id string, updatedArticle Article) error
	Remove(id string) error
	Add(article *Article) error
}

type InMemoryArticleRepository struct {
	articles []*Article
}

func NewInMemoryArticleRepository(articles []*Article) *InMemoryArticleRepository {
	return &InMemoryArticleRepository{articles: articles}
}

var NotFoundError = errors.New("Not Found Error")

func (r *InMemoryArticleRepository) FindById(id string) (Article, error) {
	for _, article := range r.articles {
		if article.Id == id {
			return *article, nil
		}
	}
	return Article{}, NotFoundError
}

func (r *InMemoryArticleRepository) FindAll() ([]*Article, error) {
	if len(r.articles) == 0 {
		return []*Article{}, nil
	}
	return r.articles, nil
}

func (r *InMemoryArticleRepository) Remove(id string) error {
	for index, article := range r.articles {
		if article.Id == id {
			r.articles = append(r.articles[:index], r.articles[index+1:]...)
			return nil
		}
	}
	return NotFoundError
}

func (r *InMemoryArticleRepository) Add(article *Article) error {
	r.articles = append(r.articles, article)
	return nil
}

func (r *InMemoryArticleRepository) Update(id string, updateArticle Article) error {
	for index, article := range r.articles {
		if article.Id == id {
			r.articles[index].Content = updateArticle.Content
			r.articles[index].Desc = updateArticle.Desc
			r.articles[index].Title = updateArticle.Title
			return nil
		}
	}
	return NotFoundError
}

type ArticleRow struct {
	tableName   struct{} `pg:"articles"` //nolint:unused
	Id          string
	Title       string
	Description string
	Content     string
}

type PgGoArticleRepository struct {
	tx *pg.Tx
}

func NewPgGoArticleRepository(tx *pg.Tx) *PgGoArticleRepository {
	return &PgGoArticleRepository{tx: tx}
}

func (r *PgGoArticleRepository) FindById(id string) (Article, error) {
	articleRow := &ArticleRow{
		Id: id,
	}
	err := r.tx.Model(articleRow).
		WherePK().
		Select()

	if err != nil {
		fmt.Println("err", err)
		return Article{}, err
	}
	return Article{
		Id:      articleRow.Id,
		Content: articleRow.Content,
		Desc:    articleRow.Description,
		Title:   articleRow.Title,
	}, nil
}

func (r *PgGoArticleRepository) FindAll() ([]*Article, error) {
	return []*Article{}, nil
}

func (r *PgGoArticleRepository) Remove(id string) error {
	fmt.Println("Called Remove")
	return nil
}

func (r *PgGoArticleRepository) Add(article *Article) error {
	fmt.Println("Called Add")
	return nil
}

func (r *PgGoArticleRepository) Update(id string, updateArticle Article) error {
	fmt.Println("Called Update")
	return nil
}
