package main

import (
	"fmt"
	"encoding/json"
	"log"
	"net/http"
)

// Item x
type Item struct {
	ID int
	Done bool
	Title string
}

// Items x
type Items []Item

var items Items

func main() {
	http.HandleFunc("/", homeController)
	http.ListenAndServe(":13000", nil)
}

func homeController(w http.ResponseWriter, r *http.Request) {
	log.Println("request")
}

func listItemsAction(w http.ResponseWriter, r *http.Request) {
	body, err = json.Marshal(items)
	if err != nil {
		log.Fatal("Error enconding Items")
	}
	fmt.Printf(body, w)
	w.Header().Add("Content-Type", "application/json")
}