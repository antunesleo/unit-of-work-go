package articles

import (
	"errors"

	"github.com/go-pg/pg/v10"
)

var NotFoundError = errors.New("Not Found Error")

type ArticleRepository interface {
	FindAll() ([]*Article, error)
	Add(article *Article) error
}

type ArticleDB struct {
	tableName   struct{} `pg:"articles"` //nolint:unused
	Id          int64
	Title       string
	Description string
	Content     string
	CategoryId  int64
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
	articleDB := &ArticleDB{
		Title:       article.Title,
		Description: article.Desc,
		Content:     article.Content,
		CategoryId:  article.Category.Id,
	}
	_, err := r.tx.Model(articleDB).Insert()
	article.Id = articleDB.Id
	return err
}

type CategoryRepository interface {
	FindByName(string) (Category, error)
	Add(*Category) error
}

type CategoryDB struct {
	tableName struct{} `pg:"categories"` //nolint:unused
	Id        int64
	Name      string
}

type goPgCategoryRepository struct {
	tx *pg.Tx
}

func NewGoPgCategoryRepository(tx *pg.Tx) *goPgCategoryRepository {
	return &goPgCategoryRepository{tx: tx}
}

func (r *goPgCategoryRepository) FindByName(name string) (Category, error) {
	categoryDB := &CategoryDB{}
	err := r.tx.Model(categoryDB).
		Where("name = ?", name).
		Select()

	if err != nil {
		if err == pg.ErrNoRows {
			return Category{}, NotFoundError
		}
		return Category{}, err
	}
	return Category{
		Id:   categoryDB.Id,
		Name: categoryDB.Name,
	}, nil
}

func (r *goPgCategoryRepository) Add(category *Category) error {
	cdb := &CategoryDB{
		Name: category.Name,
	}
	_, err := r.tx.Model(cdb).Insert()
	category.Id = cdb.Id
	return err
}
