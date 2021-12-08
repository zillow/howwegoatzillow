package main

import (
	"context"
	"log"
)

func main() {
	server, cleanup := InitializeServer()
	defer cleanup()
	if err := server.Serve(context.Background()); err != nil {
		log.Panic(err)
	}
}
