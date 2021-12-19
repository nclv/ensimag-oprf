package server

import (
	"html/template"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"github.com/oprf/go/server/controllers"
)

func NewRouter() (*echo.Echo, error) {
	router := echo.New()

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
