package memory

import (
	"log"

	"github.com/thiagolunardi/learning-golang/data/dberrors"
	"github.com/thiagolunardi/learning-golang/models"
)

// Repo -
type Repo struct {
	client string
}

var items models.Items

// NewClient -
func NewClient() (*Repo, error) {

	log.Println("Using In-Memory database")

	dataSeed()

	return &Repo{
		client: "Memory",
	}, nil
}

// List -
func (repo *Repo) List() (models.Items, error) {
	return items, nil
}

// Get -
func (repo *Repo) Get(ID int) (*models.Item, error) {

	for index := range items {
		if items[index].ID == ID {
			return &items[index], nil
		}
	}

	return nil, nil
}

// Create -
func (repo *Repo) Create(item *models.Item) (*models.Item, error) {

	id := 1

	for i := range items {
		if items[i].ID >= id {
			id = items[i].ID + 1
		}
	}

	item.ID = id

	items = append(items, *item)

	return item, nil
}

// Update -
func (repo *Repo) Update(item *models.Item) (*models.Item, error) {

	var existingItem *models.Item

	for index := range items {
		if items[index].ID == item.ID {
			existingItem = &items[index]
		}
	}

	if existingItem == nil {
		return item, dberrors.ErrItemNotFound
	}

	existingItem.Done = item.Done
	existingItem.Title = item.Title

	return existingItem, nil
}

// Delete -
func (repo *Repo) Delete(ID int) error {

	for index, item := range items {
		if item.ID == ID {
			items = append(items[:index], items[index+1:]...)
			break
		}
	}

	return nil
}

func dataSeed() {
	if items != nil {
		return
	}

	log.Println("Seeding data...")

	items = []models.Item{
		models.Item{
			ID:    1,
			Done:  false,
			Title: "Item A",
		},
		models.Item{
			ID:    2,
			Done:  false,
			Title: "Item B",
		},
	}
}
