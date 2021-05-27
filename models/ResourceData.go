package models

import (
	"github.com/lucashmorais/go-interscity/database"
	"go.mongodb.org/mongo-driver/mongo"
)

var ResourceDataCollection *mongo.Collection

type ResourceData struct {
	UUID string                 `bson: uuid`
	Data map[string]interface{} `bson: "data"`
}

var ResourceDataDataCollection *mongo.Collection

func CreateResourceDataSchema() {
	database.DB.CreateCollection(database.Ctx, "resource_data")
	ResourceDataCollection = database.DB.Collection("resource_data")
}
