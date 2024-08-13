package main

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/elastic/go-elasticsearch/v7"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"github.com/rs/cors"

	"github.com/99xtal/bluthinator/api/internal/config"
	"github.com/99xtal/bluthinator/api/internal/handlers"
)

func main() {
	config := config.New()

	// Initialize the Elasticsearch client
	var err error
	esClient, err := elasticsearch.NewClient(config.GetElasticSearchConfig())
	if err != nil {
		log.Fatalf("Error creating the client: %s", err)
	}

	// Initialize the PostgreSQL database connection
	connStr := config.GetPostgresConnString()
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("Error connecting to the database: %s", err)
	}
	defer db.Close()

	server := handlers.NewServer(db, esClient)

	router := mux.NewRouter()

	// Routes
	router.HandleFunc("/episode/{key}", server.GetEpisodeData).Methods("GET")
	router.HandleFunc("/episode/{key}/{timestamp}", server.GetEpisodeFrame).Methods("GET")
	router.HandleFunc("/nearby", server.GetNearbyFrames).Methods("GET")
	router.HandleFunc("/search", server.SearchFrames).Methods("GET")
	router.HandleFunc("/healthcheck", server.HealthCheck).Methods("GET")

	// Configure CORS
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000"},
		AllowCredentials: true,
	})
	handler := c.Handler(router)

	// Start the server
	port := ":" + config.ServerPort
	log.Println("Server listening on port", port)
	log.Fatal(http.ListenAndServe(port, handler))
}
