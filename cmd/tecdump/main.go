package main

import (
	"bytes"
	"image"
	"image/png"
	"io"
	"log"
	"os"

	"vallon.me/redshift/disk"
)

func init() {
	log.SetFlags(0)
}

func main() {
	if len(os.Args) != 2 {
		log.Fatal(os.Stderr, "usage: tecdump disk.png > out.solution")
	}

	f, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	img, err := png.Decode(f)
	if err != nil {
		log.Fatal(err)
	}

	var pix io.Reader
	switch img := img.(type) {
	case *image.NRGBA:
		pix = bytes.NewReader(img.Pix)
	case *image.RGBA:
		pix = bytes.NewReader(img.Pix)
	default:
		log.Fatal("unsupported image format")
	}

	d, err := disk.NewReader(pix)
	if err != nil {
		log.Fatal(err)
	}

	var buf bytes.Buffer
	if _, err := buf.ReadFrom(d); err != nil && err != io.EOF {
		log.Fatalf("corrupt image %s", err)
	}

	buf.WriteTo(os.Stdout)
}
