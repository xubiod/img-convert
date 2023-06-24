package main

import (
	"fmt"
	"golang.org/x/image/bmp"
	"golang.org/x/image/tiff"
	"golang.org/x/image/vp8l"
	"golang.org/x/image/webp"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"os"
	"path/filepath"

	webpEncoder "github.com/chai2010/webp"
	"github.com/leotaku/mobi/jfif"
)

// Valid decoders
var ValidInputTypes = []string{
	"bmp",  // golang.org/x/image
	"tiff", // golang.org/x/image
	"vp8l", // golang.org/x/image
	"webp", // golang.org/x/image

	"gif",  // std/image
	"jpeg", // std/image
	"jpg",  // std/image
	"png",  // std/image
}

// Valid encoders
var ValidOutputTypes = []string{
	"bmp",  // golang.org/x/image
	"tiff", // golang.org/x/image

	"gif",  // std/image
	"jpeg", // std/image
	"jpg",  // std/image
	"png",  // std/image

	"jfif", // github.com/leotaku/mobi/jfif
	"webp", // github.com/chai2010/webp
}

type QualityInformation struct {
	Lossless     bool
	QualityInt   int
	QualityFloat float32
	WebpExact    bool
}

func ConvertTo(filename string, outputFileType string, quality QualityInformation, overrideSameTypeSkip bool, overwriteFiles bool) error {
	fmt.Printf("starting %s\n", filename)

	var decodedImage image.Image

	if (filepath.Ext(filename) == "."+outputFileType) && !overrideSameTypeSkip {
		return fmt.Errorf("%s already outout type, skipping", filename)
	}

	if _, err := os.Stat(filename + "." + outputFileType); err == nil && !overwriteFiles {
		return fmt.Errorf("%s already exists, skipping", filename+"."+outputFileType)
	}

	f, err := os.Open(filename)
	if err != nil {
		f.Close()
		return fmt.Errorf("%s couldn't be opened, skipping (%s)", filename, err.Error())
	}

	inputValid := false
	for _, item := range ValidInputTypes {
		if filepath.Ext(filename) == "."+item {
			inputValid = true
		}
	}

	if !inputValid {
		f.Close()
		return fmt.Errorf("%s is not a valid input, skipping", filename)
	}

	switch filepath.Ext(filename) {
	case ".bmp":
		decodedImage, err = bmp.Decode(f)

	case ".tiff":
		decodedImage, err = tiff.Decode(f)

	case ".vp8l":
		decodedImage, err = vp8l.Decode(f)

	case ".webp":
		decodedImage, err = webp.Decode(f)

	case ".gif":
		decodedImage, err = gif.Decode(f)

	case ".jpg", ".jpeg":
		decodedImage, err = jpeg.Decode(f)

	case ".png":
		decodedImage, err = png.Decode(f)
	}

	if err != nil {
		f.Close()
		return fmt.Errorf("%s couldn't be decoded (%s), skipping", filename, err.Error())
	}

	r, err := os.Create(filename + "." + outputFileType)
	if err != nil {
		f.Close()
		r.Close()
		return fmt.Errorf("%s.%s couldn't be created, skipping (%s)", filename, outputFileType, err.Error())
	}

	switch outputFileType {
	case "bmp":
		err = bmp.Encode(r, decodedImage)
	case "tiff":
		err = tiff.Encode(r, decodedImage, &tiff.Options{
			Compression: tiff.CompressionType(quality.QualityInt),
		})
	case "gif":
		err = gif.Encode(r, decodedImage, &gif.Options{
			NumColors: quality.QualityInt,
		})
	case "jpeg", "jpg":
		err = jpeg.Encode(r, decodedImage, &jpeg.Options{
			Quality: quality.QualityInt,
		})
	case "png":
		err = png.Encode(r, decodedImage)
	case "jfif":
		err = jfif.Encode(r, decodedImage, &jpeg.Options{
			Quality: quality.QualityInt,
		})
	case "webp":
		err = webpEncoder.Encode(r, decodedImage, &webpEncoder.Options{
			Quality:  quality.QualityFloat,
			Lossless: quality.Lossless,
			Exact:    quality.WebpExact,
		})
	}

	if err != nil {
		f.Close()
		r.Close()
		return fmt.Errorf("%s.%s couldn't be encoded (%s)", filename, outputFileType, err.Error())
	}

	f.Close()
	r.Close()

	return nil
}
