package datastore

import (
	"os"

	"github.com/labstack/gommon/log"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func Connect() *mongo.Client {
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

func Collection(client *mongo.Client, col string) *mongo.Collection {
	return client.Database(db).Collection(col)
}
