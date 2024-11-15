package main

import (
	"os"
)

func main() {
	img, err := PreProcess("/Users/alexyoung/downloads/monalisa.jpg", 70, 40)
	if err != nil {
		panic(err)
	}
	out :=Convert(img)
	os.WriteFile("ascii.txt", []byte(out), 0644)
}