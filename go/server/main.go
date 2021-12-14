package main

import "github.com/labstack/echo/v4"

const (
	HOST = "localhost"
	PORT = "1323"
)

func main() {
	e := echo.New()

	// TODO: https://echo.labstack.com/cookbook/auto-tls/

	server := NewServer()
	server.Initialize()

	e.GET("/request_public_keys", server.getKeys)
	e.POST("/evaluate", server.evaluate)

	e.Logger.Fatal(e.Start(HOST + ":" + PORT))
}
