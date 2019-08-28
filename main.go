package main

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"crypto/rand"
	"encoding/hex"
	"flag"
	"fmt"
	"image"
	"image/draw"
	"image/jpeg"
	"image/png"
	"io"
	"io/ioutil"
	"regexp"
)

// Usage:
// ./golang-meetup -img gopher.jpeg -msg hey
// ./golang-meetup -img gopher-encrypted.jpeg -decrypt

// imageToRGBA converts image.Image to image.RGBA
func imageToRGBA(img image.Image) *image.RGBA {
	bounds := img.Bounds()
	m := image.NewRGBA(image.Rect(0, 0, bounds.Dx(), bounds.Dy()))
	draw.Draw(m, m.Bounds(), img, bounds.Min, draw.Src)
	return m
}

// encodeRGBA encodes a secret into an RGBAimage
func encodeRGBA(img *image.RGBA, secret []byte) {
	bounds := img.Bounds()
	i := 0

	nextBit := func() byte {
		var bit byte
		if i < len(secret)*8 {
			bit = getBit(secret[i/8], i%8)
		}
		i++
		return bit
	}

	message := make([]byte, 1, len(secret)+1)
	message[0] = byte(len(secret))
	message = append(message, secret...)

	for x := bounds.Min.X; x < bounds.Max.X; x++ {
		for y := bounds.Min.Y; y < bounds.Max.Y; y++ {

			c := img.RGBAAt(x, y)
			c.R = setLSB(c.R, nextBit())
			c.G = setLSB(c.G, nextBit())
			c.B = setLSB(c.B, nextBit())
			img.SetRGBA(x, y, c)
		}
	}
}

func decodeRGBA(img *image.RGBA) []byte {
	bounds := img.Bounds()
	secret := make([]byte, bounds.Dx()*bounds.Dy())
	i := 0
	for x := bounds.Min.X; x < bounds.Max.X; x++ {
		for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
			c := img.RGBAAt(x, y)
			secret[i/8] = setBit(secret[i/8], i%8, getLSB(c.R))
			i++
			secret[i/8] = setBit(secret[i/8], i%8, getLSB(c.G))
			i++
			secret[i/8] = setBit(secret[i/8], i%8, getLSB(c.B))
			i++
		}
	}
	length := int(secret[0])
	return secret[1 : length+1]
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

// setBit returns byte with the bit at index set to bit
func setBit(b byte, index int, bit byte) byte {
	var mask byte = 0x80
	mask = mask >> uint(index)

	if bit == 0 {
		mask = ^mask
		b = b & mask
	} else if bit == 1 {
		b = b | mask
	}
	return b
}

// getLSB returns the least significant bit of byte b
func getLSB(b byte) byte {
	return getBit(b, 0)
}

// setLSB sets the least significant bit of byte b to bit
func setLSB(b byte, bit byte) byte {
	return b&254 + bit
}

const (
	formatJpeg = iota
	formatPng

	passphrase = "chopped"
)

func encode(img image.Image, secret []byte) (image.Image, error) {
	rgba := imageToRGBA(img)
	encodeRGBA(rgba, secret)
	return rgba, nil
}

func decode(img image.Image) ([]byte, error) {
	rgba := imageToRGBA(img)
	return decodeRGBA(rgba), nil
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
		panic(fmt.Sprintf("unable to read file %s", filename))
	}

	reader := bytes.NewReader(raw)
	img, err := jpeg.Decode(reader)
	format := formatJpeg
	if err != nil {
		img, err = png.Decode(reader)
		format = formatPng
	}
	fmt.Printf("format: %d\n", format)
	if err != nil {
		panic("unable to decode image - not jpeg or png")
	}

	if *df {
		msg, err := decode(img)
		if err != nil {
			panic(fmt.Sprintf("unable to decode image: %s", err))
		}
		fmt.Println(string(msg))
		//plaintext := decrypt([]byte(msg), passphrase)
		//fmt.Println(plaintext)
	} else {
		ciphertext := encrypt([]byte(msg), passphrase)
		ciphertext = []byte(msg)
		new, err := encode(img, ciphertext)
		if err != nil {
			panic(fmt.Sprintf("unable to encrypt image: %s", err))
		}

		var buf bytes.Buffer
		writer := io.Writer(&buf)

		switch format {
		case formatJpeg:
			err = jpeg.Encode(writer, new, &jpeg.Options{Quality: 100})
			if err != nil {
				panic(fmt.Sprintf("unable to encode jpeg: %s", err))
			}
		case formatPng:
			err = png.Encode(writer, new)
			if err != nil {
				panic(fmt.Sprintf("unable to encode png: %s", err))
			}
		}

		re := regexp.MustCompile(`(\w+)\.(png|jpeg)`)
		err = ioutil.WriteFile(re.ReplaceAllString(filename, `$1-encrypted.$2`), buf.Bytes(), 0644)
		if err != nil {
			panic(err)
		}
	}
}

func encrypt(data []byte, passphrase string) []byte {
	block, _ := aes.NewCipher([]byte(_badHash(passphrase)))
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		panic(err.Error())
	}
	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		panic(err.Error())
	}
	ciphertext := gcm.Seal(nonce, nonce, data, nil)
	return ciphertext
}

func decrypt(data []byte, passphrase string) []byte {
	key := []byte(_badHash(passphrase))
	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err.Error())
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		panic(err.Error())
	}
	nonceSize := gcm.NonceSize()
	nonce, ciphertext := data[:nonceSize], data[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		panic(err.Error())
	}
	return plaintext
}

func _badHash(passphrase string) string {
	hasher := md5.New()
	hasher.Write([]byte(passphrase))
	return hex.EncodeToString(hasher.Sum(nil))
}
