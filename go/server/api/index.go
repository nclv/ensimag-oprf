package api

import (
	"log"
	"net/http"

	"github.com/ensimag-oprf/go/server/controllers"
	"github.com/ensimag-oprf/go/server/routers"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	serializedBase64KeyMap := controllers.LoadPrivateKeysFromEnv()

	router, err := routers.NewRouter(serializedBase64KeyMap)
	if err != nil {
		log.Println(err)

		return
	}

	router.ServeHTTP(w, r)
}
