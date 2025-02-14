package migrations

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/labstack/gommon/log"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func Apply(ctx context.Context, db *mongo.Database) error {
	mc := db.Collection("migrations")
	for _, m := range migrations {
		var existing bson.M
		err := mc.FindOne(ctx, bson.M{"_id": m.id}).Decode(&existing)
		if errors.Is(err, mongo.ErrNoDocuments) {
			log.Printf("Applying migration: %s", m.id)

			if err = m.up(ctx, db); err != nil {
				return fmt.Errorf("migration %s failed: %w", m.id, err)
			}

			if _, err = mc.InsertOne(
				ctx, bson.M{"_id": m.id, "applied_at": time.Now().UTC().Format(time.RFC3339)},
			); err != nil {
				return err
			}
		}
	}

	return nil
}

type migration struct {
	id   string
	up   func(ctx context.Context, db *mongo.Database) error
	down func(ctx context.Context, db *mongo.Database) error
}

var migrations = []migration{
	{
		id: "20250212072636_initial_create",
		up: func(ctx context.Context, db *mongo.Database) error {
			err := db.CreateCollection(context.Background(), "inventory")
			if err != nil {
				return err
			}

			_, err = db.Collection("inventory").InsertOne(
				ctx, bson.D{
					{"short_name", "Rare item"},
					{"full_name", "Rare item - Limited collection"},
					{"status", "In stock"},
					{"stock", 1},
					{"created_at", time.Now().UTC()},
				},
			)
			if err != nil {
				return err
			}

			invCol := db.Collection("inventory")

			invShortNameIX := mongo.IndexModel{
				Keys:    bson.D{{"short_name", 1}},
				Options: options.Index().SetUnique(true),
			}
			if _, err = invCol.Indexes().CreateOne(context.Background(), invShortNameIX); err != nil {
				return err
			}

			invFullNameIX := mongo.IndexModel{
				Keys:    bson.D{{"full_name", 1}},
				Options: options.Index().SetUnique(true),
			}
			if _, err = invCol.Indexes().CreateOne(context.Background(), invFullNameIX); err != nil {
				return err
			}

			return nil
		},
		down: func(ctx context.Context, db *mongo.Database) error {
			err := db.Collection("inventory").Drop(context.Background())
			if err != nil {
				return err
			}

			return nil
		},
	},
}
