package model

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type Inventory struct {
	ID        bson.ObjectID `bson:"_id" json:"id"`
	ShortName string        `bson:"short_name" json:"short_name"`
	FullName  string        `bson:"full_name" json:"full_name"`
	Status    string        `bson:"status" json:"status"`
	Stock     int           `bson:"stock" json:"stock"`
	CreatedAt time.Time     `bson:"created_at" json:"created_at"`
	UpdatedAt *time.Time    `bson:"updated_at" json:"updated_at"`
}
