package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()

	if err != nil {
		log.Fatal("Erreur lors du chargement du fichier .env: " + err.Error())
	}

	port := os.Getenv("APP_PORT")
	router := http.NewServeMux()

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), router))
}
