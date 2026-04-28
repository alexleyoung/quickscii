package quickscii

import (
	"fmt"
	"image"
	"image/color"
	_ "image/jpeg"
	_ "image/png"
	"os"
	"strings"

	"github.com/fogleman/gg"
	"golang.org/x/image/draw"
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
	if path == "" {
		return "", fmt.Errorf("invalid path: path cannot be empty")
	}

	if w <= 0 || l <= 0 {
		return "", fmt.Errorf("invalid dimensions: width and height must be positive")
	}

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

	f, err := os.Open(path)
	if err != nil {
		return "", fmt.Errorf("failed to read image at path: %s", path)
	}
	defer f.Close()

	src, _, err := image.Decode(f)
	if err != nil {
		return "", fmt.Errorf("failed to read image at path: %s", path)
	}

	resized := image.NewRGBA(image.Rect(0, 0, w, l))
	draw.NearestNeighbor.Scale(resized, resized.Bounds(), src, src.Bounds(), draw.Over, nil)

	var sb strings.Builder
	sb.Grow((w + 1) * l)

	for y := 0; y < l; y++ {
		for x := 0; x < w; x++ {
			i := resized.PixOffset(x, y)
			r, g, b := resized.Pix[i], resized.Pix[i+1], resized.Pix[i+2]
			// Rec. 601 luma — same coefficients OpenCV's BGR2GRAY uses.
			lum := (uint32(r)*299 + uint32(g)*587 + uint32(b)*114) / 1000
			idx := int(lum) * (len(selectedCharset) - 1) / 255
			sb.WriteRune(selectedCharset[idx])
		}
		sb.WriteByte('\n')
	}

	return sb.String(), nil
}

// Convert is an alias for Asciify, preserved for backwards compatibility.
func Convert(path string, w, l int, charset string) (string, error) {
	return Asciify(path, w, l, charset)
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
	const fontsize = 12

	dc := gg.NewContext(1, 1)
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

	ascii, err := Asciify(inputPath, w, l, charset)
	if err != nil {
		return fmt.Errorf("failed to generate ASCII art: %v", err)
	}

	lines := strings.Split(ascii, "\n")
	if len(lines) == 0 {
		return fmt.Errorf("no ASCII content generated")
	}

	width, _ := dc.MeasureString(lines[0])
	height := float64(len(lines)) * fontsize

	dc = gg.NewContext(int(width), int(height))

	dc.SetColor(color.Black)
	dc.Clear()

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

	dc.SetColor(color.White)

	for i, line := range lines {
		y := float64(i)*fontsize + fontsize
		dc.DrawString(line, 0, y)
	}

	if err := dc.SavePNG(outputPath); err != nil {
		return fmt.Errorf("failed to save image: %v", err)
	}

	return nil
}
