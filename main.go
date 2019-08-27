package main

import (
	"bytes"
	"image/jpeg"
	"image/png"
	"io/ioutil"
	"log"
)

func main() {
	filename := "test"
	raw, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatalf("unable to read file %s", filename)
	}

	reader := bytes.NewReader(raw)
	img, err := jpeg.Decode(reader)
	if err != nil {
		img, err = png.Decode(reader)
	}
	if err != nil {
		log.Fatal("unable to decode image - not jpeg or png")
		return
	}

	_ = img
}
