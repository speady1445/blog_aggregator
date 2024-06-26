package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/speady1445/blog_aggregator/internal/database"
	"github.com/speady1445/blog_aggregator/internal/scraper"
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
	mux.HandleFunc("GET /v1/users", apiCfg.middlewareAuth(apiCfg.handlerUsersGet))

	mux.HandleFunc("POST /v1/feeds", apiCfg.middlewareAuth(apiCfg.handlerFeedsCreate))
	mux.HandleFunc("GET /v1/feeds", apiCfg.handlerFeedsGet)

	mux.HandleFunc("POST /v1/feed_follows", apiCfg.middlewareAuth(apiCfg.handlerFeedFollowsCreate))
	mux.HandleFunc("GET /v1/feed_follows/", apiCfg.middlewareAuth(apiCfg.handlerFeedFollowsGet))
	mux.HandleFunc("DELETE /v1/feed_follows/{feedFollowID}", apiCfg.middlewareAuth(apiCfg.handlerFeedFollowsDelete))

	mux.HandleFunc("GET /v1/posts/", apiCfg.middlewareAuth(apiCfg.handlerPostsGet))

	corsMux := middlewareCors(mux)

	server := http.Server{
		Addr:    ":" + port,
		Handler: corsMux,
	}

	const collectionLimit = 2
	const collectionInterval = time.Minute
	go scraper.Start(dbQueries, collectionLimit, collectionInterval)

	log.Println("Listening on port " + port)
	server.ListenAndServe()
}
