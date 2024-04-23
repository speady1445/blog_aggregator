package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/speady1445/blog_aggregator/internal/database"
)

type apiConfig struct {
	DB *database.Queries
}

func main() {
	godotenv.Load()

	port, exists := os.LookupEnv("PORT")
	if !exists {
		log.Fatal("PORT environment variable not found")
	}
	dbURL, exists := os.LookupEnv("PSQL_URL")
	if !exists {
		log.Fatal("PSQL_URL environment variable not found")
	}

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal(err)
	}
	dbQueries := database.New(db)

	apiCfg := apiConfig{
		DB: dbQueries,
	}

	mux := http.NewServeMux()

	mux.HandleFunc("GET /v1/healthz", handlerHealthz)

	mux.HandleFunc("POST /v1/users", apiCfg.handlerUsersCreate)
	mux.HandleFunc("GET /v1/users", apiCfg.handlerUsersGet)

	corsMux := middlewareCors(mux)

	server := http.Server{
		Addr:    ":" + port,
		Handler: corsMux,
	}

	log.Println("Listening on port " + port)
	server.ListenAndServe()
}
