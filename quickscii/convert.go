package quickscii

func Convert(path string, w, l int) string {
	processed, err := PreProcess(path, w, l)
	if err != nil {
		return err.Error()
	}
	return Asciify(processed)
}