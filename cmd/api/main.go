package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/antunesleo/rest-api-go/articles"
	"github.com/go-pg/pg/extra/pgdebug"
	"github.com/go-pg/pg/v10"
	"github.com/gorilla/mux"
)

func startServer() {
	db := pg.Connect(&pg.Options{
		User:     "postgres",
		Password: "rest-api-go",
		Database: "rest-api-go",
	})
	db.AddQueryHook(pgdebug.DebugHook{Verbose: true})

	uow := articles.NewGoPgUow(db)
	articleUserCases := articles.NewArticleUseCases(uow)
	articleHandlers := articles.NewArticleHandlers(articleUserCases)
	router := mux.NewRouter().StrictSlash(true)
	articles.BuildHandlers(router, articleHandlers)
	log.Fatal(http.ListenAndServe(":10000", router))
}

func main() {
	fmt.Println("Unity of work GO!")
	startServer()
}
