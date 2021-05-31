package models

import (
	"github.com/lucashmorais/go-interscity/database"
	"go.mongodb.org/mongo-driver/mongo"
)

var ResourceCollection *mongo.Collection

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
	database.DB.CreateCollection(database.Ctx, "resources")
	ResourceCollection = database.DB.Collection("resources")
}
