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
	"github.com/kevin-cantwell/dotmatrix"
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

	"dotmatrix.txt", // github.com/kevin-cantwell/dotmatrix
}

type QualityInformation struct {
	Lossless      bool
	Quality       int
	WebpExact     bool
	TiffPredictor bool
}

func ConvertTo(filename string, outputFileType string, quality QualityInformation, overrideSameTypeSkip bool, overwriteFiles bool) (err error) {
	fmt.Printf("starting %s\n", filename)

	if (filepath.Ext(filename) == "."+outputFileType) && !overrideSameTypeSkip {
		return fmt.Errorf("%s already output type, skipping", filename)
	}

	if _, err := os.Stat(filename + "." + outputFileType); err == nil && !overwriteFiles {
		return fmt.Errorf("%s already exists, skipping", filename+"."+outputFileType)
	}

	imported, err := Import(filename)

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
		err = bmp.Encode(r, imported)

	case "tiff":
		err = tiff.Encode(r, imported, &tiff.Options{
			Compression: tiff.CompressionType(quality.Quality), Predictor: quality.TiffPredictor,
		})

	case "gif":
		err = gif.Encode(r, imported, &gif.Options{NumColors: quality.Quality})

	case "jpeg", "jpg":
		err = jpeg.Encode(r, imported, &jpeg.Options{Quality: quality.Quality})

	case "png":
		err = png.Encode(r, imported)

	case "jfif":
		err = jfif.Encode(r, imported, &jpeg.Options{Quality: quality.Quality})

	case "webp":
		err = webp.Encode(r, imported, &webp.Options{
			Quality: float32(quality.Quality), Lossless: quality.Lossless, Exact: quality.WebpExact,
		})

	case "pbm":
		err = pnm.Encode(r, imported, pnm.PBM)

	case "pgm":
		err = pnm.Encode(r, imported, pnm.PGM)

	case "ppm":
		err = pnm.Encode(r, imported, pnm.PPM)

	case "pcx":
		err = pcx.Encode(r, imported)

	case "megasd":
		err = megaSD.Encode(r, imported)

	case "qoi":
		err = qoi.Encode(r, imported)

	case "tga":
		err = tga.Encode(r, imported)

	case "xpm":
		err = xpm.Encode(r, imported)

	case "xcf":
		err = xcf.Encode(r, imported)

	case "dotmatrix.txt":
		err = (*dotmatrix.NewPrinter(r, &dotmatrix.Config{})).Print(imported)
	}

	if err != nil {
		_ = r.Close()
		return fmt.Errorf("%s.%s couldn't be encoded (%s)", filename, outputFileType, err.Error())
	}

	_ = r.Close()

	return nil
}
