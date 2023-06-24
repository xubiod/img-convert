package jfif

import (
	"bytes"
	"image"
	"image/jpeg"
	"io"
)

// From https://github.com/leotaku/mobi/blob/master/jfif/jfif.go
var naiveJFIFHeader = []byte{
	0xFF, 0xD8, // SOI
	0xFF, 0xE0, // APP0 Marker
	0x00, 0x10, // Length
	0x4A, 0x46, 0x49, 0x46, 0x00, // JFIF\0
	0x01, 0x02, // 1.02
	0x00,       // Density type
	0x00, 0x01, // X Density
	0x00, 0x01, // Y Density
	0x00, 0x00, // No Thumbnail
}

func Decode(r io.Reader) (img image.Image, err error) {
	rawData, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}

	// Hijack the JFIF data to mash in a "Start of Image" marker
	// This technically makes it a valid JPEG
	var overwriteHere = int64(len(naiveJFIFHeader)) - 2

	// Cut off all the JFIF header except for two bytes that
	// make the "Start of Image" marker
	rawData = rawData[overwriteHere:]
	rawData[0] = 0xFF
	rawData[1] = 0xD8

	rawRead := bytes.NewReader(rawData)

	img, err = jpeg.Decode(rawRead)
	return img, err
}
