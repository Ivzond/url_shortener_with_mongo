package app

import (
	"context"
	"fmt"
	"github.com/gocraft/web"
	"net/http"
)

func Run(ctx context.Context) error {
	service := NewService()
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
