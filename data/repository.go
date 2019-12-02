package data

import (
	"os"
	"context"
	"github.com/thiagolunardi/learning-golang/data/mongodb"
	"github.com/thiagolunardi/learning-golang/data/memory"
	"github.com/thiagolunardi/learning-golang/models"
)

// TodoRepository -
type TodoRepository interface {
	List(ctx context.Context) (models.Items, error)
	Get(ctx context.Context, ID int) (*models.Item, error)
	Create(ctx context.Context, item *models.Item) (*models.Item, error)
	Update(ctx context.Context, item *models.Item) (*models.Item, error)
	Delete(ctx context.Context, ID int) error
}

// NewClient -
func NewClient() (TodoRepository, error) {	
	switch dbType := os.Getenv("DbType"); dbType {
	case "MongoDb": 
		return mongodb.NewClient()
	default:
		return memory.NewClient()
	}
}
