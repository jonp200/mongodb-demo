package model

import "go.mongodb.org/mongo-driver/v2/bson"

type Inventory struct {
	ID        bson.ObjectID `bson:"_id" json:"id" validate:"required"`
	ShortName string        `bson:"short_name" json:"short_name" validate:"not_blank"`
	FullName  string        `bson:"full_name" json:"full_name" validate:"not_blank"`
	Status    string        `bson:"status" json:"status" validate:"not_blank"`
	Stock     int           `bson:"stock" json:"stock" validate:"gte=0"`
	CreatedAt string        `bson:"created_at" json:"created_at" validate:"not_blank,format=2006-01-02T15:04:05Z07:00"`
	UpdatedAt *string       `bson:"updated_at" json:"updated_at" validate:"format=2006-01-02T15:04:05Z07:00"`
}
