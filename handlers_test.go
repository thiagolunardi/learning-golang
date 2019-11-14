package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/matryer/is"
	"github.com/thiagolunardi/learning-golang/models"
)

func TestListItems(t *testing.T) {
	req, err := http.NewRequest(http.MethodGet, "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	router := NewRouter()
	router.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v instead of %v",
			status, http.StatusOK)
	}

	body := rr.Body.Bytes()
	var items models.Items
	err = json.Unmarshal(body, &items)
	if err != nil {
		t.Errorf("handler returned unexpected body: got %v", rr.Body.String())
	}

	if len(items) != 2 {
		t.Errorf("handler returned unexpected items: got %v instead of %v",
			len(items), 2)
	}
}

func TestGetItem(t *testing.T) {
	req, err := http.NewRequest(http.MethodGet, "/1", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	router := NewRouter()
	router.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v instead of %v",
			status, http.StatusOK)
	}

	body := rr.Body.Bytes()
	var item models.Item
	err = json.Unmarshal(body, &item)
	if err != nil {
		t.Errorf("handler returned unexpected body: got %v", rr.Body.String())
	}

	is := is.New(t)

	is.Equal(item.ID, 1)
	is.Equal(item.Done, false)
	is.Equal(item.Title, "Item A")
}

func TestUpdateItem(t *testing.T) {
	req, err := http.NewRequest(http.MethodPut, "/1", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	router := NewRouter()
	router.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v instead of %v",
			status, http.StatusOK)
	}

	body := rr.Body.Bytes()
	var item models.Item
	err = json.Unmarshal(body, &item)
	if err != nil {
		t.Errorf("handler returned unexpected body: got %v", rr.Body.String())
	}

	is := is.New(t)

	is.Equal(item.ID, 1)
	is.Equal(item.Done, true)
	is.Equal(item.Title, "Item A")
}

func TestDeleteItem(t *testing.T) {
	req, err := http.NewRequest(http.MethodDelete, "/1", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	router := NewRouter()
	router.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v instead of %v",
			status, http.StatusOK)
	}

	req, err = http.NewRequest(http.MethodGet, "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr = httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	var items models.Items
	err = json.Unmarshal(rr.Body.Bytes(), &items)
	if err != nil {
		t.Errorf("handler returned unexpected body: got %v", rr.Body.String())
	}

	is := is.New(t)

	is.Equal(len(items), 1)
	is.Equal(items[0].ID, 2)
	is.Equal(items[0].Done, false)
	is.Equal(items[0].Title, "Item B")
}

func TestAddItem(t *testing.T) {
	newItem := models.Item{
		Title: "Item X",
	}
	buffer, _ := json.Marshal(newItem)

	req, err := http.NewRequest(http.MethodPost, "/", bytes.NewReader(buffer))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	router := NewRouter()
	router.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v instead of %v",
			status, http.StatusOK)
	}

	req, err = http.NewRequest(http.MethodGet, "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr = httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	var items models.Items
	err = json.Unmarshal(rr.Body.Bytes(), &items)
	if err != nil {
		t.Errorf("handler returned unexpected body: got %v", rr.Body.String())
	}

	is := is.New(t)

	is.Equal(len(items), 2)
	is.Equal(items[1].ID, 3)
	is.Equal(items[1].Done, false)
	is.Equal(items[1].Title, "Item X")
}
