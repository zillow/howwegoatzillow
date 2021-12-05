package main

import (
	"context"
	"log"
)

func main() {
	w, cleanup := InitializeServer()
	defer cleanup()
	if err := w.Serve(context.Background()); err != nil {
		log.Panic(err)
	}
}
