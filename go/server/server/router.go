package server

import (
	"fmt"
	"html/template"
	"log"
	"os"

	"github.com/ensimag-oprf/go/server/controllers"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func NewRouter() (*echo.Echo, error) {
	router := echo.New()

	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(dir)

	// Template renderer
	renderer := &Template{
		templates: template.Must(template.ParseGlob("public/*.html")),
	}
	router.Renderer = renderer

	// Middlewares
	router.Use(middleware.Logger())
	router.Use(middleware.Recover())
	router.Use(middleware.Gzip())
	router.Use(middleware.CORS())

	// Endpoints
	router.GET("/", controllers.IndexHandler)

	oprfServerController := controllers.NewOPRFServerController()
	oprfServerController.Initialize()

	router.GET("/request_public_keys", oprfServerController.GetKeysHandler)
	router.POST("/evaluate", oprfServerController.EvaluateHandler)

	// Static files
	router.Static("/static", "./static")

	return router, nil
}
