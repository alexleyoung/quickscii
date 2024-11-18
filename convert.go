package quickscii

// Convert converts an image to ascii art
// INPUT: path to image, width, height
// OUTPUT: string of ascii art
func Convert(path string, w, l int) string {
	processed, err := PreProcess(path, w, l)
	if err != nil {
		return err.Error()
	}
	return Asciify(processed)
}