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
            JOIN episodes e on f.episode = e.key
        WHERE s.episode = $1
            AND f.id = (
                SELECT MIN(f2.id)
                FROM frames f2
                WHERE f2.episode = s.episode
                AND f2.timestamp BETWEEN s.start_timestamp AND s.end_timestamp
            )
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
            JOIN subtitles s on f.episode = s.episode
        WHERE f.timestamp = $1
            AND f.episode = $2
            AND f.timestamp BETWEEN s.start_timestamp AND s.end_timestamp;
    `

    var frameResponse models.FrameResponse
    var frameData models.Frame
    var subtitleData models.Subtitle
    var episodeData models.Episode

    err := s.DB.QueryRow(query, timestamp, episode).Scan(
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