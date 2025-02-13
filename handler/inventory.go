package handler

import (
	"errors"
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
		ID   string `query:"id" form:"id" validate:"required_without=Name|uuid4"`
		Name string `query:"name" form:"name" validate:"required_without=ID"`
	}
	if err := c.Bind(&m); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	if err := c.Validate(m); err != nil {
		return echo.NewHTTPError(http.StatusUnprocessableEntity, err.Error())
	}

	// Select database and collection
	col := h.Client.Database(datastore.DbHobbyShop).Collection("inventory")

	// Initialise with ID filter
	f := bson.M{"_id": bson.M{"$eq": m.ID}}

	// If ID filter is not provided, use the other filters combined
	if m.ID == "" {
		f = bson.M{"short_name": filter.BeginsWith(m.Name)}
	}

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
		return echo.NewHTTPError(http.StatusNotFound, "Inventory item not found")
	}

	return c.JSON(
		http.StatusOK, map[string]any{
			"data": results,
		},
	)
}
