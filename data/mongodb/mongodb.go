package mongodb

import (
	"strings"
	"context"
	"log"
	"time"

	"github.com/thiagolunardi/learning-golang/data/dberrors"
	"github.com/thiagolunardi/learning-golang/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
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

	pingServer()
	log.Println("Using MongoDb database")

	itemsCollection = dbClient.Database("todo").Collection("items")

	dataSeed()

	return &Repo{
		client: "MongoDB",
	}, nil
}

func openConnection () {
	var err error
	dbClient, err = mongo.Connect(getContext(), options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil { log.Fatal(err)}
}

func closeConnection() {
	err := dbClient.Disconnect(dbContext)
	if err != nil { log.Fatal(err) } 
}

func pingServer() {
	var err error
	dbClient, err = mongo.Connect(getContext(), options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil { log.Fatal(err)}

	err = dbClient.Ping(getContext(), readpref.Primary())
}

func getContext() (context.Context) {
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)

	return ctx
}

// List -
func (repo *Repo) List() (models.Items, error) {
	openConnection()
	defer closeConnection()

	cur, err := itemsCollection.Find(getContext(), bson.D{})
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
	openConnection()
	defer closeConnection()

	filter := bson.M { "id": ID }
	
	item := models.Item{}
	err := itemsCollection.FindOne(dbContext, filter).Decode(&item)
	
	if err != nil { 
		if strings.Contains(err.Error(), "no documents in result") {
			return nil, nil
		}
		log.Fatal(err) 
	}

	return &item, err
}

// Create -
func (repo *Repo) Create(item *models.Item) (*models.Item, error) {
	openConnection()
	defer closeConnection()

	res, err := itemsCollection.InsertOne(dbContext, item)
	if err != nil { log.Fatal(err) }

	item.ID = res.InsertedID.(int)

	return item, nil
}

// Update -
func (repo *Repo) Update(item *models.Item) (*models.Item, error) {
	openConnection()
	defer closeConnection()

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
	openConnection()
	defer closeConnection()

	filter := bson.D {{ "id", ID }}
	
	_, err := itemsCollection.DeleteOne(dbContext, filter)
	if err != nil { log.Fatal(err) }

	return err
}

func any() bool {
	openConnection()
	defer closeConnection()

	filter := bson.D {{ "_id", bson.D{{"$exists", true}}}}
	ctx := getContext()
	cur, err := itemsCollection.Find(ctx, filter)
	if err != nil { log.Fatal(err) }
	defer cur.Close(ctx)

	return cur.Next(ctx)
}

func dataSeed() {
	
	if any() { return }

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
