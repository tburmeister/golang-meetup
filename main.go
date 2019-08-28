package main

import (
	"bytes"
	"flag"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"io/ioutil"
	"log"
)

const (
	formatJpeg = iota
	formatPng
)

func encrypt(img image.Image, str string) (image.Image, error) {
	return nil, fmt.Errorf("unimplemented")
}

func decrypt(img image.Image) (string, error) {
	return "", fmt.Errorf("unimplemented")
}

func main() {
	ff := flag.String("img", "", "The image to process")
	df := flag.Bool("decrypt", false, "Flag on to decrypt image (default encrypts)")
	mf := flag.String("msg", "", "Message to encrypt in image")
	flag.Parse()

	filename := *ff
	msg := *mf

	raw, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatalf("unable to read file %s", filename)
	}

	reader := bytes.NewReader(raw)
	img, err := jpeg.Decode(reader)
	format := formatJpeg
	if err != nil {
		img, err = png.Decode(reader)
		format = formatPng
	}
	if err != nil {
		log.Fatal("unable to decode image - not jpeg or png")
		return
	}

	if *df {
		msg, err := decrypt(img)
		if err != nil {
			log.Fatalf("unable to decrypt image: %s", err)
			return
		}
		fmt.Println(msg)
	} else {
		new, err := encrypt(img, msg)
		if err != nil {
			log.Fatalf("unable to encrypt image: %s", err)
			return
		}

		raw := make([]byte, 0)
		ext := ""

		switch format {
		case formatJpeg:
			jpeg.Encode(raw, new, &jpeg.Options{100})
		case formatPng:
			png.Encode(raw, new)
		}

		err = ioutil.WriteFile(filename[:len(filename)-4] + "-encrypted" + ext, raw, 0644)
		if err != nil {
			panic(err)
		}
	}
}
