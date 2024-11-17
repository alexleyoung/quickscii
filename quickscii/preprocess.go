package quickscii

import (
	"errors"
	"fmt"
	"image"

	"gocv.io/x/gocv"
)

// Resizes image to given width and height and converts to grayscale
func PreProcess(path string, width int, height int) (*gocv.Mat, error) {
	img := gocv.IMRead(path, gocv.IMReadColor)

	if img.Empty() {
		fmt.Println("Error reading image")
		return nil, errors.New("error reading image")
	}

	gocv.Resize(img, &img, image.Point{width, height}, 0, 0, gocv.InterpolationNearestNeighbor)
	gocv.CvtColor(img, &img, gocv.ColorBGRToGray)

	return &img, nil
}