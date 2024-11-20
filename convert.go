package quickscii

import (
	"fmt"
	"image"

	"gocv.io/x/gocv"
)

// Convert image to ASCII art string
// path: path to the image file
// w: width of the output ASCII art
// l: height of the output ASCII art
// charset: character set to use
//     - "block": block characters
//     - "poly": polygon characters
//     - "mix": text/blocks characters
// returns: ASCII art string, error
func Convert(path string, w, l int, charset string) (string, error) {
	// SECTION 1: Image Preprocessing
	// Read and validate the image
	if path == "" {
		return "", fmt.Errorf("invalid path: path cannot be empty")
	}
	if w <= 0 || l <= 0 {
		return "", fmt.Errorf("invalid dimensions: width and height must be positive")
	}

	img := gocv.IMRead(path, gocv.IMReadColor)
	if img.Empty() {
		return "", fmt.Errorf("failed to read image at path: %s", path)
	}
	defer img.Close()

	// Resize the image to specified dimensions
	resized := gocv.NewMat()
	defer resized.Close()
	gocv.Resize(img, &resized, image.Point{X: w, Y: l}, 0, 0, gocv.InterpolationNearestNeighbor)
	if resized.Empty() {
		return "", fmt.Errorf("failed to resize image")
	}

	// Convert to grayscale for ASCII conversion
	gray := gocv.NewMat()
	defer gray.Close()
	gocv.CvtColor(resized, &gray, gocv.ColorBGRToGray)
	if gray.Empty() {
		return "", fmt.Errorf("failed to convert image to grayscale")
	}

	// SECTION 2: ASCII Conversion
	// Define ASCII characters for different intensity levels
	block := []rune{'▪', '▦', '▥', '▤', '▓', '▒', '░', '▁', '▂', '▃', '▄', '▅', '▆', '▇', '█'}
	poly := []rune{'▫', '▧', '▨', '▥', '▲', '▱', '▯', '▰', '△', '▲', '△', '▲', '△', '▲', '△', '▲'}
	mix := []rune{' ', '.', ':', '-', '=', '+', '*', '#', '%', '@', '▁', '▂', '▃', '▄', '▅', '█'}

	var selectedCharset []rune
	switch charset {
	case "block":
		selectedCharset = block
	case "poly":
		selectedCharset = poly
	case "mix":
		selectedCharset = mix
	default:
		return "", fmt.Errorf("invalid charset: %s", charset)
	}
	out := ""

	// Iterate through each pixel and map to corresponding ASCII character
	rows := gray.Rows()
	cols := gray.Cols()
	if rows == 0 || cols == 0 {
		return "", fmt.Errorf("invalid image dimensions after processing")
	}

	for i := 0; i < rows; i++ {
		for j := 0; j < cols; j++ {
			pixel := gray.GetUCharAt(i, j)
			// Map pixel intensity to ASCII character
			index := int(pixel) * (len(selectedCharset) - 1) / 255
			out += string(selectedCharset[index])
		}
		out += "\n"
	}

	return out, nil
}