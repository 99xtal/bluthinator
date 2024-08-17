package ffmpeg

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os/exec"
)

type FFProbeOutput struct {
    Streams []Stream `json:"streams"`
}

type Stream struct {
    Width  int `json:"width"`
    Height int `json:"height"`
}

func ProbeVideo(filePath string) (FFProbeOutput, error) {
	var probeOutput FFProbeOutput

	probeCmd := exec.Command("ffprobe", "-v", "error", "-select_streams", "v:0", "-show_entries", "stream=width,height", "-of", "json", filePath)
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