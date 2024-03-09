package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
	"github.com/nearrivers/go-light/json"
	"github.com/nearrivers/go-light/lights"
)

func LightRoutes() *chi.Mux {
	router := chi.NewRouter()
	router.Get("/light", func(w http.ResponseWriter, r *http.Request) {
		bulbsState, err := lights.GetBulbsState()
		if err != nil {
			json.RespondWithError(w, 500, err.Error())
			return
		}

		json.RespondWithJSON(w, 200, bulbsState)
	})

	router.Put("/light", func(w http.ResponseWriter, r *http.Request) {
		err := lights.ControlLights()
		if err != nil {
			json.RespondWithError(w, 500, err.Error())
		}

		json.RespondWithJSON(w, 200, "")
	})

	return router
}

func main() {
	err := godotenv.Load()

	if err != nil {
		log.Fatal("Erreur lors du chargement du fichier .env: " + err.Error())
	}

	port := os.Getenv("APP_PORT")
	router := LightRoutes()

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), router))
}
