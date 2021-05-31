package models

import (
	"github.com/lucashmorais/go-interscity/database"
	"go.mongodb.org/mongo-driver/mongo"
)

type CoreSubscription struct {
	Uuid         string   `bson: uuid`
	Capabilities []string `bson: capabilities`
	Url          string   `bson: url`
}

type Subscription struct {
	Subscription CoreSubscription `bson: subscription`
}

var SubscriptionCollection *mongo.Collection

func CreateSubscriptionSchema() {
	database.DB.CreateCollection(database.Ctx, "subscriptions")
	SubscriptionCollection = database.DB.Collection("subscriptions")
}
