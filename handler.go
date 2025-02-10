package main

import (
	"errors"
	"fmt"
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

	// Select database and collection
	col := collection(h.client, "movies")

	// Define filter (empty filter retrieves all documents)
	filter := bson.M{"title": beginsWith(m.Title)}

	// Execute Find query
	cursor, err := col.Find(c.Request().Context(), filter)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return echo.NewHTTPError(http.StatusNotFound, err)
		}
		log.Panic(err)
	}

	var results []bson.M

	// Iterate through the cursor
	for cursor.Next(c.Request().Context()) {
		var result bson.M
		if err = cursor.Decode(&result); err != nil {
			log.Panic(err)
		}
		results = append(results, result)
	}

	// Check for errors during iteration
	if err = cursor.Err(); err != nil {
		log.Panic(err)
	}

	if len(results) == 0 {
		return echo.NewHTTPError(http.StatusNotFound, fmt.Sprintf("No results found for title: %s", m.Title))
	}

	return c.JSON(http.StatusOK, results)
}
