package ttl

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func MakeTTLIndex(collection *mongo.Collection) string {
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
