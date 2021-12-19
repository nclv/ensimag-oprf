package api

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"net/http"
)

func hello(c echo.Context) error {
	return c.String(http.StatusOK, "Hello World")
}

func hello2(c echo.Context) error {
	return c.String(http.StatusOK, "Hello World2")
}

func Handler(w http.ResponseWriter, r *http.Request) {
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.GET("/api/", hello)
	e.GET("/api/2", hello2)

	e.ServeHTTP(w, r)
}

//
//func Handler(w http.ResponseWriter, r *http.Request) {
//	router, err := server.NewRouter()
//	if err != nil {
//		log.Println(err)
//
//		return
//	}
//
//	router.ServeHTTP(w, r)
//}
