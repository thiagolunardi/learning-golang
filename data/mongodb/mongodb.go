package mongodb

import (
	"log"
	"go.mongodb.org/mongo-driver/bson"
	"context"
	"time"
	"github.com/thiagolunardi/learning-golang/data/dberrors"
	"github.com/thiagolunardi/learning-golang/models"
	"go.mongodb.org/mongo-driver/mongo/readpref"
    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
)

// Repo -
type Repo struct {
	client string
}

var dbClient *mongo.Client

var dbContext context.Context
var cancelFunc context.CancelFunc

var itemsCollection *mongo.Collection

// NewClient -
func NewClient() (*Repo, error) {

	dbClient, err := mongo.NewClient(options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil { log.Fatal(err) }

	dbClient.Ping(context.Background(), readpref.Primary())
	log.Println("Using MongoDb database")

	itemsCollection = dbClient.Database("todo").Collection("items")

	return &Repo{
		client: "MongoDB",
	}, nil
}

func getConnection () {
	dbContext, cancelFunc = context.WithTimeout(context.Background(), 3*time.Second)
	err := dbClient.Connect(dbContext)
	if err != nil { log.Fatal(err)}
}

func closeConnection() {
	cancelFunc()
	err := dbClient.Disconnect(dbContext)
	if err != nil { log.Fatal(err) } 
}

// List -
func (repo *Repo) List() (models.Items, error) {
	getConnection()
	defer closeConnection()

	cur, err := itemsCollection.Find(context.Background(), bson.D{})
	if err != nil { log.Fatal(err) }
	defer cur.Close(dbContext)

	items := models.Items{}
	for cur.Next(dbContext) {
		var item models.Item
		err := cur.Decode(&item)
		if err != nil { log.Fatal(err) }

		items = append(items, item)
	}

	return items, nil
}

// Get -
func (repo *Repo) Get(ID int) (*models.Item, error) {

	filter := bson.M { "id": ID }
	
	item := models.Item{}
	err := itemsCollection.FindOne(dbContext, filter).Decode(&item)
	
	if err != nil { log.Fatal(err) }

	return &item, err
}

// Create -
func (repo *Repo) Create(item *models.Item) (*models.Item, error) {

	res, err := itemsCollection.InsertOne(dbContext, item)
	if err != nil { log.Fatal(err) }

	item.ID = res.InsertedID.(int)

	return item, nil
}

// Update -
func (repo *Repo) Update(item *models.Item) (*models.Item, error) {

	filter := bson.D {{ "id", item.ID }}
	update := bson.D{
		{
			"$set", 
			bson.D {
				{"done", item.Done },
				{"title", item.Title},
		}}}

	result, err := itemsCollection.UpdateOne(dbContext, filter, update)
	if err != nil { log.Fatal(err) }

	if result.ModifiedCount == 0 {
		log.Fatal(dberrors.ErrItemNotFound)
	}

	return item, err
}

// Delete -
func (repo *Repo) Delete(ID int) error {

	filter := bson.D {{ "id", ID }}
	
	_, err := itemsCollection.DeleteOne(dbContext, filter)
	if err != nil { log.Fatal(err) }

	return err
}

func (repo *Repo) dataSeed() {
	log.Println("Seeding data...")

	items := []models.Item{
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

	var documents []interface{}
	for _, item := range items {
		documents = append(documents, item)
	}

	opts := options.InsertMany().SetOrdered(false)
	res, _ := itemsCollection.InsertMany(dbContext, documents, opts)

	log.Printf("Seeded items with IDs %v\n", res.InsertedIDs)
}
