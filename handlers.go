package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/thiagolunardi/learning-golang/data"
	"github.com/thiagolunardi/learning-golang/models"

	"github.com/gorilla/mux"
)

var dbClient data.TodoRepository

// Initialize - 
func Initialize() {
	if dbClient != nil { return }

	var err error
	dbClient, err = data.NewClient()
	if err != nil { log.Fatal(err) }
}

// ListItems -
func ListItems(w http.ResponseWriter, r *http.Request) {
	items, _ := dbClient.List(r.Context())

	respondOK(w, items)
}

// AddItem -
func AddItem(w http.ResponseWriter, r *http.Request) {

	var newItem models.Item
	err := json.NewDecoder(r.Body).Decode(&newItem)
	if err != nil {
		log.Fatal("Cannot decode Item from JSON ", err)
	}
	defer r.Body.Close()

	createdItem, _ := dbClient.Create(r.Context(), &newItem)

	respondOK(w, createdItem)
}

// UpdateItem -
func UpdateItem(w http.ResponseWriter, r *http.Request) {

	id := getIDValue(r)

	dbClient, _ := data.NewClient()
	item, _ := dbClient.Get(r.Context(), id)

	if item == nil {
		respondNotFound(w)
		return
	}

	item.SetAsDone()

	item, _ = dbClient.Update(r.Context(), item)

	respondOK(w, item)
}

// DeleteItem -
func DeleteItem(w http.ResponseWriter, r *http.Request) {

	id := getIDValue(r)
	dbClient.Delete(r.Context(), id)

	w.WriteHeader(http.StatusOK)
}

// GetItem -
func GetItem(w http.ResponseWriter, r *http.Request) {

	id := getIDValue(r)
	itemFound, _ := dbClient.Get(r.Context(), id)

	if itemFound != nil {
		respondOK(w, itemFound)
	} else {
		respondNotFound(w)
	}
}

func getIDValue(r *http.Request) int {

	id, err := strconv.ParseInt(mux.Vars(r)["id"], 0, 0)
	if err != nil {
		log.Fatal("ID is not integer")
		return 0
	}

	return int(id)
}

func respondOK(w http.ResponseWriter, data interface{}) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(data)
}

func respondNotFound(w http.ResponseWriter) {
	w.WriteHeader(http.StatusNotFound)
}
