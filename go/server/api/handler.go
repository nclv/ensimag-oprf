package api

import (
	"log"
	"net/http"

	"github.com/ensimag-oprf/go/server/server"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	log.Println("hello")
	router, err := server.NewRouter()
	if err != nil {
		log.Println(err)

		return
	}
	log.Println("router")

	router.ServeHTTP(w, r)
}
