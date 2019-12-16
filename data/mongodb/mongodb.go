package mongodb

import (
	"context"
	"log"
	"strings"
	"time"

	"github.com/thiagolunardi/learning-golang/data/dberrors"
	"github.com/thiagolunardi/learning-golang/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

// Repo -
type Repo struct {
	client string
}

type disconnectFunc func() error

var isInitialized bool

// NewClient -
func NewClient() (*Repo, error) {

	if !isInitialized {

		dbClient, err := getClient()
		if err != nil {
			log.Fatal(err)
		}

		ctx, _ := getContext()
		dbClient.Connect(ctx)
		defer dbClient.Disconnect(ctx)

		err = dbClient.Ping(ctx, readpref.Primary())
		if err != nil {
			log.Fatal(err)
		}

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

func getItemCollection(dbClient *mongo.Client) *mongo.Collection {
	return dbClient.Database("todo").Collection("items")
}

func getContext() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), 30*time.Second)
}

func getCollection(ctx context.Context) (*mongo.Collection, disconnectFunc) {
	dbClient, err := getClient()
	if err != nil {
		log.Fatal()
	}

	err = dbClient.Connect(ctx)
	if err != nil {
		log.Fatal()
	}

	itemsCollection := getItemCollection(dbClient)

	return itemsCollection, func() error {
		return dbClient.Disconnect(ctx)
	}
}

// List -
func (repo *Repo) List(ctx context.Context) (models.Items, error) {
	itemsCollection, disconnect := getCollection(ctx)
	defer disconnect()

	c := make(chan models.Item)
	go listAsync(ctx, itemsCollection, c)

	items := models.Items{}
	for item := range c {
		items = append(items, item)
	}

	return items, nil
}

func listAsync(
	ctx context.Context,
	itemsCollection *mongo.Collection,
	c chan models.Item) {

	cur, err := itemsCollection.Find(ctx, bson.D{})
	if err != nil {
		log.Fatal(err)
	}
	defer cur.Close(ctx)

	for cur.Next(ctx) {
		var item models.Item
		err := cur.Decode(&item)
		if err != nil {
			log.Fatal(err)
		}

		c <- item
	}
	close(c)
}

// Get -
func (repo *Repo) Get(ctx context.Context, ID int) (*models.Item, error) {
	itemsCollection, disconnect := getCollection(ctx)
	defer disconnect()

	c := make(chan *models.Item, 1)
	go getAsync(ctx, itemsCollection, ID, c)
	item := <-c

	return item, nil
}

func getAsync(
	ctx context.Context,
	itemsCollection *mongo.Collection,
	ID int,
	c chan *models.Item) {
	filter := bson.M{"id": ID}

	item := models.Item{}
	err := itemsCollection.FindOne(ctx, filter).Decode(&item)

	if err == nil {
		c <- &item
	} else {
		if strings.Contains(err.Error(), "no documents in result") {
			c <- nil
		}
		close(c)
		log.Fatal(err)
	}
	close(c)
}

// Create -
func (repo *Repo) Create(ctx context.Context, item *models.Item) (*models.Item, error) {
	itemsCollection, disconnect := getCollection(ctx)
	defer disconnect()

	c := make(chan *models.Item, 1)
	go createAsync(ctx, itemsCollection, item, c)

	createdItem := <-c

	return createdItem, nil
}

func createAsync(
	ctx context.Context,
	itemsCollection *mongo.Collection,
	item *models.Item,
	c chan *models.Item) {
	filter := bson.D{}
	opts := options.FindOne().SetSort(bson.D{primitive.E{Key: "id", Value: -1}})
	lastItem := models.Item{}
	err := itemsCollection.FindOne(ctx, filter, opts).Decode(&lastItem)
	if err != nil {
		log.Fatal(err)
	}

	item.ID = lastItem.ID + 1

	_, err = itemsCollection.InsertOne(ctx, item)
	if err != nil {
		log.Fatal(err)
	}

	c <- item
	close(c)
}

// Update -
func (repo *Repo) Update(ctx context.Context, item *models.Item) (*models.Item, error) {
	itemsCollection, disconnect := getCollection(ctx)
	defer disconnect()

	c := make(chan *models.Item, 1)
	go updateAsync(ctx, itemsCollection, item, c)

	updatedItem := <-c
	return updatedItem, nil
}

// Update -
func updateAsync(
	ctx context.Context,
	itemsCollection *mongo.Collection,
	item *models.Item,
	c chan *models.Item) {

	filter := bson.D{primitive.E{Key: "id", Value: item.ID}}
	update := bson.D{primitive.E{
		Key: "$set",
		Value: bson.D{
			primitive.E{Key: "done", Value: item.Done},
			primitive.E{Key: "title", Value: item.Title},
		}}}

	result, err := itemsCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		close(c)
		log.Fatal(err)
	}

	if result.ModifiedCount == 0 {
		close(c)
		log.Fatal(dberrors.ErrItemNotFound)
	}

	c <- item
	close(c)
}

// Delete -
func (repo *Repo) Delete(ctx context.Context, ID int) error {
	itemsCollection, disconnect := getCollection(ctx)
	defer disconnect()

	c := make(chan int64, 1)
	go deleteAsync(ctx, itemsCollection, ID, c)
	deleted := <- c

	log.Printf("%d deleted", deleted)

	return nil
}

func deleteAsync(
	ctx context.Context, 
	itemsCollection *mongo.Collection,
	ID int,
	c chan int64) {

	filter := bson.D{primitive.E{Key: "id", Value: ID}}

	r, err := itemsCollection.DeleteOne(ctx, filter)
	if err != nil {
		close(c)
		log.Fatal(err)
	}
	
	c <- r.DeletedCount
}


func any(ctx context.Context) bool {
	itemsCollection, disconnect := getCollection(ctx)
	defer disconnect()

	filter := bson.D{primitive.E{Key: "_id", Value: bson.D{primitive.E{Key: "$exists", Value: true}}}}
	cur, err := itemsCollection.Find(ctx, filter)
	if err != nil {
		log.Fatal(err)
	}
	defer cur.Close(ctx)

	return cur.Next(ctx)
}

func dataSeed(ctx context.Context) {

	if any(ctx) {
		return
	}

	log.Println("Seeding data...")

	dbClient, err := getClient()
	if err != nil {
		log.Fatal()
	}

	err = dbClient.Connect(ctx)
	if err != nil {
		log.Fatal()
	}
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
