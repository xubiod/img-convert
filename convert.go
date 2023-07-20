package main

import (
	"fmt"
	"image/gif"
	"image/jpeg"
	"image/png"
	"os"
	"path/filepath"

	"github.com/blezek/tga"
	"github.com/hhrutter/tiff"
	pnm "github.com/jbuchbinder/gopnm"
	"github.com/samuel/go-pcx/pcx"
	"github.com/xyproto/xpm"
	"golang.org/x/image/bmp"
	"lelux.net/x/image/qoi"
	"vimagination.zapto.org/limage/xcf"

	megaSD "github.com/bodgit/megasd/image"
	"github.com/chai2010/webp"
	"github.com/leotaku/mobi/jfif"
)

// ValidOutputTypes
//
// A list of valid encoders by file type.
var ValidOutputTypes = []string{
	"png",  // std/image
	"gif",  // std/image
	"jpeg", // std/image

	"bmp",  // golang.org/x/image
	"tiff", // golang.org/x/image

	"jfif",   // github.com/leotaku/mobi/jfif
	"webp",   // github.com/chai2010/webp
	"pbm",    // github.com/jbuchbinder/gopnm
	"pgm",    // github.com/jbuchbinder/gopnm
	"ppm",    // github.com/jbuchbinder/gopnm
	"pcx",    // github.com/samuel/go-pcx/pcx
	"megasd", // github.com/bodgit/megasd/image
	"qoi",    // lelux.net/x/image/qoi
	"tga",    // github.com/blezek/tga
	"xpm",    // github.com/xyproto/xpm
	"xcf",    // vimagination.zapto.org/limage/xcf
}

type QualityInformation struct {
	Lossless      bool
	QualityInt    int
	QualityFloat  float32
	WebpExact     bool
	TiffPredictor bool
}

func ConvertTo(filename string, outputFileType string, quality QualityInformation, overrideSameTypeSkip bool, overwriteFiles bool) (err error) {
	fmt.Printf("starting %s\n", filename)

	if (filepath.Ext(filename) == "."+outputFileType) && !overrideSameTypeSkip {
		return fmt.Errorf("%s already outout type, skipping", filename)
	}

	if _, err := os.Stat(filename + "." + outputFileType); err == nil && !overwriteFiles {
		return fmt.Errorf("%s already exists, skipping", filename+"."+outputFileType)
	}

	decodedImage, err := Import(filename)

	if err != nil {
		return fmt.Errorf("%s import failed, skipping (%s)", filename, err.Error())
	}

	r, err := os.Create(filename + "." + outputFileType)
	if err != nil {
		_ = r.Close()
		return fmt.Errorf("%s.%s couldn't be created, skipping (%s)", filename, outputFileType, err.Error())
	}

	switch outputFileType {
	case "bmp":
		err = bmp.Encode(r, decodedImage)
	case "tiff":
		err = tiff.Encode(r, decodedImage, &tiff.Options{
			Compression: tiff.CompressionType(quality.QualityInt),
			Predictor:   quality.TiffPredictor,
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
		err = webp.Encode(r, decodedImage, &webp.Options{
			Quality:  quality.QualityFloat,
			Lossless: quality.Lossless,
			Exact:    quality.WebpExact,
		})
	case "pbm":
		err = pnm.Encode(r, decodedImage, pnm.PBM)
	case "pgm":
		err = pnm.Encode(r, decodedImage, pnm.PGM)
	case "ppm":
		err = pnm.Encode(r, decodedImage, pnm.PPM)
	case "pcx":
		err = pcx.Encode(r, decodedImage)
	case "megasd":
		err = megaSD.Encode(r, decodedImage)
	case "qoi":
		err = qoi.Encode(r, decodedImage)
	case "tga":
		err = tga.Encode(r, decodedImage)
	case "xpm":
		err = xpm.Encode(r, decodedImage)
	case "xcf":
		err = xcf.Encode(r, decodedImage)
	}

	if err != nil {
		_ = r.Close()
		return fmt.Errorf("%s.%s couldn't be encoded (%s)", filename, outputFileType, err.Error())
	}

	_ = r.Close()

	return nil
}
