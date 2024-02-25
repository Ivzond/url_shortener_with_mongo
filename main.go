package main

import (
	"context"
	"url_shortener_with_mongo/app"
)

func main() {
	if err := app.Run(context.Background()); err != nil {
		panic(err)
	}
}
