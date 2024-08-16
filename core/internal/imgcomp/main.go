package imgcomp

import (
	"image"
	"image/color"
	"math"
)

var (
	L = 1.0
	K1 = 0.01
	K2 = 0.03
	C1 = math.Pow((K1*L), 2.0)
	C2 = math.Pow((K2*L), 2.0)
	C3 = C2 / 2.0
  )

func SSIM(img1, img2 image.Image) float64 {
	windowSize := 11
	var localSSIMs []float64

	bounds := img1.Bounds()
	width, height := bounds.Dx(), bounds.Dy()

	for y := 0; y < height; y += windowSize {
		for x := 0; x < width; x += windowSize {
			windowRect := image.Rect(x, y, x+windowSize, y+windowSize)

            subImg1 := img1.(interface {
                SubImage(r image.Rectangle) image.Image
            }).SubImage(windowRect)

            subImg2 := img2.(interface {
                SubImage(r image.Rectangle) image.Image
            }).SubImage(windowRect)

			localSSIMs = append(localSSIMs, localSSIM(subImg1, subImg2))
		}
	}

	return mean(localSSIMs)
}

func localSSIM(img1, img2 image.Image) float64 {
    luminances1 := pixelLuminances(img1)
    luminances2 := pixelLuminances(img2)

    mean1 := mean(luminances1)
    mean2 := mean(luminances2)
    variance1 := variance(luminances1)
    variance2 := variance(luminances2)
    covariance := covariance(luminances1, luminances2)

    luminanceComparison := (2 * mean1 * mean2 + C1) / (mean1*mean1 + mean2*mean2 + C1)
    contrastComparison := (2 * math.Sqrt(variance1) * math.Sqrt(variance2) + C2) / (variance1 + variance2 + C2)
    structureComparison := (covariance + C3) / (math.Sqrt(variance1) * math.Sqrt(variance2) + C3)

    return luminanceComparison * contrastComparison * structureComparison
}

func pixelLuminances(img image.Image) []float64 {
    bounds := img.Bounds()
    luminances := make([]float64, 0, bounds.Dx()*bounds.Dy())

    for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
        for x := bounds.Min.X; x < bounds.Max.X; x++ {
            luminances = append(luminances, luminance(img.At(x, y)))
        }
    }

    return luminances
}

func luminance(c color.Color) float64 {
    r, g, b, _ := c.RGBA()
    normR := float64(r) / 65535.0
    normG := float64(g) / 65535.0
    normB := float64(b) / 65535.0
    return 0.299*normR + 0.587*normG + 0.114*normB
}

func mean(nums []float64) float64 {
    sum := 0.0
    for _, num := range nums {
        sum += num
    }
    return sum / float64(len(nums))
}

func variance(nums []float64) float64 {
    mean := mean(nums)
    sum := 0.0

    for _, num := range nums {
        sum += (num - mean) * (num - mean)
    }

    return sum / float64(len(nums))
}

func covariance(nums1, nums2 []float64) float64 {
    if len(nums1) != len(nums2) {
        panic("slices must have the same length")
    }

    n := len(nums1)
    if n == 0 {
        panic("slices must not be empty")
    }

    mean1 := mean(nums1)
    mean2 := mean(nums2)
    sum := 0.0

    for i, num1 := range nums1 {
        num2 := nums2[i]
        sum += (num1 - mean1) * (num2 - mean2)
    }

    return sum / float64(len(nums1))
}