package data

import (
	"github.com/thiagolunardi/learning-golang/data/memory"
	"github.com/thiagolunardi/learning-golang/models"
)

// TodoRepository -
type TodoRepository interface {
	List() (models.Items, error)
	Get(ID int) (*models.Item, error)
	Create(item *models.Item) (*models.Item, error)
	Update(item *models.Item) (*models.Item, error)
	Delete(ID int) error
}

// NewClient -
func NewClient() (TodoRepository, error) {
	return memory.NewClient()
}
