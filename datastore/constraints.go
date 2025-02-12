package datastore

import (
	"context"

	"github.com/labstack/gommon/log"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

func Index(client *mongo.Client) {
	movies := Collection(client, "movies")

	titleIX := mongo.IndexModel{Keys: bson.D{{"title", 1}}}
	if _, err := movies.Indexes().CreateOne(context.Background(), titleIX); err != nil {
		log.Fatal(err)
	}
}
