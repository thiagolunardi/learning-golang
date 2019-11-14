package main

import (
	"net/http"

	"github.com/gorilla/mux"
)

// NewRouter -
func NewRouter() *mux.Router {
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", ListItems).Methods(http.MethodGet)
	router.HandleFunc("/", AddItem).Methods(http.MethodPost)
	router.HandleFunc("/{id}", GetItem).Methods(http.MethodGet)
	router.HandleFunc("/{id}", UpdateItem).Methods(http.MethodPut)
	router.HandleFunc("/{id}", DeleteItem).Methods(http.MethodDelete)

	return router
}