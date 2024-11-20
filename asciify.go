package quickscii

import (
	"fmt"
	"image"
	"image/color"
	"strings"

	"github.com/fogleman/gg"
	"gocv.io/x/gocv"
)

// Asciify converts an image file into ASCII art.
//
// Parameters:
//   - path: Path to the image file
//   - w: Width of the output ASCII art
//   - l: Height of the output ASCII art
//   - charset: Character set to use for conversion:
//        • "block": Block characters (█▀▄)
//        • "poly": Polygon characters (◢◣◤◥)
//        • "mix": Mixed text and block characters
//
// Returns:
//   - string: The generated ASCII art
//   - error: Error if conversion fails
func Asciify(path string, w, l int, charset string) (string, error) {
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

// AsciifyToImage converts an image file into ASCII art and saves it as a PNG image.
//
// Parameters:
//   - inputPath: Path to the input image file
//   - outputPath: Path where the ASCII art image will be saved
//   - w: Width of the output ASCII art
//   - l: Height of the output ASCII art
//   - charset: Character set to use for conversion:
//        • "block": Block characters (█▀▄)
//        • "poly": Polygon characters (◢◣◤◥)
//        • "mix": Mixed text and block characters
//
// Returns:
//   - error: Error if conversion fails
func AsciifyToImage(inputPath, outputPath string, w, l int, charset string) error {
	// Get ASCII art string
	ascii, err := Asciify(inputPath, w, l, charset)
	if err != nil {
		return fmt.Errorf("failed to generate ASCII art: %v", err)
	}

	// Create a new image context
	// We'll make the image larger to accommodate the text
	const fontsize = 12
	const padding = 20
	imgWidth := float64(w * fontsize)
	imgHeight := float64(l * fontsize)
	
	dc := gg.NewContext(int(imgWidth + 2*padding), int(imgHeight + 2*padding))
	
	// Set white background
	dc.SetColor(color.White)
	dc.Clear()
	
	// Set text properties
	if err := dc.LoadFontFace("Go-Mono", fontsize); err != nil {
		// Fallback to a basic monospace font if Go-Mono is not available
		if err := dc.LoadFontFace("Courier", fontsize); err != nil {
			return fmt.Errorf("failed to load font: %v", err)
		}
	}
	
	dc.SetColor(color.Black)
	
	// Draw the ASCII art
	lines := strings.Split(ascii, "\n")
	for i, line := range lines {
		y := float64(i)*fontsize + padding
		dc.DrawString(line, padding, y)
	}
	
	// Save the image
	if err := dc.SavePNG(outputPath); err != nil {
		return fmt.Errorf("failed to save image: %v", err)
	}
	
	return nil
}