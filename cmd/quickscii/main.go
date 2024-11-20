package main

import (
	"fmt"
	"log"
	"path/filepath"
	"runtime"

	"github.com/alexleyoung/quickscii"
)

func main() {
	// Get the absolute path to the project root
	_, filename, _, _ := runtime.Caller(0)
	projectRoot := filepath.Join(filepath.Dir(filename), "../..")
	
	// Construct absolute paths
	testImagePath := filepath.Join(projectRoot, "testdata", "nerd.jpg")
	outputPath := filepath.Join(projectRoot, "output.png")

	// Test Asciify function
	ascii, err := quickscii.Asciify(testImagePath, 80, 40, "mix")
	if err != nil {
		log.Fatalf("Error converting to ASCII: %v", err)
	}
	fmt.Println("ASCII Art Result:")
	fmt.Println(ascii)

	// Test AsciifyToImage function
	err = quickscii.AsciifyToImage(testImagePath, outputPath, 450, 200, "mix")
	if err != nil {
		log.Fatalf("Error creating ASCII image: %v", err)
	}
	fmt.Println("ASCII image has been saved to output.png")
}
