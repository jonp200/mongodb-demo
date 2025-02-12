package handler

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/jonp200/mongodb-demo/datastore"
	"github.com/jonp200/mongodb-demo/datastore/filter"
	"github.com/jonp200/mongodb-demo/model"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type Handler struct {
	Client *mongo.Client
}

func (h *Handler) FindByTitle(c echo.Context) error {
	var m model.Movie
	if err := c.Bind(&m); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	if err := c.Validate(m); err != nil {
		return echo.NewHTTPError(http.StatusUnprocessableEntity, err.Error())
	}

	// Select database and collection
	col := datastore.Collection(h.Client, "movies")

	// Define filter (empty filter retrieves all documents)
	f := bson.M{"title": filter.BeginsWith(m.Title)}

	// Execute Find query
	cursor, err := col.Find(c.Request().Context(), f)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return echo.NewHTTPError(http.StatusNotFound, err)
		}
		log.Panic(err)
	}

	var results []bson.D
	if err = cursor.All(c.Request().Context(), &results); err != nil {
		log.Panic(err)
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
