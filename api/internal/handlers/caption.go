package handlers

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"image"
	"image/png"
	"net/http"
	"net/url"
	"strings"

	"github.com/fogleman/gg"
	"github.com/gorilla/mux"
)

func (s *Server) GetCaptionedFrame(w http.ResponseWriter, r *http.Request) {
	// Get the parameters
	params := mux.Vars(r)
	key := params["key"]
	timestamp := params["timestamp"]

	base64Caption := r.URL.Query().Get("b")

	// URL decode the base64-encoded caption
	urlDecodedCaption, err := url.QueryUnescape(base64Caption)
	if err != nil {
		http.Error(w, "Error decoding the URL-encoded caption", http.StatusBadRequest)
		return
	}

	// Ensure the base64 string is properly padded
	paddedCaption := urlDecodedCaption
	if len(paddedCaption)%4 != 0 {
		paddedCaption += strings.Repeat("=", 4-len(paddedCaption)%4)
	}

	// Decode the base64-encoded caption
	captionBytes, err := base64.StdEncoding.DecodeString(paddedCaption)
	if err != nil {
		http.Error(w, "Error decoding the caption", http.StatusBadRequest)
		return
	}
	caption := string(captionBytes)

	data, err := s.ObjectStorage.GetObject(fmt.Sprintf("bluthinator/frames/%s/%s/medium.png", key, timestamp))
	if err != nil {
		http.Error(w, "Error fetching the image", http.StatusInternalServerError)
		return
	}

	img, _, err := image.Decode(bytes.NewReader(data))
	if err != nil {
		http.Error(w, "Error decoding the image", http.StatusInternalServerError)
		return
	}

	// Prepare drawing context
	dc := gg.NewContextForImage(img)
	imgWidth := float64(dc.Width())
	imgHeight := float64(dc.Height())

	fontColor := [3]float64{0.97254, 0.89803, 0.67843}
	fontSize := float64(32)
	if err := dc.LoadFontFace("static/fonts/DeFonteReducedNormal.ttf", fontSize); err != nil {
		panic(err)
	}

	splitTextIntoLines := func(text string, maxWidth float64) []string {
		words := strings.Fields(text)
		var lines []string
		var currentLine string
		for _, word := range words {
			testLine := currentLine + " " + word
			testWidth, _ := dc.MeasureString(testLine)
			if testWidth > maxWidth {
				lines = append(lines, strings.TrimSpace(currentLine))
				currentLine = word
			} else {
				currentLine = testLine
			}
		}
		lines = append(lines, strings.TrimSpace(currentLine))
		return lines
	}

	lines := splitTextIntoLines(caption, imgWidth-20) // 20 pixels padding

	// Draw each line of text
	lineHeight := fontSize * 1.2
	startY := imgHeight - float64(len(lines))*lineHeight - 48 // 10 pixels padding from the bottom

	for i, line := range lines {
		y := startY + float64(i)*lineHeight

		// Draw the shadow with reduced opacity
		shadowColor := [3]float64{0, 0, 0} // Black color for the shadow
		shadowOpacity := 0.5               // 50% opacity
		shadowOffsetX := 2.0
		shadowOffsetY := 2.0

		dc.SetRGBA(shadowColor[0], shadowColor[1], shadowColor[2], shadowOpacity)
		dc.DrawStringAnchored(line, imgWidth/2+shadowOffsetX, y+shadowOffsetY, 0.5, 0.5)

		// Draw the main text
		dc.SetRGB(fontColor[0], fontColor[1], fontColor[2])
		dc.DrawStringAnchored(line, imgWidth/2, y, 0.5, 0.5)
	}

	// Write the response
	w.Header().Set("Content-Type", "image/png")
	err = png.Encode(w, dc.Image())
	if err != nil {
		http.Error(w, "Error encoding the image", http.StatusInternalServerError)
		return
	}
}
