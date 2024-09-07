package main

import (
	"fmt"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"io"
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

var genericExporters = map[string]func(io.Writer, image.Image) error{
	"bmp":    bmp.Encode,
	"png":    png.Encode,
	"pcx":    pcx.Encode,
	"megasd": megaSD.Encode,
	"qoi":    qoi.Encode,
	"tga":    tga.Encode,
	"xpm":    xpm.Encode,
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
	case "tiff":
		err = tiff.Encode(r, imported, &tiff.Options{
			Compression: tiff.CompressionType(quality.Quality), Predictor: quality.TiffPredictor,
		})

	case "gif":
		err = gif.Encode(r, imported, &gif.Options{NumColors: quality.Quality})

	case "jpeg", "jpg":
		err = jpeg.Encode(r, imported, &jpeg.Options{Quality: quality.Quality})

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

	case "xcf":
		err = xcf.Encode(r, imported)

	case "dotmatrix.txt":
		err = (*dotmatrix.NewPrinter(r, &dotmatrix.Config{})).Print(imported)

	_:
		err = genericExporters[outputFileType](r, imported)
	}

	if err != nil {
		_ = r.Close()
		return fmt.Errorf("%s.%s couldn't be encoded (%s)", filename, outputFileType, err.Error())
	}

	_ = r.Close()

	return nil
}
