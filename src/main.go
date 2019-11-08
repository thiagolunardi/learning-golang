package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"./models"
	"strconv"

	"github.com/gorilla/mux"
)

var items = models.Items{}

func main() {
	httpPort := 13000

	seedData()

	log.Println("Server available at following address:")
	log.Printf("    http://localhost:%d/", httpPort)

  handler := requestHandler()
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", httpPort), logRequest(handler)))

	log.Println("Terminated.")
}

func requestHandler () http.Handler {
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", listItemsAction).Methods(http.MethodGet)
	router.HandleFunc("/", addItemAction).Methods(http.MethodPost)
	router.HandleFunc("/{id}", getItemAction).Methods(http.MethodGet)
	router.HandleFunc("/{id}", updateItemAction).Methods(http.MethodPut)
	router.HandleFunc("/{id}", deleteItemAction).Methods(http.MethodDelete)

	return router
}

func listItemsAction(w http.ResponseWriter, r *http.Request) {
	log.Println("listItemsAction")
	respondJSON(w, items)
}

func addItemAction(w http.ResponseWriter, r *http.Request) {
	log.Println("addItemAction")
	var item models.Item
	err := json.NewDecoder(r.Body).Decode(&item)
	if err != nil {
		log.Fatal("Cannot decode Item from JSON ", err)
	}
	defer r.Body.Close()

	item.ID = len(items) + 1
	items = append(items, item)

	respondJSON(w, item)
}

func updateItemAction(w http.ResponseWriter, r *http.Request) {
	log.Println("updateItemAction")
	id := getIDValue(r)
	
	var item *models.Item
	item = findByID(id)
	item.Done = true

	log.Println(items)

	respondJSON(w, item);
}

func deleteItemAction(w http.ResponseWriter, r *http.Request) {
	log.Println("deleteItemAction")
	id := getIDValue(r)
	removeByID(id)
	w.WriteHeader(http.StatusOK)
}

func getItemAction(w http.ResponseWriter, r *http.Request) {
	log.Println("getItemAction")
	id := getIDValue(r)

	itemFound := findByID(int(id))

	respondJSON(w, itemFound)
}

func replaceItem(item models.Item) {
	for index, item := range items {
        if item.ID == item.ID {
			items[index] = item
            break
        }
	}
}


func removeByID(id int) {
	for index, item := range items {
        if item.ID == id {
            items = append(items[:index], items[index+1:]...)
            break
        }
	}
}

func getIDValue(r *http.Request) int {
	id, err := strconv.ParseInt( mux.Vars(r)["id"], 0, 0 )
	if err != nil {
		log.Fatal("ID is not integer")		
		return 0
	}
	return int(id)
}

func findByID(id int) *models.Item {
  for _, item := range items {
    if item.ID == id {
      return &item
    }
  }
  return nil
}

func respondJSON(w http.ResponseWriter, data interface{}) {	
	var defaultValue *models.Item
	if data == defaultValue {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(data)
}

func seedData() {
	log.Println("Seeding data sample...")
	items = models.DataSeed()
}

func logRequest(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s %s\n", r.RemoteAddr, r.Method, r.URL)
		handler.ServeHTTP(w, r)
	})
}
