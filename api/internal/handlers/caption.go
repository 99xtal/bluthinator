package handlers

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

func (s *Server) GetCaptionedFrame(w http.ResponseWriter, r *http.Request) {
	// Get the parameters
	params := mux.Vars(r)
	key := params["key"]
	timestamp := params["timestamp"]

	data, err := s.ObjectStorage.GetObject(fmt.Sprintf("bluthinator/frames/%s/%s/medium.png", key, timestamp))
	if err != nil {
		http.Error(w, "Error fetching the image", http.StatusInternalServerError)
		return
	}

	// Write the response
	w.Header().Set("Content-Type", "image/png")
	w.Write(data)
}
