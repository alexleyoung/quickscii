package main

import (
	"fmt"
	"os"
)

func main() {
	img, err := PreProcess("/Users/alexyoung/downloads/monalisa.jpg", 70, 40)
	if err != nil {
		panic(err)
	}
	out :=Convert(img)
	fmt.Println(out)
	os.WriteFile("ascii.txt", []byte(out), 0644)
}