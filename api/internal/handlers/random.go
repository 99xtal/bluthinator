package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/99xtal/bluthinator/api/internal/models"
)

func (s *Server) GetRandomFrame(w http.ResponseWriter, r *http.Request) {
	query := `
		SELECT
			id,
			episode,
			timestamp
		FROM frames
		TABLESAMPLE SYSTEM (1)
		ORDER BY RANDOM()
		LIMIT 1;
	`

	var frame models.Frame
	err := s.DB.QueryRow(query).Scan(&frame.ID, &frame.Episode, &frame.Timestamp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(frame); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
