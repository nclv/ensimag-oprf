package server

import (
	"github.com/ensimag-oprf/go/server/controllers"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func NewRouter() (*echo.Echo, error) {
	router := echo.New()

	// Template renderer
	//renderer := &Template{
	//	templates: template.Must(template.ParseGlob("public/*.html")), // Path issue with vercel
	//}
	//router.Renderer = renderer

	// Middlewares
	router.Use(middleware.Logger())
	router.Use(middleware.Recover())
	router.Use(middleware.Gzip())
	router.Use(middleware.CORS())

	// Endpoints
	// router.GET("/", controllers.IndexHandler)
	// router.File("/", "public/main.html")

	oprfServerController := controllers.NewOPRFServerController()
	oprfServerController.Initialize()

	router.GET("/api/request_public_keys", oprfServerController.GetKeysHandler)
	router.POST("/api/evaluate", oprfServerController.EvaluateHandler)

	// Static files
	// router.Static("/static", "./static")

	return router, nil
}
