package quickscii

import (
	"gocv.io/x/gocv"
)

// Asciify converts a grayscale image to ascii art
// INPUT: grayscale opencv matrix
// OUTPUT: string of ascii art
func Asciify(mat *gocv.Mat) string {
	
	ascii := []rune{'█', '▇', '▆', '▅', '▄', '▃', '▂', '▁', '░', '▒', '▓', '▤', '▥', '▪'}
	out := ""

	for i := 0; i < mat.Rows(); i++ {
		for j := 0; j < mat.Cols(); j++ {
			out += string(ascii[int(mat.GetUCharAt(i, j))*(len(ascii)-1)/255])
		}
		out += "k\n"
	}
	return out
}