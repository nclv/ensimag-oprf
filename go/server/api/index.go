package api

import (
	"github.com/ensimag-oprf/go/server/server"
	"log"
	"net/http"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	router, err := server.NewRouter()
	if err != nil {
		log.Println(err)

		return
	}

	// w.Header().Set("Content-Type", "application/json")

	router.ServeHTTP(w, r)
}
