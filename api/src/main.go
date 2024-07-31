package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func searchFrames(w http.ResponseWriter, r *http.Request) {
	// Get the query string
	query := r.URL.Query().Get("q")

	fmt.Println("Query: ", query)
}

func main() {
	router := mux.NewRouter()

	// Routes
	router.HandleFunc("/search", searchFrames).Methods("GET")

	// Start the server
	log.Println("Server is running on port 8000")
	log.Fatal(http.ListenAndServe(":8000", router))
}