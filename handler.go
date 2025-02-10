package main

import (
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type handler struct {
	client *mongo.Client
}

func (h *handler) FindByTitle(c echo.Context) error {
	var m Movie
	if err := c.Bind(&m); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	if err := c.Validate(m); err != nil {
		return echo.NewHTTPError(http.StatusUnprocessableEntity, err.Error())
	}

	col := collection(h.client, "movies")

	var result bson.M
	if err := col.FindOne(c.Request().Context(), bson.D{{"title", m.Title}}).Decode(&result); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			log.Printf("No document was found with the title %q", m.Title)
			return c.NoContent(http.StatusNotFound)
		}
		log.Panic(err)
	}

	return c.JSON(http.StatusOK, result)
}
