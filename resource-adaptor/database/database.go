package database

import (
	"context"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Ctx | @desc: mongo context interface
var Ctx context.Context

// Cancel | @desc: mongo context cancel function
var Cancel context.CancelFunc

// Client | @desc: mongo client struct
var Client *mongo.Client

// DB | @desc: mongo database struct
var DB *mongo.Database

// var DB_PORT int

func getMongoURI() string {
	base := os.Getenv(("MONGO_URI_BASE"))
	port := os.Getenv(("MONGO_PORT"))
	suffix := os.Getenv(("MONGO_URI_SUFFIX"))

	return base + port + suffix
}

// Connect | @desc: connects to mongoDB
func Connect() {
	var err error

	Ctx, Cancel = context.WithTimeout(context.Background(), 30*time.Second)
	// Client, err = mongo.Connect(Ctx, options.Client().ApplyURI(os.Getenv("MONGO_URI")))
	Client, err = mongo.Connect(Ctx, options.Client().ApplyURI(getMongoURI()))
	if err != nil {
		panic(err)
	}

	// Connect to mongo database
	DB = Client.Database("go-test")
}
