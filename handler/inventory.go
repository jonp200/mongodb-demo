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

func (h *Handler) FindInventory(c echo.Context) error {
	var m struct {
		Name string `query:"name" form:"name" validate:"not_blank"`
	}
	if err := c.Bind(&m); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	if err := c.Validate(m); err != nil {
		return echo.NewHTTPError(http.StatusUnprocessableEntity, err.Error())
	}

	// Select database and collection
	col := h.Client.Database(datastore.DbHobbyShop).Collection("inventory")

	// Define filter (empty filter retrieves all documents)
	f := bson.M{"short_name": filter.BeginsWith(m.Name)}

	// Execute Find query
	cursor, err := col.Find(c.Request().Context(), f)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return echo.NewHTTPError(http.StatusNotFound, err)
		}
		log.Panic(err)
	}

	var results []model.Inventory
	if err = cursor.All(c.Request().Context(), &results); err != nil {
		log.Panic(err)
	}

	// Check for errors during iteration
	if err = cursor.Err(); err != nil {
		log.Panic(err)
	}

	if len(results) == 0 {
		return echo.NewHTTPError(http.StatusNotFound, fmt.Sprintf("No results found for: %s", m.Name))
	}

	return c.JSON(
		http.StatusOK, map[string]any{
			"data": results,
		},
	)
}
