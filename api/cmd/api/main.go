package main

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/elastic/go-elasticsearch/v7"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)


type Frame struct {
    ID        int    `json:"id"`
    Episode   string `json:"episode"`
    Timestamp int    `json:"timestamp"`
}

type Subtitle struct {
    ID              int    `json:"id"`
    Episode         string `json:"episode"`
    Text            string `json:"text"`
    StartTimestamp  int    `json:"start_timestamp"`
    EndTimestamp    int    `json:"end_timestamp"`
    FrameTimestamp  int    `json:"frame_timestamp"`
}

type Episode struct {
    EpisodeNumber    int        `json:"episode_number"`
    Season           int        `json:"season"`
    Title            string     `json:"title"`
    Director         string     `json:"director"`
    Subtitles        []Subtitle `json:"subtitles"`
}

type FrameResponse struct {
    Frame   Frame     `json:"frame"`
    Episode Episode   `json:"episode"`
    Subtitle Subtitle `json:"subtitle"`
}

var db *sql.DB

var esClient *elasticsearch.Client

var cfg = elasticsearch.Config{
	Addresses: []string{
	  "http://elasticsearch:9200",
	},
	Username: os.Getenv("ELASTIC_USERNAME"),
	Password: os.Getenv("ELASTIC_PASSWORD"),
  }


func episodeHandler(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    episode := vars["key"]

    query := `
        SELECT 
            e.episode_number,
            e.season,
            e.title,
            e.director,
            s.id,
            s.episode,
            s.text,
            s.start_timestamp,
            s.end_timestamp,
            f.timestamp
        FROM subtitles s
            JOIN frames f ON s.episode = f.episode
                AND f.timestamp BETWEEN s.start_timestamp AND s.end_timestamp
            JOIN episodes e on f.episode = e.key
        WHERE s.episode = $1
            AND f.id = (
                SELECT MIN(f2.id)
                FROM frames f2
                WHERE f2.episode = s.episode
                AND f2.timestamp BETWEEN s.start_timestamp AND s.end_timestamp
            )
    `

    rows, err := db.Query(query, episode)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    defer rows.Close()

    var episodeData Episode
    var subtitles []Subtitle

    for rows.Next() {
        var subtitle Subtitle
        if err := rows.Scan(&episodeData.EpisodeNumber, &episodeData.Season, &episodeData.Title, &episodeData.Director, &subtitle.ID, &subtitle.Episode, &subtitle.Text, &subtitle.StartTimestamp, &subtitle.EndTimestamp, &subtitle.FrameTimestamp); err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }
        subtitles = append(subtitles, subtitle)
    }

    if err := rows.Err(); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    if len(subtitles) == 0 {
        http.Error(w, "Episode not found", http.StatusNotFound)
        return
    }

    episodeData.Subtitles = subtitles

    w.Header().Set("Content-Type", "application/json")
    if err := json.NewEncoder(w).Encode(episodeData); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
}

func frameHandler(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    episode := vars["key"]
    timestamp := vars["timestamp"]

    query := `
        SELECT
            f.id,
            f.timestamp,
            f.episode,
            s.id,
            s.episode,
            s.text,
            s.start_timestamp,
            s.end_timestamp,
            e.episode_number,
            e.season,
            e.title,
            e.director
        FROM frames f
            JOIN episodes e ON e.key = f.episode
            JOIN subtitles s on f.episode = s.episode
        WHERE f.timestamp = $1
            AND f.episode = $2
            AND f.timestamp BETWEEN s.start_timestamp AND s.end_timestamp;
    `

    var frameResponse FrameResponse
    var frameData Frame
    var subtitleData Subtitle
    var episodeData Episode

    err := db.QueryRow(query, timestamp, episode).Scan(
        &frameData.ID,
        &frameData.Timestamp,
        &frameData.Episode,
        &subtitleData.ID,
        &subtitleData.Episode,
        &subtitleData.Text,
        &subtitleData.StartTimestamp,
        &subtitleData.EndTimestamp,
        &episodeData.EpisodeNumber,
        &episodeData.Season,
        &episodeData.Title,
        &episodeData.Director,
    )
    if err != nil {
        if err == sql.ErrNoRows {
            http.Error(w, "No matching records found", http.StatusNotFound)
        } else {
            http.Error(w, err.Error(), http.StatusInternalServerError)
        }
        return
    }

    frameResponse.Frame = frameData
    frameResponse.Episode = episodeData
    frameResponse.Subtitle = subtitleData

    w.Header().Set("Content-Type", "application/json")
    if err := json.NewEncoder(w).Encode(frameResponse); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
}

func searchHandler(w http.ResponseWriter, r *http.Request) {
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
    // Initialize the Elasticsearch client
	var err error
	esClient, err = elasticsearch.NewClient(cfg)
    if err != nil {
        log.Fatalf("Error creating the client: %s", err)
    }

    // Initialize the PostgreSQL database connection
    connStr := "host=" + os.Getenv("POSTGRES_HOST") + " user=" + os.Getenv("POSTGRES_USER") + " dbname=" + os.Getenv("POSTGRES_DB") + " sslmode=disable password=" + os.Getenv("POSTGRES_PASSWORD") + " port=" + os.Getenv("POSTGRES_PORT")
    fmt.Println(connStr)
    db, err = sql.Open("postgres", connStr)
    if err != nil {
        log.Fatalf("Error connecting to the database: %s", err)
    }
    defer db.Close()

	router := mux.NewRouter()

	// Routes
    router.HandleFunc("/episode/{key}", episodeHandler).Methods("GET")
    router.HandleFunc("/episode/{key}/{timestamp}", frameHandler).Methods("GET")
	router.HandleFunc("/search", searchHandler).Methods("GET")

	// Start the server
    port := ":8000"
	log.Println("Server listening on port", port)
	log.Fatal(http.ListenAndServe(port, router))
}