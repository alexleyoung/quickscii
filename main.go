package main

func main() {
	img, err := PreProcess("test.png", 100, 100)
	if err != nil {
		panic(err)
	}
	Convert(img)
}