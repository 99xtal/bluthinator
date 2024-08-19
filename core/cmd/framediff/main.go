package main

import (
	"fmt"
	"image"
	"image/jpeg"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/99xtal/bluthinator/core/internal/ffmpeg"
	"github.com/99xtal/bluthinator/core/internal/ssim"
	"github.com/nfnt/resize"
	"github.com/vbauerster/mpb/v7"
	"github.com/vbauerster/mpb/v7/decor"
)

var (
	similarityThreshold = 0.70
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: framediff <video_dir>")
		os.Exit(1)
	}
	inputDirPath := os.Args[1]

	videoFiles, err := filepath.Glob(filepath.Join(inputDirPath, "*.mkv"))
	if err != nil {
		log.Fatal(err)
	}
	if len(videoFiles) == 0 {
		fmt.Println("No video files found in the input directory")
		os.Exit(1)
	}

	videoChan := make(chan string, len(videoFiles))
	var wg sync.WaitGroup
	numWorkers := 3

	p := mpb.New(mpb.WithWaitGroup(&wg))

	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for videoPath := range videoChan {
				err := extractFrames(videoPath, p)
				if err != nil {
					log.Fatal(err)
				}
			}
		}()
	}

	for _, videoPath := range videoFiles {
		videoChan <- videoPath
	}
	close(videoChan)

	wg.Wait()
}

func extractFrames(videoPath string, p *mpb.Progress) error {
	base := filepath.Base(videoPath)
	episodeKey := strings.TrimSuffix(filepath.Base(videoPath), filepath.Ext(base))
	outputDir := fmt.Sprintf("frames/%s", episodeKey)

	probe, err := ffmpeg.ProbeVideo(videoPath)
	if err != nil {
		return err
	}

	frameRate, err := probe.FrameRate()
	if err != nil {
		return err
	}

	totalFrames, err := probe.TotalFrames()
	if err != nil {
		return err
	}

	bar := newProgressBar(p, totalFrames, episodeKey)

	var significantFrame *image.RGBA
	err = ffmpeg.ReadFrames(videoPath, func(img *image.RGBA, frameNumber int) error {
		bar.SetCurrent(int64(frameNumber))
		start := time.Now()

		if frameNumber == 1 {
			significantFrame = img
			return nil
		}

		mean_ssim := ssim.MeanSSIM(significantFrame, img)
		if mean_ssim < similarityThreshold {
			significantFrame = img

			size := map[string]uint{
				"small":  240,
				"medium": 480,
				"large":  720,
			}
			for sizeName, imgWidth := range size {
				resizedImg := resize.Resize(imgWidth, 0, img, resize.Lanczos3)
				timestamp := frameNumberToMs(frameNumber, frameRate)
				filePath := fmt.Sprintf("%s/%d/%s.jpg", outputDir, timestamp, sizeName)

				err := saveAsJPEG(resizedImg, filePath)
				if err != nil {
					return err
				}
			}
		}

		bar.DecoratorEwmaUpdate(time.Since(start))

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

func newProgressBar(p *mpb.Progress, totalFrames int, episodeKey string) *mpb.Bar {
	return p.New(int64(totalFrames),
		mpb.BarStyle().Lbound("|").Filler("=").Tip(">").Padding("-").Rbound("|"),
		mpb.PrependDecorators(
			decor.Name(fmt.Sprintf("Processing %s: ", episodeKey)),
			decor.CountersNoUnit("%d/%d"),
			decor.Name(" ("),
			decor.Percentage(),
			decor.Name(")"),
		),
		mpb.AppendDecorators(
			decor.Elapsed(decor.ET_STYLE_GO),
			decor.EwmaSpeed(0, " %.2f ops/s", 60),
			decor.Name(" (ETA: "),
			decor.EwmaETA(decor.ET_STYLE_GO, 60),
			decor.Name(")"),
		),
	)
}
