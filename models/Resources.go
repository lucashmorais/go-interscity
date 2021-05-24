package models

import (
	"github.com/lucashmorais/go-interscity/database"
	"go.mongodb.org/mongo-driver/mongo"
)

// ResourceCollection | @desc: the resource ccollection on the database
var ResourceCollection *mongo.Collection

// Resource | @desc: resource model struct
/*
type Resource struct {
	ID    primitive.ObjectID `bson:"_id,omitempty"`
	Name  string             `bson:"name,omitempty"`
	Email string             `bson:"email,omitempty"`
}
*/

type DBResource struct {
	Data string `bson: "data"`
	// Data map[string]interface{} `json: "data"`
}

type Resource struct {
	UUID string                 `_id`
	Data map[string]interface{} `bson: "data"`
}

// CreateResourceSchema | @desc: adds schema validation and indexes to collection
func CreateResourceSchema() {
	/*
		jsonSchema := bson.M{
			"bsonType": "object",
			"required": []string{"name", "email"},
			"properties": bson.M{
				"name": bson.M{
					"bsonType":    "string",
					"description": "must be a string and is required",
				},
				"email": bson.M{
					"bsonType":    "string",
					"description": "must be a string and is required",
				},
			},
		}

		validator := bson.M{
			"$jsonSchema": jsonSchema,
		}

		database.DB.CreateCollection(database.Ctx, "resources", options.CreateCollection().SetValidator(validator))
	*/

	database.DB.CreateCollection(database.Ctx, "resources")
	ResourceCollection = database.DB.Collection("resources")
}
