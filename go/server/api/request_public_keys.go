package api

import (
	"github.com/ensimag-oprf/go/server/controllers"
	"github.com/labstack/echo/v4"
	"net/http"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	router := echo.New()

	oprfServerController := controllers.NewOPRFServerController()
	oprfServerController.Initialize()

	router.GET("/api/request_public_keys", oprfServerController.GetKeysHandler)

	router.ServeHTTP(w, r)
}
