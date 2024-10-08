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
	"github.com/99xtal/bluthinator/api/internal/services"
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

	// Initialize object storage client
	storageClient := services.NewStorageClient(config.ObjectStorageEndpoint)

	server := handlers.NewServer(db, esClient, storageClient)

	router := mux.NewRouter()

	// Routes
	router.HandleFunc("/caption/{key}/{timestamp}", server.GetCaptionedFrame).Methods("GET")
	router.HandleFunc("/episode/{key}", server.GetEpisodeData).Methods("GET")
	router.HandleFunc("/episode/{key}/{timestamp}", server.GetEpisodeFrame).Methods("GET")
	router.HandleFunc("/healthcheck", server.HealthCheck).Methods("GET")
	router.HandleFunc("/nearby", server.GetNearbyFrames).Methods("GET")
	router.HandleFunc("/random", server.GetRandomFrame).Methods("GET")
	router.HandleFunc("/search", server.SearchFrames).Methods("GET")

	// Configure CORS
	c := cors.New(cors.Options{
		AllowedOrigins:   config.AllowedOrigins,
		AllowCredentials: true,
	})
	handler := c.Handler(router)

	// Start the server
	port := ":" + config.ServerPort
	log.Println("Server listening on port", port)
	log.Fatal(http.ListenAndServe(port, handler))
}
