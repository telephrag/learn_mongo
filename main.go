package main

import (
	"context"
	"db_test/config"
	"db_test/models"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// db.createCollection("player", {size: 1048576, capped: true})

func insertAndDisplay(collection *mongo.Collection) *[]models.Player { // inserts documents into collection and shows them

	niki := models.Player{
		Username: "niki",
		Expire:   primitive.NewDateTimeFromTime(time.Now().Add(time.Second * 30).UTC()),
	}
	charon := models.Player{
		Username: "charon",
		Expire:   primitive.NewDateTimeFromTime(time.Now().Add(time.Second * 30).UTC()),
	}
	skif := models.Player{
		Username: "skif",
		Expire:   primitive.NewDateTimeFromTime(time.Now().Add(time.Second * 30).UTC()),
	}

	players := []interface{}{niki, charon, skif}
	imr, err := collection.InsertMany(context.TODO(), players)

	// imr stores ids of inserted documents
	if err != nil {
		log.Panic(err)
	}
	fmt.Println(imr)

	coursor, err := collection.Find( // find docs with ids in imr
		context.TODO(),
		bson.M{
			"_id": bson.M{
				"$in": imr.InsertedIDs,
			},
		},
	)
	if err != nil {
		log.Panic(err)
	}
	var documents []models.Player
	err = coursor.All(context.TODO(), &documents)
	if err != nil {
		log.Panic(err)
	}

	idToUpdate := imr.InsertedIDs[0]
	_, err = collection.UpdateOne(
		context.TODO(),
		bson.M{"_id": idToUpdate},
		bson.A{
			bson.M{
				"$set": bson.D{{
					"name", bson.M{"$concat": bson.A{"[GBC]", "$name"}},
				}},
			},
		},
	)
	if err != nil {
		log.Panic(err)
	}

	return &documents
}

func findLongBoi(collection *mongo.Collection) *[]models.Player { // find players with names longer than 4 symbols

	longName := bson.M{ // make filter
		"name": bson.M{"$exists": true}, // find name
		"$expr": bson.M{ // that satisfy expression
			"$gt": bson.A{ // greater than
				bson.M{"$strLenCP": "$name"}, // length of name
				4,                            // greater than 4
			},
		},
	}
	coursor, err := collection.Find(context.TODO(), longName)
	if err != nil {
		log.Panic(err)
	}
	// since find returns coursor we need to use All() to iterate over documents
	var longBoi []models.Player
	err = coursor.All(context.TODO(), &longBoi)
	if err != nil {
		log.Panic(err)
	}

	return &longBoi
}

func makeTTLIndex(collection *mongo.Collection) string {
	index := mongo.IndexModel{
		Keys:    bson.D{{"expire", 1}},
		Options: options.Index().SetExpireAfterSeconds(0),
	}

	str, err := collection.Indexes().CreateOne(context.TODO(), index)
	if err != nil {
		log.Panic(err)
	}

	return str
}

func iterateChangeStream(ctx context.Context, stream *mongo.ChangeStream, cancel context.CancelFunc) {
	defer stream.Close(ctx)

	for stream.Next(ctx) {
		var event bson.M
		err := stream.Decode(&event)
		if err != nil {
			log.Panic(err)
		}
		fmt.Println(event)
	}
}

func watchEvents(collection *mongo.Collection) {

	pipeline := mongo.Pipeline{
		bson.D{{
			"$match",
			bson.D{{
				"$or", bson.A{
					bson.D{{"operationType", "insert"}},
					bson.D{{"operationType", "delete"}},
				},
			}},
		}},
	}

	stream, err := collection.Watch(context.TODO(), pipeline)
	if err != nil {
		log.Panic(err)
	}
	// defer stream.Close(context.TODO()) ???

	ctx, cancelFunc := context.WithCancel(context.Background())
	go iterateChangeStream(ctx, stream, cancelFunc)
}

func main() {

	clientOptions := options.Client().ApplyURI(
		"mongodb://localhost:27017",
	)

	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Panic(err)
	}
	defer client.Disconnect(context.TODO())

	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Panic(err)
	}

	// create collection and fill it with data with expiration times
	collection := client.Database(config.DBName).Collection(config.CollectionName)

	watchEvents(collection)

	documents := insertAndDisplay(collection)
	fmt.Println("Inserted documents:\n", documents)

	longBoi := findLongBoi(collection)
	fmt.Printf("Longboi: %v\n", longBoi)

	// create TTL index
	indexName := makeTTLIndex(collection)
	fmt.Println("Created index: ", indexName)

	// collection contents can be tracked via db.player.find() command from mongo console
	// where db is database that contains this collection

	interrupt := make(chan os.Signal, 1) // graceful shutdown
	signal.Notify(interrupt, syscall.SIGTERM, syscall.SIGINT)
	<-interrupt

	err = collection.Drop(context.TODO()) // cleanup
	if err != nil {
		log.Panic(err)
	}
}
