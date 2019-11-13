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

// ListItems -
func ListItems(w http.ResponseWriter, r *http.Request) {

	dbClient, _ := data.NewClient()
	items, _ := dbClient.List()

	respondJSON(w, items)
}

// AddItem -
func AddItem(w http.ResponseWriter, r *http.Request) {

	var newItem models.Item
	err := json.NewDecoder(r.Body).Decode(&newItem)
	if err != nil {
		log.Fatal("Cannot decode Item from JSON ", err)
	}
	defer r.Body.Close()

	dbClient, _ := data.NewClient()
	createdItem, _ := dbClient.Create(&newItem)

	respondJSON(w, createdItem)
}

// UpdateItem -
func UpdateItem(w http.ResponseWriter, r *http.Request) {

	id := getIDValue(r)

	dbClient, _ := data.NewClient()
	item, _ := dbClient.Get(id)

	if item == nil {
		respondJSON(w, nil)
		return
	}

	item.SetAsDone()

	item, _ = dbClient.Update(item)

	respondJSON(w, item)
}

// DeleteItem -
func DeleteItem(w http.ResponseWriter, r *http.Request) {

	id := getIDValue(r)
	dbClient, _ := data.NewClient()
	dbClient.Delete(id)

	w.WriteHeader(http.StatusOK)
}

// GetItem -
func GetItem(w http.ResponseWriter, r *http.Request) {

	id := getIDValue(r)
	dbClient, _ := data.NewClient()
	itemFound, _ := dbClient.Get(id)

	respondJSON(w, itemFound)
}

func getIDValue(r *http.Request) int {

	id, err := strconv.ParseInt(mux.Vars(r)["id"], 0, 0)
	if err != nil {
		log.Fatal("ID is not integer")
		return 0
	}

	return int(id)
}

func respondJSON(w http.ResponseWriter, data interface{}) {
	if data == nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(data)
}
