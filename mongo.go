package main

import (
	"context"
	"os"

	"github.com/labstack/gommon/log"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func connect() *mongo.Client {
	uri := os.Getenv("MONGODB_URI")
	docs := "www.mongodb.com/docs/drivers/go/current/"
	if uri == "" {
		log.Fatal(
			"Set your 'MONGODB_URI' environment variable. " +
				"See: " + docs +
				"usage-examples/#environment-variable",
		)
	}
	client, err := mongo.Connect(options.Client().ApplyURI(uri))
	if err != nil {
		log.Panic(err)
	}

	return client
}

const db = "sample_mflix"

func collection(client *mongo.Client, col string) *mongo.Collection {
	return client.Database(db).Collection(col)
}

func beginsWith(a string) bson.M {
	return bson.M{"$regex": "^" + a, "$options": "i"}
}

func index(client *mongo.Client) {
	movies := collection(client, "movies")

	titleIX := mongo.IndexModel{Keys: bson.D{{"title", 1}}}
	if _, err := movies.Indexes().CreateOne(context.Background(), titleIX); err != nil {
		log.Fatal(err)
	}
}
