package api

import (
	"html/template"
	"io"
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/oprf/server/controllers"
)

type Template struct {
	templates *template.Template
}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data) //nolint:wrapcheck
}

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

func Handler(w http.ResponseWriter, r *http.Request) {
	router, err := NewRouter()
	if err != nil {
		log.Println(err)

		return
	}

	router.ServeHTTP(w, r)
}
