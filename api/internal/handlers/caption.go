package handlers

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"image"
	"image/jpeg"
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

	caption, err := decodeBase64String(urlDecodedCaption)
	if err != nil {
		http.Error(w, "Error decoding the base64-encoded caption", http.StatusBadRequest)
		return
	}
	if len(caption) > 500 {
		http.Error(w, "Caption is too long", http.StatusBadRequest)
		return
	}

	data, err := s.ObjectStorage.GetObject(fmt.Sprintf("frames/%s/%s/large.jpg", key, timestamp))
	if err != nil {
		http.Error(w, "Error fetching the image", http.StatusInternalServerError)
		return
	}

	img, _, err := image.Decode(bytes.NewReader(data))
	if err != nil {
		http.Error(w, "Error decoding the image", http.StatusInternalServerError)
		return
	}

	captionedImage, err := drawCaption(img, caption)
	if err != nil {
		http.Error(w, "Error drawing the caption", http.StatusInternalServerError)
		return
	}

	// Write the response
	w.Header().Set("Content-Type", "image/jpeg")
	err = jpeg.Encode(w, captionedImage, nil)
	if err != nil {
		http.Error(w, "Error encoding the image", http.StatusInternalServerError)
		return
	}
}

func decodeBase64String(base64String string) (string, error) {
	// Ensure the base64 string is properly padded
	if len(base64String)%4 != 0 {
		base64String += strings.Repeat("=", 4-len(base64String)%4)
	}

	captionBytes, err := base64.StdEncoding.DecodeString(base64String)
	if err != nil {
		return "", err
	}
	return string(captionBytes), nil
}

func drawCaption(img image.Image, caption string) (image.Image, error) {
	// Prepare drawing context
	dc := gg.NewContextForImage(img)
	imgWidth := float64(dc.Width())
	imgHeight := float64(dc.Height())
	imgPadding := 8.0

	// Draw each line of text
	fontColor := [3]float64{0.97254, 0.89803, 0.67843}
	fontSize := float64(32)
	if err := dc.LoadFontFace("static/fonts/FFBlurProMedium/font.ttf", fontSize); err != nil {
		return nil, err
	}

	lines := splitTextIntoLines(dc, caption, imgWidth-imgPadding)

	lineHeight := fontSize * 1.0
	startY := imgHeight - float64(len(lines))*lineHeight - 8 // pixels padding from the bottom

	for i, line := range lines {
		y := startY + float64(i)*lineHeight
		fmt.Printf("Drawing line %d at y=%f\n", i, y)

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

	return dc.Image(), nil
}

// set font size before calling this
func splitTextIntoLines(dc *gg.Context, text string, maxWidth float64) []string {
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
