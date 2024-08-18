package ffmpeg

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math"
	"os/exec"
	"strconv"
	"strings"
)

type FFProbeOutput struct {
	Format 	Format 		`json:"format"`
    Streams []Stream 	`json:"streams"`
}

type Format struct {
	Duration string `json:"duration"`
}

type Stream struct {
    Width  int `json:"width"`
    Height int `json:"height"`
	RFrameRate string `json:"r_frame_rate"`
}

func ProbeVideo(filePath string) (FFProbeOutput, error) {
	var probeOutput FFProbeOutput

	probeCmd := exec.Command("ffprobe", "-v", "error", "-select_streams", "v:0", "-show_entries", "stream=width,height,r_frame_rate", "-show_entries", "format=duration", "-of", "json", filePath)
	var out bytes.Buffer
	probeCmd.Stdout = &out
	if err := probeCmd.Run(); err != nil {
		return FFProbeOutput{}, err
	}

	if err := json.Unmarshal(out.Bytes(), &probeOutput); err != nil {
		return FFProbeOutput{}, err
	}

	if len(probeOutput.Streams) == 0 {
		return FFProbeOutput{}, fmt.Errorf("no video stream found in %s", filePath)
	}

	return probeOutput, nil
}

func (f *FFProbeOutput) DurationMS() (float64, error) {
	durationSec, err := strconv.ParseFloat(f.Format.Duration, 64)
	if err != nil {
		return 0, err
	}

	durationMS := durationSec * 1000
	return durationMS, nil
}

func (f *FFProbeOutput) FrameRate() (int, error) {
	parts := strings.Split(f.Streams[0].RFrameRate, "/")
	numerator, err := strconv.ParseFloat(parts[0], 64)
	if err != nil {
		return 0, err
	}

	denominator, err := strconv.ParseFloat(parts[1], 64)
	if err != nil {
		return 0, err
	}

	frameRate := int(math.Ceil(numerator / denominator))
	return frameRate, nil
}

func (f *FFProbeOutput) TotalFrames() (int, error) {
	durationMS, err := f.DurationMS()
	if err != nil {
		return 0, err
	}

	frameRate, err := f.FrameRate()
	if err != nil {
		return 0, err
	}

	totalFrames := int(durationMS) * frameRate / 1000
	return totalFrames, nil
}