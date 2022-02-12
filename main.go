package main

import (
	"context"
	"db_test/changestream"
	"db_test/config"
	"db_test/crud"
	"db_test/longboi"
	"db_test/ttl"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// db.createCollection("player", {size: 1048576, capped: true})

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

	ctx, cancel := context.WithCancel(context.Background())
	go changestream.WatchEvents(collection, ctx, cancel)

	crud.InsertAndDisplay(collection)
	//fmt.Println("Inserted documents:\n", documents)

	longboi.FindLongBoi(collection)
	//fmt.Println("Longboi:", longBoi)

	// create TTL index
	ttl.MakeTTLIndex(collection)
	//fmt.Println("Created index: ", indexName)

	// collection contents can be tracked via db.player.find() command from mongo console
	// where db is database that contains this collection

	interrupt := make(chan os.Signal, 1) // graceful shutdown
	signal.Notify(interrupt, syscall.SIGTERM, syscall.SIGINT)
	<-interrupt

	err = collection.Drop(context.TODO()) // cleanup
	if err != nil {
		log.Panic(err)
	}

	// handling of operation do not occur if <-time.After() is used
	timeout := time.After(10 * time.Second)
	for {
		select {
		case <-ctx.Done():
			fmt.Println("handled invalidation at shutdown")
			return
		case <-timeout:
			fmt.Println("invalidation wasn't handled or didn't occur")
			return
		default:
		}
		time.Sleep(time.Millisecond * 500)
	}
}

// make program shutdown after writing to channel that invalidate event was handled
