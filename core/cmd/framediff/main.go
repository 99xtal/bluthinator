package main

import (
	"fmt"
	"image"
	"image/jpeg"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/99xtal/bluthinator/core/internal/ffmpeg"
	"github.com/99xtal/bluthinator/core/internal/ssim"
	"github.com/nfnt/resize"
)

var (
	similarityThreshold = 0.70
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: framediff <video_file>")
		os.Exit(1)
	}
	videoFilePath := os.Args[1]

	start := time.Now()

	err := extractFrames(videoFilePath)
	if err != nil {
		log.Fatal(err)
	}

	elapsed := time.Since(start)
	fmt.Printf("Time taken: %s\n", elapsed)
}

func extractFrames(videoPath string) error {
	base := filepath.Base(videoPath)
	key := strings.TrimSuffix(base, filepath.Ext(base))

	outputDir := fmt.Sprintf("frames/%s", key)
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return err
	}

	var significantFrame *image.RGBA
	
	err := ffmpeg.ReadFrames(videoPath, func(img *image.RGBA, frameNumber int) error {
		if frameNumber == 1 {
			significantFrame = img
			return nil
		}

		mean_ssim := ssim.MeanSSIM(significantFrame, img)
		if (mean_ssim < similarityThreshold) {
			significantFrame = img

			size := map[string]uint{
				"small": 240,
				"medium": 480,
			}
			for sizeName, imgWidth := range size {
				resizedImg := resize.Resize(imgWidth, 0, img, resize.Lanczos3)
				filePath := fmt.Sprintf("%s/%d/%s.jpg", outputDir, frameNumberToMs(frameNumber, 24), sizeName)

				err := saveAsJPEG(resizedImg, filePath)
				if err != nil {
					return err
				}
			}
		}
		
		return nil
	})
	if err != nil {
		return err
	}

	return nil
}

func frameNumberToMs(frameNumber int, fps int) int {
	return frameNumber * 1000 / fps
}

func saveAsJPEG(img image.Image, fileName string) error {
	// Ensure the directory exists
	dir := filepath.Dir(fileName)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	file, err := os.Create(fileName)
	if err != nil {
		return err
	}

	return jpeg.Encode(file, img, nil)
}