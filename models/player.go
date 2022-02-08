package models

import "go.mongodb.org/mongo-driver/bson/primitive"

// bson tags you specified here will be used by database

type Player struct {
	ID       primitive.ObjectID `bson:"_id,omitempty"`
	Username string             `bson:"name"`
	Expire   primitive.DateTime `bson:"expire"`
}
