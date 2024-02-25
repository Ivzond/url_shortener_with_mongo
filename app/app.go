package app

import (
	"context"
	"fmt"
	"github.com/gocraft/web"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"net/http"
)

func Run(ctx context.Context) error {
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://127.0.0.1:27017"))
	if err != nil {
		return err
	}

	defer func() {
		if err := client.Disconnect(ctx); err != nil {
			log.Println(err)
		}
	}()

	shortUrlDAO, err := NewUrlDAO(ctx, client)
	if err != nil {
		return err
	}
	service := NewService(shortUrlDAO)
	httpHandler := NewHandler(service)
	fmt.Println("App is working on port :8000")
	return http.ListenAndServe("localhost:8000", initEndpoints(httpHandler))
}

func initEndpoints(h *Handler) *web.Router {
	router := web.New(*h)
	router.Post("/shorten", WrapEndpoint(h.Shorten))
	router.Get("/:shortUrl", WrapEndpoint(h.GetFullURL))
	router.Post("/update/:shortUrl", WrapEndpoint(h.Update))
	router.Delete("/:shortUrl", WrapEndpoint(h.Delete))
	router.Get("/ping", WrapEndpoint(h.Ping))
	return router
}
