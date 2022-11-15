package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/antunesleo/rest-api-go/articles"
	"github.com/gorilla/mux"
)

func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Welcome to the home page!")
	fmt.Println("Endpoint hit: homePage")
}

func startServer() {
	router := mux.NewRouter().StrictSlash(true)
	articles.BuildHandlers(router)
	log.Fatal(http.ListenAndServe(":10000", router))
}

func main() {
	fmt.Println("Rest API v2.0 - Mux Routers")
	startServer()
}
