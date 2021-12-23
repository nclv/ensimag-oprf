package main

import (
	"context"
	"github.com/cloudflare/circl/oprf"
	"github.com/ensimag-oprf/go/server/controllers"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/ensimag-oprf/go/server/routers"
)

const (
	HOST = "localhost"
	PORT = "1323"
)

// loadPrivateKeysFromEnv load the base64 serialized private keys from the environment variables.
func loadPrivateKeysFromEnv() controllers.SerializedBase64KeyMap {
	serializedBase64KeyMap := make(controllers.SerializedBase64KeyMap)

	serializedBase64P256PrivateKey, ok := os.LookupEnv("P256_PRIVATE_KEY")
	if ok {
		serializedBase64KeyMap[oprf.OPRFP256] = serializedBase64P256PrivateKey
	}

	serializedBase64P384PrivateKey, ok := os.LookupEnv("P384_PRIVATE_KEY")
	if ok {
		serializedBase64KeyMap[oprf.OPRFP384] = serializedBase64P384PrivateKey
	}

	serializedBase64P521PrivateKey, ok := os.LookupEnv("P521_PRIVATE_KEY")
	if ok {
		serializedBase64KeyMap[oprf.OPRFP521] = serializedBase64P521PrivateKey
	}

	return serializedBase64KeyMap
}

func main() {
	serializedBase64KeyMap := loadPrivateKeysFromEnv()

	router, err := routers.NewRouter(serializedBase64KeyMap)
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
