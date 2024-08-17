package ffmpeg

import (
	"image"
	"image/color"
	"io"
	"os/exec"
)

func ReadFrames(videoPath string, callback func(*image.RGBA, int) error) error {
	probeOutput, err := ProbeVideo(videoPath)
	if err != nil {
		return err
	}

	frameWidth := probeOutput.Streams[0].Width
	frameHeight := probeOutput.Streams[0].Height

	cmd := exec.Command("ffmpeg", "-i", videoPath, "-f", "rawvideo", "-pix_fmt", "rgb24", "pipe:")
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}
	if err := cmd.Start(); err != nil {
		return err
	}

	frameSize := frameWidth * frameHeight * 3
	frameNumber := 0

	for {
		buf := make([]byte, frameSize)
		_, err = io.ReadFull(stdout, buf)
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		frameNumber++

		img := image.NewRGBA(image.Rect(0, 0, frameWidth, frameHeight))
		for y := 0; y < frameHeight; y++ {
			for x := 0; x < frameWidth; x++ {
				offset := (y*frameWidth + x) * 3
				r := buf[offset]
				g := buf[offset+1]
				b := buf[offset+2]
				img.Set(x, y, color.RGBA{r, g, b, 255})
			}
		}
	
		if err := callback(img, frameNumber); err != nil {
			return err
		}
	}

	if err := cmd.Wait(); err != nil {
        return err
    }

	return nil
}