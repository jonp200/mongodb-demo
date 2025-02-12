package main

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"

	"github.com/joho/godotenv"
	"github.com/jonp200/mongodb-demo/datastore"
	"github.com/jonp200/mongodb-demo/handler"
	"github.com/jonp200/mongodb-demo/model"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
)

func main() {
	log.SetHeader("${time_rfc3339} ${level} ${short_file}:${line}")

	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}

	// Initialise MongoDB client once
	client := datastore.Connect()

	defer func() {
		log.Print("Disconnecting from MongoDB...")
		if err := client.Disconnect(context.Background()); err != nil {
			log.Fatal(err)
		}

		log.Print("MongoDB disconnected.")
	}()

	// Create indexes
	datastore.Index(client)

	e := echo.New()

	e.Validator = model.Validator()

	e.GET(
		"/", func(c echo.Context) error {
			return c.String(http.StatusOK, "MongoDB demo")
		},
	)

	h := handler.Handler{Client: client}

	e.GET("/movies", h.FindByTitle)

	// Handle shutdown gracefully
	go func() {
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, os.Interrupt)
		<-sigChan

		log.Print("Shutting down gracefully...")
		if err := e.Shutdown(context.Background()); err != nil {
			log.Fatal(err)
		}
	}()

	if err := e.Start(":8080"); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatal(err)
	}
}
