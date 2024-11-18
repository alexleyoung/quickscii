package main

import (
	"os"

	"github.com/alexleyoung/quickscii/quickscii"
)

func main() {
	img, err := quickscii.PreProcess("/Users/alexyoung/downloads/monalisa.jpg", 140, 70)
	if err != nil {
		panic(err)
	}
	out := quickscii.Convert(img)
	os.WriteFile("ascii.txt", []byte(out), 0644)
}