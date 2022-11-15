package articles

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-pg/pg/v10"
	"github.com/gorilla/mux"
)

var Articles []Article

type UpdateArticleReq struct {
	Title   string `json:"Title"`
	Desc    string `json:"desc"`
	Content string `json:"content"`
}

type CreateArticleReq struct {
	Id      string `json:"Id"`
	Title   string `json:"Title"`
	Desc    string `json:"desc"`
	Content string `json:"content"`
}

func handleRoot(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Welcome to the home page!")
	fmt.Println("Endpoint hit: homePage")
}

func BuildHandlers(router *mux.Router) {
	db := pg.Connect(&pg.Options{
		User:     "postgres",
		Password: "rest-api-go",
		Database: "rest-api-go",
	})

	Articles := []*Article{}
	Articles = append(Articles, &Article{Id: "1", Title: "Hello", Desc: "Article Description", Content: "Article Content"})
	Articles = append(Articles, &Article{Id: "2", Title: "Hello 2", Desc: "Article Description", Content: "Article Content"})

	// articleRepository := &InMemoryArticleRepository{Articles}

	uow := NewPgGoUow(db)
	articleUserCases := NewArticleUseCases(Articles, uow)
	articleHandlers := ArticleHandlers{articleUseCases: articleUserCases}
	router.HandleFunc("/", handleRoot)
	router.HandleFunc("/articles", articleHandlers.handleGetArticles).Methods("GET")
	router.HandleFunc("/articles", articleHandlers.handleCreateArticle).Methods("POST")
	router.HandleFunc("/articles/{id}", articleHandlers.handleGetArticle).Methods("GET")
	router.HandleFunc("/articles/{id}", articleHandlers.handleDeleteArticle).Methods("DELETE")
	router.HandleFunc("/articles/{id}", articleHandlers.handleUpdateArticle).Methods("PUT")
}

type ArticleHandlers struct {
	articleUseCases ArticleUseCasesExecutor
}

func (h *ArticleHandlers) handleGetArticle(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	vars := mux.Vars(r)

	article, err := h.articleUseCases.GetArticle(ctx, vars["id"])

	if err != nil {
		if err == NotFoundError {
			http.Error(w, "Not Found", http.StatusNotFound)
		} else {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
		return
	}

	json.NewEncoder(w).Encode(article)
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
		createArticleReq.Id,
		createArticleReq.Title,
		createArticleReq.Content,
		createArticleReq.Desc,
	)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(article)
}

func (h *ArticleHandlers) handleDeleteArticle(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	vars := mux.Vars(r)

	err := h.articleUseCases.DeleteArticle(ctx, vars["id"])
	if err != nil {
		if err == NotFoundError {
			http.Error(w, "Not Found", http.StatusNotFound)
		} else {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *ArticleHandlers) handleUpdateArticle(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	vars := mux.Vars(r)

	var updateArticleReq UpdateArticleReq
	json.NewDecoder(r.Body).Decode(&updateArticleReq)

	article, err := h.articleUseCases.UpdateArticle(
		ctx,
		vars["id"],
		updateArticleReq.Title,
		updateArticleReq.Content,
		updateArticleReq.Desc,
	)
	if err != nil {
		if err == NotFoundError {
			http.Error(w, "Not Found", http.StatusNotFound)
		} else {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
	} else {
		json.NewEncoder(w).Encode(article)
	}
}
