package main

import (
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()

	port, exists := os.LookupEnv("PORT")
	if !exists {
		log.Fatal("PORT environment variable not found")
	}

	mux := http.NewServeMux()

	corsMux := middlewareCors(mux)

	server := http.Server{
		Addr:    ":" + port,
		Handler: corsMux,
	}

	log.Println("Listening on port " + port)
	server.ListenAndServe()
}
