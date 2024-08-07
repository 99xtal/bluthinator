package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/99xtal/bluthinator/api/internal/models"
)

func (s *Server) GetNearbyFrames(w http.ResponseWriter, r *http.Request) {
	episode := r.URL.Query().Get("e")
	timestamp := r.URL.Query().Get("t")

	if (episode == "" || timestamp == "") {
		http.Error(w, "Missing query parameters", http.StatusBadRequest)
		return
	}

    query := `
        (SELECT id, episode, timestamp FROM frames
         WHERE episode = $1 AND timestamp < $2
         ORDER BY timestamp DESC
         LIMIT 3)
        UNION
        (SELECT id, episode, timestamp FROM frames
         WHERE episode = $1 AND timestamp >= $2
         ORDER BY timestamp ASC
         LIMIT 4)
        ORDER BY timestamp;
    `

    rows, err := s.DB.Query(query, episode, timestamp)
    if err != nil {
        http.Error(w, "Database query failed", http.StatusInternalServerError)
        return
    }
    defer rows.Close()

	var nearbyFrames []models.Frame
	for rows.Next() {
		var frame models.Frame
		if err := rows.Scan(&frame.ID, &frame.Episode, &frame.Timestamp); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		nearbyFrames = append(nearbyFrames, frame)
	}

	if err := rows.Err(); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

	w.Header().Set("Content-Type", "application/json")
    if err := json.NewEncoder(w).Encode(nearbyFrames); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
}