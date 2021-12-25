package routers

import (
	"github.com/ensimag-oprf/go/server/controllers"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func NewRouter(serializedBase64KeyMap controllers.SerializedBase64KeyMap) (*echo.Echo, error) {
	router := echo.New()

	// Middlewares
	router.Use(middleware.Logger())
	router.Use(middleware.Recover())
	router.Use(middleware.Gzip())
	router.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"https://ensimag-oprf.vercel.app"},
		AllowMethods: []string{http.MethodGet, http.MethodPut, http.MethodPost, http.MethodDelete},
	}))

	// Endpoints
	router.File("/", "public/index.html")

	oprfServerController := controllers.NewOPRFServerController()
	if err := oprfServerController.Initialize(serializedBase64KeyMap); err != nil {
		return nil, err
	}

	router.GET("/api/request_public_keys", oprfServerController.GetKeysHandler)
	router.POST("/api/evaluate", oprfServerController.EvaluateHandler)

	// Static files
	router.Static("/static", "./public/static")

	return router, nil
}
