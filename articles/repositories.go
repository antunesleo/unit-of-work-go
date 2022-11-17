package articles

import (
	"github.com/go-pg/pg/v10"
)

type ArticleRepository interface {
	FindAll() ([]*Article, error)
	Add(article *Article) error
}

type ArticleDB struct {
	tableName   struct{} `pg:"articles"` //nolint:unused
	Id          string
	Title       string
	Description string
	Content     string
}

type goPgArticleRepository struct {
	tx *pg.Tx
}

func NewGoPgArticleRepository(tx *pg.Tx) *goPgArticleRepository {
	return &goPgArticleRepository{tx: tx}
}

func (r *goPgArticleRepository) FindAll() ([]*Article, error) {
	articles := []*Article{}

	articleRows := []*ArticleDB{}
	err := r.tx.Model(&articleRows).Select()
	if err != nil {
		return articles, err
	}

	for _, articleRow := range articleRows {
		article := &Article{
			Id:      articleRow.Id,
			Content: articleRow.Content,
			Desc:    articleRow.Description,
			Title:   articleRow.Title,
		}
		articles = append(articles, article)
	}

	return articles, nil
}

func (r *goPgArticleRepository) Add(article *Article) error {
	articleRow := &ArticleDB{
		Id:          article.Id,
		Title:       article.Title,
		Description: article.Desc,
		Content:     article.Content,
	}
	_, err := r.tx.Model(articleRow).Insert()
	return err
}
