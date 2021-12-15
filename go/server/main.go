package main

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

const (
	HOST = "localhost"
	PORT = "1323"
)

func main() {
	// TODO: https://echo.labstack.com/cookbook/auto-tls/

	server := NewServer()
	server.Initialize()

	e := echo.New()

	e.Use(middleware.Logger())

	e.GET("/request_public_keys", server.getKeys)
	e.POST("/evaluate", server.evaluate)

	e.Static("/static", "./public")

	e.Logger.Fatal(e.Start(HOST + ":" + PORT))
}
