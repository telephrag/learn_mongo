package longboi

import (
	"context"
	"db_test/models"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func FindLongBoi(collection *mongo.Collection) *[]models.Player { // find players with names longer than 4 symbols

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
