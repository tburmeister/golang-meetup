package main

import (
	"bytes"
	"flag"
	"image/jpeg"
	"image/png"
	"io/ioutil"
	"log"
)

func main() {
	ff := flag.String("img", "", "The image to process")
	flag.Parse()

	filename := *ff
	log.Printf("image: %s\n", filename)

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
