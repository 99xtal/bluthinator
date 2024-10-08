package commands

import (
	"encoding/csv"
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
	"github.com/spf13/cobra"
	"github.com/vbauerster/mpb/v7"
)

var (
	similarityThreshold float64
	numWorkers          uint
	outputDir           string
)

// extractCmd represents the extract command
var extractCmd = &cobra.Command{
	Use:   "extract [video_dir]",
	Short: "Extract perceptually distinct frames and metadata from a video",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		inputDirPath := args[0]

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

		p := mpb.New(mpb.WithWaitGroup(&wg))

		for i := uint(0); i < numWorkers; i++ {
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
	},
}

func init() {
	rootCmd.AddCommand(extractCmd)

	extractCmd.Flags().Float64VarP(&similarityThreshold, "threshold", "t", 0.70, "Threshold for frame similarity")
	extractCmd.Flags().UintVarP(&numWorkers, "workers", "w", 3, "Number of workers to use")
	extractCmd.Flags().StringVarP(&outputDir, "output", "o", "./output", "Output directory for extracted frames and metadata")
}

func extractFrames(videoPath string, p *mpb.Progress) error {
	fileName := filename(videoPath)
	frameDir := fmt.Sprintf("%s/%s/frames", outputDir, fileName)
	if err := os.MkdirAll(frameDir, 0755); err != nil {
		return err
	}

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

	// Open CSV file for writing
	csvFile, err := os.Create(fmt.Sprintf("%s/index.csv", frameDir))
	if err != nil {
		return err
	}
	defer csvFile.Close()

	// Create CSV writer
	csvWriter := csv.NewWriter(csvFile)
	defer csvWriter.Flush()

	header := []string{"episode", "timestamp"}
	if err := csvWriter.Write(header); err != nil {
		return err
	}

	bar := newProgressBar(p, int64(totalFrames), fmt.Sprintf("Processing %s: ", fileName))

	var lastSignificantFrame *image.RGBA
	err = ffmpeg.ReadFrames(videoPath, func(img *image.RGBA, frameNumber int) error {
		bar.SetCurrent(int64(frameNumber))
		start := time.Now()

		if frameNumber == 1 {
			lastSignificantFrame = img
			return nil
		}

		mean_ssim := ssim.MeanSSIM(lastSignificantFrame, img)
		if mean_ssim < similarityThreshold {
			lastSignificantFrame = img
			timestamp := frameNumberToMs(frameNumber, frameRate)

			// write to index file
			record := []string{fileName, fmt.Sprintf("%d", timestamp)}
			if err := csvWriter.Write(record); err != nil {
				return err
			}

			// write images
			imageDir := fmt.Sprintf("%s/%d", frameDir, timestamp)
			err := writeImages(lastSignificantFrame, imageDir)
			if err != nil {
				return err
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

func writeImages(img image.Image, outputDir string) error {
	imgSizes := map[string]uint{
		"small":  240,
		"medium": 480,
		"large":  720,
	}

	var wg sync.WaitGroup
	errChan := make(chan error, len(imgSizes))

	for sizeName, imgWidth := range imgSizes {
		wg.Add(1)
		go func(sizeName string, imgWidth uint) {
			defer wg.Done()

			resizedImg := resize.Resize(imgWidth, 0, img, resize.Lanczos3)
			filePath := fmt.Sprintf("%s/%s.jpg", outputDir, sizeName)

			err := saveAsJPEG(resizedImg, filePath)
			if err != nil {
				errChan <- err
			}
		}(sizeName, imgWidth)
	}

	wg.Wait()
	close(errChan)

	for err := range errChan {
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

func filename(path string) string {
	return strings.TrimSuffix(filepath.Base(path), filepath.Ext(path))
}
