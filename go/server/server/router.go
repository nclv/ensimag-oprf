package server

import (
	"github.com/ensimag-oprf/go/server/controllers"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func NewRouter() (*echo.Echo, error) {
	router := echo.New()

	// Middlewares
	router.Use(middleware.Logger())
	router.Use(middleware.Recover())
	router.Use(middleware.Gzip())
	router.Use(middleware.CORS())

	// Endpoints
	// router.GET("/", controllers.IndexHandler)
	router.File("/", "public/index.html")

	oprfServerController := controllers.NewOPRFServerController()
	oprfServerController.Initialize()

	router.GET("/api/request_public_keys", oprfServerController.GetKeysHandler)
	router.POST("/api/evaluate", oprfServerController.EvaluateHandler)

	// Static files
	router.Static("/static", "./public/static")

	return router, nil
}
