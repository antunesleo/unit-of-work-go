package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/antunesleo/rest-api-go/articles"
	"github.com/antunesleo/rest-api-go/core"
)

func startServer() {
	db := core.NewDB()
	uow := articles.NewGoPgUow(db)
	articleUserCases := articles.NewArticleUseCases(uow)
	articleHandlers := articles.NewArticleHandlers(articleUserCases)
	router := core.NewRouter()
	articles.AssignHandlers(router, articleHandlers)
	log.Fatal(http.ListenAndServe(":10000", router))
}

func main() {
	fmt.Println("Unit of work GO!")
	startServer()
}
