package main

import (
	"context"
	"github.com/oprf/go/server/api"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

const (
	HOST = "localhost"
	PORT = "1323"
)

func main() {
	router, err := api.NewRouter()
	if err != nil {
		log.Println(err)

		return
	}

	// Start the server
	go func() {
		if err := router.Start(HOST + ":" + PORT); err != nil && err != http.ErrServerClosed {
			router.Logger.Fatal("shutting down the server")
		}
	}()

	// Wait for interrupt signal to gracefully shut down the server with a timeout of 10 seconds.
	// Use a buffered channel to avoid missing signals as recommended for signal.Notify
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := router.Shutdown(ctx); err != nil {
		router.Logger.Fatal(err)
	}
}
