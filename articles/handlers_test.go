package articles_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/antunesleo/rest-api-go/articles"
	mock_usecases "github.com/antunesleo/rest-api-go/articles/mocks"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestArticlesHandlers(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	t.Run("should return status created when no errors", func(t *testing.T) {
		createArticleReq := articles.CreateArticleReq{
			Title:        "title",
			Content:      "content",
			Desc:         "desc",
			CategoryName: "category name",
		}
		body, _ := json.Marshal(createArticleReq)
		request := httptest.NewRequest(http.MethodPost, "/articles", bytes.NewBuffer(body))
		response := httptest.NewRecorder()

		articleId := int64(1)
		categoryId := int64(1)
		article := &articles.Article{
			Id:      articleId,
			Title:   createArticleReq.Title,
			Content: createArticleReq.Content,
			Desc:    createArticleReq.Desc,
			Category: articles.Category{
				Id:   categoryId,
				Name: createArticleReq.CategoryName,
			},
		}
		mockUseCase := mock_usecases.NewMockArticleUseCases(mockCtrl)
		mockUseCase.
			EXPECT().
			CreateArticle(
				request.Context(),
				createArticleReq.Title,
				createArticleReq.Content,
				createArticleReq.Desc,
				createArticleReq.CategoryName,
			).
			Return(article, nil)

		articleHandlers := articles.NewArticleHandlers(mockUseCase)

		articleHandlers.HandleCreateArticle(response, request)

		assert.Equal(t, http.StatusCreated, response.Result().StatusCode)
		var returnedArticle articles.Article
		err := json.NewDecoder(response.Body).Decode(&returnedArticle)

		assert.NoError(t, err)
		assert.Equal(t, articleId, returnedArticle.Id)
		assert.Equal(t, article.Title, returnedArticle.Title)
		assert.Equal(t, article.Category, returnedArticle.Category)
		assert.Equal(t, article.Content, returnedArticle.Content)
		assert.Equal(t, article.Category.Id, returnedArticle.Category.Id)
		assert.Equal(t, article.Category.Name, returnedArticle.Category.Name)
	})
}
