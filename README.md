# Summary:
Documents with expiration can be created via creating document with expiration date as one of its field
and then creating index over these documents that make use of said field like so:
        index := mongo.IndexModel{
	        Keys:    bson.D{{"expire", 1}},
		Options: options.Index().SetExpireAfterSeconds(0),
	}
db.col.getIndexes() will display it as: 
        {
		"v" : 2,
		"key" : {
			"expire" : 1
		},
		"name" : "expire_1",
		"expireAfterSeconds" : 0
	}

If server shuts down before expiration of documents the documents remain undeleted.
!!! Indexes are not deleted on sudden shutdown. 
They will be deleted once the server starts up again. 

NOTE: Mongo updates indexes every 60 seconds. Hence deletion (and some other actions as well)
can occur not momentarily. 