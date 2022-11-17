package articles

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

var Articles []Article

type UpdateArticleReq struct {
	Title   string `json:"Title"`
	Desc    string `json:"desc"`
	Content string `json:"content"`
}

type CreateArticleReq struct {
	Title        string `json:"title"`
	Desc         string `json:"desc"`
	Content      string `json:"content"`
	CategoryName string `json:"categoryName"`
}

func AssignHandlers(router *mux.Router, articleHandlers *ArticleHandlers) {
	router.HandleFunc("/articles", articleHandlers.handleGetArticles).Methods("GET")
	router.HandleFunc("/articles", articleHandlers.handleCreateArticle).Methods("POST")
}

type ArticleHandlers struct {
	articleUseCases ArticleUseCases
}

func NewArticleHandlers(useCases ArticleUseCases) *ArticleHandlers {
	return &ArticleHandlers{articleUseCases: useCases}
}

func (h *ArticleHandlers) handleGetArticles(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	articles, err := h.articleUseCases.GetArticles(ctx)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(articles)
}

func (h *ArticleHandlers) handleCreateArticle(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var createArticleReq CreateArticleReq
	json.NewDecoder(r.Body).Decode(&createArticleReq)

	article, err := h.articleUseCases.CreateArticle(
		ctx,
		createArticleReq.Title,
		createArticleReq.Content,
		createArticleReq.Desc,
		createArticleReq.CategoryName,
	)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(article)
}
