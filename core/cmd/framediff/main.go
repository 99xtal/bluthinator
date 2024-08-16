package main

import (
	"fmt"
	"image/jpeg"
	"log"
	"os"

	"github.com/99xtal/bluthinator/core/internal/ssim"
)

func main() {
	imgFile1, err := os.Open("image1.jpg")
	if err != nil {
		log.Fatal(err)
	}
	defer imgFile1.Close()

	imgFile2, err := os.Open("image2.jpg")
	if err != nil {
		log.Fatal(err)
	}
	defer imgFile2.Close()

	img1, err := jpeg.Decode(imgFile1)
	if err != nil {
		log.Fatal(err)
	}

	img2, err := jpeg.Decode(imgFile2)
	if err != nil {
		log.Fatal(err)
	}
	
	ssim_index := ssim.GetSSIMIndex(img1, img2)
	fmt.Printf("SSIM index: %f\n", ssim_index)
}
