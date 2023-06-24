package jfif

import (
	"bytes"
	"image"
	"image/jpeg"
	"io"
)

const jfifHeaderSize = 20

func Decode(r io.Reader) (img image.Image, err error) {
	rawData, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}

	// Hijack the JFIF data to mash in a "Start of Image" marker
	// This technically makes it a valid JPEG
	var overwriteHere = int64(jfifHeaderSize) - 2

	// Cut off all the JFIF header except for two bytes that
	// make the "Start of Image" marker
	rawData = rawData[overwriteHere:]
	rawData[0] = 0xFF
	rawData[1] = 0xD8

	rawRead := bytes.NewReader(rawData)

	img, err = jpeg.Decode(rawRead)
	return img, err
}
