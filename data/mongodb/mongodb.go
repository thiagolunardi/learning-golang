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

var isInitialized bool

// NewClient -
func NewClient() (*Repo, error) {

	if !isInitialized {
	
		dbClient, err := getClient()
		if err != nil { log.Fatal(err) }

		ctx, _ := getContext()
		dbClient.Connect(ctx)
		defer dbClient.Disconnect(ctx)

		err = dbClient.Ping(ctx, readpref.Primary())
		if err != nil { log.Fatal(err)}
		
		log.Println("Using MongoDb database")

		//dataSeed()

		isInitialized = true

		log.Println("MongoDb initialized.")
	}

	return &Repo{
		client: "MongoDB",
	}, nil
}

func getClient() (*mongo.Client, error) {
	return mongo.NewClient(options.Client().ApplyURI("mongodb://0.0.0.0:27017"))
}

func getItemCollection(dbClient *mongo.Client) (*mongo.Collection) {
	return dbClient.Database("todo").Collection("items")
}

func getContext() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), 30*time.Second)
}

// List -
func (repo *Repo) List(ctx1 context.Context) (models.Items, error) {
	ctx, _ := getContext()
	dbClient, err := getClient()
	if err != nil { log.Fatal() }

	err = dbClient.Connect(ctx)
	if err != nil { log.Fatal() }
	defer dbClient.Disconnect(ctx)

	itemsCollection := getItemCollection(dbClient)

	cur, err := itemsCollection.Find(ctx, bson.D{})
	if err != nil { log.Fatal(err) }
	defer cur.Close(ctx)

	items := models.Items{}
	for cur.Next(ctx) {
		var item models.Item
		err := cur.Decode(&item)
		if err != nil { log.Fatal(err) }

		items = append(items, item)
	}

	return items, nil
}

// Get -
func (repo *Repo) Get(ctx1 context.Context, ID int) (*models.Item, error) {
	ctx, _ := getContext()
	dbClient, err := getClient()
	if err != nil { log.Fatal() }

	err = dbClient.Connect(ctx)
	if err != nil { log.Fatal() }
	defer dbClient.Disconnect(ctx)

	itemsCollection := getItemCollection(dbClient)

	filter := bson.M { "id": ID }
	
	item := models.Item{}
	err = itemsCollection.FindOne(ctx, filter).Decode(&item)
	
	if err != nil { 
		if strings.Contains(err.Error(), "no documents in result") {
			return nil, nil
		}
		log.Fatal(err) 
	}

	return &item, err
}

// Create -
func (repo *Repo) Create(ctx1 context.Context,item *models.Item) (*models.Item, error) {

	ctx, _ := getContext()
	dbClient, err := getClient()
	if err != nil { log.Fatal() }

	err = dbClient.Connect(ctx)
	if err != nil { log.Fatal() }
	defer dbClient.Disconnect(ctx)

	itemsCollection := getItemCollection(dbClient)

	filter := bson.D {}
	opts := options.FindOne().SetSort(bson.D{{"id", -1}})
	lastItem := models.Item{}
	err = itemsCollection.FindOne(ctx, filter, opts).Decode(&lastItem)
	if err != nil { log.Fatal(err) }

	item.ID = lastItem.ID + 1

	_, err = itemsCollection.InsertOne(ctx, item)
	if err != nil { log.Fatal(err) }
	
	return item, nil
}

// Update -
func (repo *Repo) Update(ctx1 context.Context, item *models.Item) (*models.Item, error) {
	ctx, _ := getContext()
	dbClient, err := getClient()
	if err != nil { log.Fatal() }

	err = dbClient.Connect(ctx)
	if err != nil { log.Fatal() }
	defer dbClient.Disconnect(ctx)

	itemsCollection := getItemCollection(dbClient)

	filter := bson.D {{ "id", item.ID }}
	update := bson.D{
		{
			"$set", 
			bson.D {
				{"done", item.Done },
				{"title", item.Title},
		}}}

	result, err := itemsCollection.UpdateOne(ctx, filter, update)
	if err != nil { log.Fatal(err) }

	if result.ModifiedCount == 0 {
		log.Fatal(dberrors.ErrItemNotFound)
	}

	return item, err
}

// Delete -
func (repo *Repo) Delete(ctx1 context.Context,ID int) error {
	ctx, _ := getContext()
	dbClient, err := getClient()
	if err != nil { log.Fatal() }

	err = dbClient.Connect(ctx)
	if err != nil { log.Fatal() }
	defer dbClient.Disconnect(ctx)

	itemsCollection := getItemCollection(dbClient)

	filter := bson.D {{ "id", ID }}
	
	_, err = itemsCollection.DeleteOne(ctx, filter)
	if err != nil { log.Fatal(err) }

	return err
}

func any() bool {
	ctx, _ := getContext()
	dbClient, err := getClient()
	if err != nil { log.Fatal() }

	err = dbClient.Connect(ctx)
	if err != nil { log.Fatal() }
	defer dbClient.Disconnect(ctx)

	itemsCollection := getItemCollection(dbClient)

	filter := bson.D {{ "_id", bson.D{{"$exists", true}}}}
	cur, err := itemsCollection.Find(ctx, filter)
	if err != nil { log.Fatal(err) }
	defer cur.Close(ctx)

	return cur.Next(ctx)
}

func dataSeed() {
	
	if any() { return }

	log.Println("Seeding data...")

	ctx, _ := getContext()
	dbClient, err := getClient()
	if err != nil { log.Fatal() }

	err = dbClient.Connect(ctx)
	if err != nil { log.Fatal() }
	defer dbClient.Disconnect(ctx)

	itemsCollection := getItemCollection(dbClient)

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
	res, _ := itemsCollection.InsertMany(ctx, documents, opts)

	log.Printf("Seeded items with IDs %v\n", res.InsertedIDs)
}
