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
	"image/jpeg"
	"image/png"
	"io"
	"io/ioutil"
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

var passphrase = "chopped"

func testmain() {
	var plaintext = "this is a test string"
	fmt.Println(encryptText([]byte(plaintext), passphrase))

	var crypted = []byte{56, 93, 224, 168, 120, 26, 240, 96, 173, 87, 39, 111, 204, 188, 151, 201, 171, 250, 49, 44, 39, 241, 164, 129, 64, 132, 158, 247, 88, 113, 172, 241, 216, 144, 94, 189, 197, 244, 119, 221, 57, 176, 132, 195, 211, 148, 4, 29, 72}
	fmt.Println(string(decryptText(crypted, passphrase)))
}

func main() {
	ff := flag.String("img", "", "The image to process")
	df := flag.Bool("decrypt", false, "Flag on to decrypt image (default encrypts)")
	mf := flag.String("msg", "", "Message to encrypt in image")
	flag.Parse()

	filename := *ff
	fmt.Printf("image: %s\n", filename)

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
		// Decrypt image
	} else {
		// Encrypt image
		_ = *mf
	}

	_ = img
}

func encryptText(data []byte, passphrase string) []byte {
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

func decryptText(data []byte, passphrase string) []byte {
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
