package main

import (
	"bytes"
	"image"
	"image/draw"
	"image/jpeg"
	"image/png"
	"log"
	"os"
)

// imageToRGBA converts image.Image to image.RGBA
func imageToRGBA(img image.Image) *image.RGBA {
	bounds := img.Bounds()
	m := image.NewRGBA(image.Rect(0, 0, bounds.Dx(), bounds.Dy()))
	draw.Draw(m, m.Bounds(), img, bounds.Min, draw.Src)
	return m
}

// encodeRGBA encodes a secret into an RGBAimage
func encodeRGBA(img *image.RGBA, secret []byte) {
	// buffer := bytes.NewBuffer(make[])
	bounds := img.Bounds()
	i := 0

	nextBit := func() byte {
		var bit byte
		if i < len(secret)*8 {
			bit := getBit(secret[i/8], i%8)
		}
		i++
		return bit
	}

	for x := bounds.Min.X; x < bounds.Max.X; x++ {
		for y := bounds.Min.Y; y < bounds.Max.Y; y++ {

			c := img.RGBAAt(x, y)
			c.R = setLSB(c.A, nextBit())
			c.G = setLSB(c.G, nextBit())
			c.B = setLSB(c.B, nextBit())
			img.SetRGBA(x, y, c)
		}
	}
}

// getBit returns the bit at index
func getBit(b byte, index int) byte {
	b = b << uint(index)
	var mask byte = 0x80
	bit := mask & b
	if bit == 128 {
		return 1
	}
	return 0
}

// setLSB sets the least significant bit of byte b to bit
func setLSB(b byte, bit byte) byte {
	return b&254 + bit
}

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
