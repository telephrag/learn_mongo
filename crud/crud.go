package crud

import (
	"context"
	"db_test/models"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func InsertAndDisplay(collection *mongo.Collection) *[]models.Player { // inserts documents into collection and shows them

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
