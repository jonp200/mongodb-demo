package handler

import "go.mongodb.org/mongo-driver/v2/mongo"

type Handler struct {
	*mongo.Client
}
