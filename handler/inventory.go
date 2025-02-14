package handler

import (
	"errors"
	"net/http"

	"github.com/jonp200/mongodb-demo/datastore"
	"github.com/jonp200/mongodb-demo/model"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func (h *Handler) FindInventory(c echo.Context) error {
	var m struct {
		Name string `query:"name" validate:"not_blank"`
	}
	if err := c.Bind(&m); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	if err := c.Validate(m); err != nil {
		return echo.NewHTTPError(http.StatusUnprocessableEntity, err.Error())
	}

	coll := h.Client.Database(datastore.DbHobbyShop).Collection("inventory")

	filter := bson.M{"short_name": bson.M{"$regex": m.Name, "$options": "i"}}

	opts := options.Find().SetSort(bson.M{"full_name": 1})

	cursor, err := coll.Find(c.Request().Context(), filter, opts)
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

func (h *Handler) FindInventoryByID(c echo.Context) error {
	var m struct {
		ID bson.ObjectID `param:"id" validate:"not_blank"`
	}
	if err := c.Bind(&m); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	if err := c.Validate(m); err != nil {
		return echo.NewHTTPError(http.StatusUnprocessableEntity, err.Error())
	}

	coll := h.Client.Database(datastore.DbHobbyShop).Collection("inventory")

	filter := bson.M{"_id": bson.M{"$eq": m.ID}}

	var result model.Inventory
	if err := coll.FindOne(c.Request().Context(), filter).Decode(&result); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return echo.NewHTTPError(http.StatusNotFound, err)
		}
		log.Panic(err)
	}

	return c.JSON(http.StatusOK, result)
}

func (h *Handler) AddInventory(c echo.Context) error {
	var m struct {
		ShortName string `json:"short_name" validate:"not_blank"`
		FullName  string `json:"full_name" validate:"not_blank"`
		Status    string `json:"status" validate:"not_blank"`
		Stock     int    `json:"stock" validate:"gte=0"`
	}
	if err := c.Bind(&m); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	if err := c.Validate(m); err != nil {
		return echo.NewHTTPError(http.StatusUnprocessableEntity, err.Error())
	}

	coll := h.Client.Database(datastore.DbHobbyShop).Collection("inventory")

	doc := model.Inventory{
		ID:        bson.NewObjectID(),
		ShortName: m.ShortName,
		FullName:  m.FullName,
		Status:    m.Status,
		Stock:     m.Stock,
		CreatedAt: h.Time.Now(),
	}

	result, err := coll.InsertOne(c.Request().Context(), doc)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(
		http.StatusCreated, map[string]any{
			"id": result.InsertedID,
		},
	)
}
