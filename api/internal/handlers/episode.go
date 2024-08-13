package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/99xtal/bluthinator/api/internal/models"
)

func (s *Server) GetEpisodeData(w http.ResponseWriter, r *http.Request) {
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
            JOIN episodes e ON f.episode = e.key
        WHERE s.episode = $1
            AND f.id = (
                SELECT MIN(f2.id)
                FROM frames f2
                WHERE f2.episode = s.episode
                AND f2.timestamp BETWEEN s.start_timestamp AND s.end_timestamp
                ORDER BY ABS(f2.timestamp - ((s.start_timestamp + s.end_timestamp) / 2))
            );
    `

	rows, err := s.DB.Query(query, episode)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var episodeData models.EpisodeResponse
	var subtitles []models.Subtitle

	for rows.Next() {
		var subtitle models.Subtitle
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

func (s *Server) GetEpisodeFrame(w http.ResponseWriter, r *http.Request) {
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
            LEFT JOIN subtitles s ON f.episode = s.episode
                AND f.timestamp BETWEEN s.start_timestamp AND s.end_timestamp
        WHERE f.timestamp = $1
            AND f.episode = $2
    `

	var frameResponse models.FrameResponse
	var frameData models.Frame
	var subtitleData models.Subtitle
	var episodeData models.Episode

	var subtitleID sql.NullInt64
	var subtitleEpisode sql.NullString
	var subtitleText sql.NullString
	var subtitleStartTimestamp sql.NullInt64
	var subtitleEndTimestamp sql.NullInt64

	err := s.DB.QueryRow(query, timestamp, episode).Scan(
		&frameData.ID,
		&frameData.Timestamp,
		&frameData.Episode,
		&subtitleID,
		&subtitleEpisode,
		&subtitleText,
		&subtitleStartTimestamp,
		&subtitleEndTimestamp,
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
			return
		}
	}

	// Check if all subtitle fields are NULL
	if !subtitleID.Valid && !subtitleEpisode.Valid && !subtitleText.Valid && !subtitleStartTimestamp.Valid && !subtitleEndTimestamp.Valid {
		frameResponse.Subtitle = nil
	} else {
		// Handle NULL subtitle fields
		if subtitleID.Valid {
			subtitleData.ID = int(subtitleID.Int64)
		}
		if subtitleEpisode.Valid {
			subtitleData.Episode = subtitleEpisode.String
		}
		if subtitleText.Valid {
			subtitleData.Text = subtitleText.String
		}
		if subtitleStartTimestamp.Valid {
			subtitleData.StartTimestamp = int(subtitleStartTimestamp.Int64)
		}
		if subtitleEndTimestamp.Valid {
			subtitleData.EndTimestamp = int(subtitleEndTimestamp.Int64)
		}
		frameResponse.Subtitle = &subtitleData
	}

	// Populate frameResponse with the data
	frameResponse.Frame = frameData
	frameResponse.Episode = episodeData

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(frameResponse); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
