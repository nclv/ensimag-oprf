package api

import (
	"log"
	"net/http"

	"github.com/ensimag-oprf/go/server/server"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	router, err := server.NewRouter()
	if err != nil {
		log.Println(err)

		return
	}

	router.ServeHTTP(w, r)
}
