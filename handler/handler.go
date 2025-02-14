package handler

import (
	"github.com/jonp200/mongodb-demo/helpers"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type Handler struct {
	*mongo.Client
	Time helpers.Time
}
