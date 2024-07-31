package main

import (
	"bytes"
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/elastic/go-elasticsearch/v7"
	"github.com/gorilla/mux"
)

var esClient *elasticsearch.Client

var cfg = elasticsearch.Config{
	Addresses: []string{
	  "http://elasticsearch:9200",
	},
	Username: "elastic",
	Password: "elastic",
  }

func searchFrames(w http.ResponseWriter, r *http.Request) {
    // Get the query string
    query := r.URL.Query().Get("q")

    // Build the search request body
    searchBody := map[string]interface{}{
        "_source": []string{"episode", "subtitle", "timestamp"},
        "query": map[string]interface{}{
            "match": map[string]interface{}{
                "subtitle": map[string]interface{}{
                    "query": "*" + query + "*",
                },
            },
        },
    }

    // Encode the search body to JSON
    var buf bytes.Buffer
    if err := json.NewEncoder(&buf).Encode(searchBody); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    // Execute the search query
    res, err := esClient.Search(
        esClient.Search.WithContext(context.Background()),
        esClient.Search.WithIndex("frames"),
        esClient.Search.WithBody(&buf),
        esClient.Search.WithTrackTotalHits(true),
    )
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    defer res.Body.Close()

    // Decode the search results
    var searchResult map[string]interface{}
    if err := json.NewDecoder(res.Body).Decode(&searchResult); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    // Extract the hits.hits attribute
    hits, ok := searchResult["hits"].(map[string]interface{})["hits"].([]interface{})
    if !ok {
        http.Error(w, "Failed to extract hits", http.StatusInternalServerError)
        return
    }

    // Map over the hits to extract only the _source attribute
    hitSources := make([]interface{}, len(hits))
    for i, hit := range hits {
        hitMap, ok := hit.(map[string]interface{})
        if !ok {
            http.Error(w, "Failed to parse hit", http.StatusInternalServerError)
            return
        }
        hitSources[i] = hitMap["_source"]
    }

	// Encode the search results to JSON and write to response
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(hitSources); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func main() {
	var err error
	esClient, err = elasticsearch.NewClient(cfg)
    if err != nil {
        log.Fatalf("Error creating the client: %s", err)
    }

	router := mux.NewRouter()

	// Routes
	router.HandleFunc("/search", searchFrames).Methods("GET")

	// Start the server
	log.Println("Server is running on port 8000")
	log.Fatal(http.ListenAndServe(":8000", router))
}