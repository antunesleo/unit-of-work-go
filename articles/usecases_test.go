package articles_test

import (
	"context"
	"errors"
	"testing"

	"github.com/antunesleo/rest-api-go/articles"
	mock_uow "github.com/antunesleo/rest-api-go/articles/manual_mocks"
	mock_articles "github.com/antunesleo/rest-api-go/articles/mocks"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestUseCaseArticle(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	ctx := context.TODO()

	t.Run("should create article and category", func(t *testing.T) {
		category := &articles.Category{
			Name: "Category",
		}
		article := &articles.Article{
			Title:    "title",
			Content:  "Content",
			Desc:     "Desc",
			Category: *category,
		}

		categoryRepository := mock_articles.NewMockCategoryRepository(mockCtrl)
		categoryRepository.EXPECT().FindByName(category.Name).Return(articles.Category{}, articles.NotFoundError)
		categoryRepository.EXPECT().Add(category).Return(nil)
		articleRepository := mock_articles.NewMockArticleRepository(mockCtrl)
		articleRepository.EXPECT().Add(article).Return(nil)

		uow := mock_uow.NewSpyUow(articleRepository, categoryRepository)
		articleUseCases := articles.NewArticleUseCases(uow)

		actualArticle, err := articleUseCases.CreateArticle(
			ctx,
			article.Title,
			article.Content,
			article.Desc,
			category.Name,
		)

		assert.Equal(t, uow.WithinTxCalls[0].Ctx, ctx)
		assert.Equal(t, actualArticle, article)
		assert.NoError(t, err)
	})

	t.Run("should create article and associate category", func(t *testing.T) {
		category := articles.Category{
			Id:   1,
			Name: "Category",
		}
		article := &articles.Article{
			Title:    "title",
			Content:  "Content",
			Desc:     "Desc",
			Category: category,
		}

		categoryRepository := mock_articles.NewMockCategoryRepository(mockCtrl)
		categoryRepository.EXPECT().FindByName(category.Name).Return(category, nil)
		articleRepository := mock_articles.NewMockArticleRepository(mockCtrl)
		articleRepository.EXPECT().Add(article).Return(nil)

		uow := mock_uow.NewSpyUow(articleRepository, categoryRepository)
		articleUseCases := articles.NewArticleUseCases(uow)

		actualArticle, err := articleUseCases.CreateArticle(
			ctx,
			article.Title,
			article.Content,
			article.Desc,
			category.Name,
		)

		assert.Equal(t, uow.WithinTxCalls[0].Ctx, ctx)
		assert.Equal(t, actualArticle, article)
		assert.NoError(t, err)
	})

	t.Run("should fail to create article", func(t *testing.T) {
		category := articles.Category{
			Id:   1,
			Name: "Category",
		}
		article := &articles.Article{
			Title:    "title",
			Content:  "Content",
			Desc:     "Desc",
			Category: category,
		}

		errorToCreate := errors.New("Error to create")
		categoryRepository := mock_articles.NewMockCategoryRepository(mockCtrl)
		categoryRepository.EXPECT().FindByName(category.Name).Return(category, nil)
		articleRepository := mock_articles.NewMockArticleRepository(mockCtrl)
		articleRepository.EXPECT().Add(article).Return(errorToCreate)

		uow := mock_uow.NewSpyUow(articleRepository, categoryRepository)
		articleUseCases := articles.NewArticleUseCases(uow)

		_, err := articleUseCases.CreateArticle(
			ctx,
			article.Title,
			article.Content,
			article.Desc,
			category.Name,
		)

		assert.Equal(t, uow.WithinTxCalls[0].Ctx, ctx)
		assert.Error(t, err)
		assert.Equal(t, err, errorToCreate)
	})

	t.Run("should fail to create category", func(t *testing.T) {
		category := &articles.Category{
			Name: "Category",
		}
		article := &articles.Article{
			Title:    "title",
			Content:  "Content",
			Desc:     "Desc",
			Category: *category,
		}

		errorToCreateCategory := errors.New("Create Error")
		categoryRepository := mock_articles.NewMockCategoryRepository(mockCtrl)
		categoryRepository.EXPECT().FindByName(category.Name).Return(articles.Category{}, articles.NotFoundError)
		categoryRepository.EXPECT().Add(category).Return(errorToCreateCategory)
		articleRepository := mock_articles.NewMockArticleRepository(mockCtrl)

		uow := mock_uow.NewSpyUow(articleRepository, categoryRepository)
		articleUseCases := articles.NewArticleUseCases(uow)

		_, err := articleUseCases.CreateArticle(
			ctx,
			article.Title,
			article.Content,
			article.Desc,
			category.Name,
		)

		assert.Equal(t, uow.WithinTxCalls[0].Ctx, ctx)
		assert.Error(t, err)
		assert.Equal(t, err, errorToCreateCategory)
	})

	t.Run("should fail to get category", func(t *testing.T) {
		category := &articles.Category{
			Name: "Category",
		}
		article := &articles.Article{
			Title:    "title",
			Content:  "Content",
			Desc:     "Desc",
			Category: *category,
		}

		errorToGetCategory := errors.New("Get Error")
		categoryRepository := mock_articles.NewMockCategoryRepository(mockCtrl)
		categoryRepository.EXPECT().FindByName(category.Name).Return(articles.Category{}, errorToGetCategory)
		articleRepository := mock_articles.NewMockArticleRepository(mockCtrl)

		uow := mock_uow.NewSpyUow(articleRepository, categoryRepository)
		articleUseCases := articles.NewArticleUseCases(uow)

		_, err := articleUseCases.CreateArticle(
			ctx,
			article.Title,
			article.Content,
			article.Desc,
			category.Name,
		)

		assert.Equal(t, uow.WithinTxCalls[0].Ctx, ctx)
		assert.Error(t, err)
		assert.Equal(t, err, errorToGetCategory)
	})
}
