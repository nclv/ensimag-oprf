package api

import (
	"log"
	"net/http"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	router, err := NewRouter()
	if err != nil {
		log.Println(err)

		return
	}

	router.ServeHTTP(w, r)
}
