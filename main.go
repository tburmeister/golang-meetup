package main

import (
	"image/jpeg"
	"image/png"
	"log"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal("must supply an image")
		return
	}

	file, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatalf("unable to open %s", os.Args[1])
		return
	}
	defer file.Close()

	img, err := jpeg.Decode(file)
	if err != nil {
		img, err = png.Decode(file)
	}
	if err != nil {
		log.Fatal("unable to decode image - not jpeg or png")
		return
	}

	_ = img
}
