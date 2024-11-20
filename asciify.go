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
	// Create a new image context
	// Calculate dimensions based on font metrics
	const fontsize = 12
	
	dc := gg.NewContext(1, 1) // Temporary context to measure text
	fonts := []string{
		"/System/Library/Fonts/Monaco.ttf",
		"/System/Library/Fonts/Menlo.ttf",
		"/Library/Fonts/Courier New.ttf",
		"/Library/Fonts/Courier New Bold.ttf",
	}

	fontLoaded := false
	for _, font := range fonts {
		if err := dc.LoadFontFace(font, fontsize); err == nil {
			fontLoaded = true
			break
		}
	}

	if !fontLoaded {
		return fmt.Errorf("failed to load any monospace font")
	}

	// Get the ASCII art to measure dimensions
	ascii, err := Asciify(inputPath, w, l, charset)
	if err != nil {
		return fmt.Errorf("failed to generate ASCII art: %v", err)
	}

	lines := strings.Split(ascii, "\n")
	if len(lines) == 0 {
		return fmt.Errorf("no ASCII content generated")
	}

	// Measure the exact width and height needed
	width, _ := dc.MeasureString(lines[0])
	height := float64(len(lines)) * fontsize

	// Create the final context with exact dimensions
	dc = gg.NewContext(int(width), int(height))
	
	// Set black background
	dc.SetColor(color.Black)
	dc.Clear()
	
	// Load the font again for the new context
	fontLoaded = false
	for _, font := range fonts {
		if err := dc.LoadFontFace(font, fontsize); err == nil {
			fontLoaded = true
			break
		}
	}

	if !fontLoaded {
		return fmt.Errorf("failed to load any monospace font")
	}
	
	// Set white text color
	dc.SetColor(color.White)
	
	// Draw the ASCII art
	for i, line := range lines {
		y := float64(i)*fontsize + fontsize // Add fontsize to y to account for baseline
		dc.DrawString(line, 0, y)
	}
	
	// Save the image
	if err := dc.SavePNG(outputPath); err != nil {
		return fmt.Errorf("failed to save image: %v", err)
	}
	
	return nil
}