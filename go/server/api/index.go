package api

import (
	"github.com/ensimag-oprf/go/server/routers"
	"log"
	"net/http"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	router, err := routers.NewRouter()
	if err != nil {
		log.Println(err)

		return
	}

	router.ServeHTTP(w, r)
}
