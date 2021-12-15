package main

import (
	"html/template"
	"io"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

const (
	HOST = "localhost"
	PORT = "1323"
)

func indexHandler(c echo.Context) error {
	return c.Render(http.StatusOK, "index.html", nil) //nolint:wrapcheck
}

type Template struct {
	templates *template.Template
}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data) //nolint:wrapcheck
}

func main() {
	// TODO: https://echo.labstack.com/cookbook/auto-tls/

	serverManager := NewOPRFServerManager()
	serverManager.Initialize()

	e := echo.New()

	// Template renderer
	renderer := &Template{
		templates: template.Must(template.ParseGlob("templates/*.html")),
	}
	e.Renderer = renderer

	// Middlewares
	e.Use(middleware.Logger())

	// Endpoints
	e.GET("/", indexHandler)
	e.GET("/request_public_keys", serverManager.getKeysHandler)
	e.POST("/evaluate", serverManager.evaluateHandler)

	// Static files
	e.Static("/static", "./public")

	e.Logger.Fatal(e.Start(HOST + ":" + PORT))
}
