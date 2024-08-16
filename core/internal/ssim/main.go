package ssim

import (
	"image"
	"image/color"
	"math"
	"sync"
)

var (
	L  = 1.0
	K1 = 0.01
	K2 = 0.03
	C1 = math.Pow((K1 * L), 2.0)
	C2 = math.Pow((K2 * L), 2.0)
	C3 = C2 / 2.0
	windowSize = 11
)

func MeanSSIM(img1, img2 image.Image) float64 {
	var wg sync.WaitGroup
	localSSIMs := make(chan float64, (img1.Bounds().Dx()/windowSize)*(img1.Bounds().Dy()/windowSize))

	bounds := img1.Bounds()
	width, height := bounds.Dx(), bounds.Dy()

	for y := 0; y < height; y += windowSize {
		for x := 0; x < width; x += windowSize {
			wg.Add(1)
			go func(x, y int) {
				defer wg.Done()
				windowRect := image.Rect(x, y, x+windowSize, y+windowSize)

				subImg1 := subImage(img1, windowRect)
				subImg2 := subImage(img2, windowRect)

				localSSIMs <- SSIM(subImg1, subImg2)
			}(x, y)
		}
	}

	go func() {
		wg.Wait()
		close(localSSIMs)
	}()

	var results []float64
	for ssim := range localSSIMs {
		results = append(results, ssim)
	}

	return Mean(results)
}

func SSIM(img1, img2 image.Image) float64 {
	bounds := img1.Bounds()
	luminances1 := make([]float64, 0, bounds.Dx()*bounds.Dy())
	luminances2 := make([]float64, 0, bounds.Dx()*bounds.Dy())

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			luminances1 = append(luminances1, luminance(img1.At(x, y)))
			luminances2 = append(luminances2, luminance(img2.At(x, y)))
		}
	}

	mean1 := Mean(luminances1)
	mean2 := Mean(luminances2)
	variance1 := Variance(luminances1)
	variance2 := Variance(luminances2)
	covariance := Covariance(luminances1, luminances2)

	luminanceComparison := (2*mean1*mean2 + C1) / (mean1*mean1 + mean2*mean2 + C1)
	contrastComparison := (2*math.Sqrt(variance1)*math.Sqrt(variance2) + C2) / (variance1 + variance2 + C2)
	structureComparison := (covariance + C3) / (math.Sqrt(variance1)*math.Sqrt(variance2) + C3)

	return luminanceComparison * contrastComparison * structureComparison
}

func luminance(c color.Color) float64 {
	r, g, b, _ := c.RGBA()
	normR := float64(r) / 65535.0
	normG := float64(g) / 65535.0
	normB := float64(b) / 65535.0
	return 0.299*normR + 0.587*normG + 0.114*normB
}

func subImage(img image.Image, r image.Rectangle) image.Image {
	return img.(interface {
		SubImage(r image.Rectangle) image.Image
	}).SubImage(r)
}