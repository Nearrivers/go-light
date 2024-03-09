package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/nearrivers/go-light/lights"
)

func main() {
	err := godotenv.Load()

	if err != nil {
		log.Fatal("Erreur lors du chargement du fichier .env: " + err.Error())
	}

	port := os.Getenv("APP_PORT")
	router := http.NewServeMux()

	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		devices, err := lights.GetDeviceList()

		if len(err) > 0 {
			w.WriteHeader(500)
			w.Write([]byte(err[0].Error()))
		}

		fmt.Println(devices)
	})

	devices, errors := lights.GetDeviceList()
	if len(errors) > 0 {
		panic(errors[0].Error())
	}

	fmt.Println(devices)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), router))
}
