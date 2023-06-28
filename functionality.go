package main

import (
	"fmt"
	"github.com/askeladdk/aseprite"
	"github.com/blezek/tga"
	pnm "github.com/jbuchbinder/gopnm"
	"github.com/mokiat/goexr/exr"
	"github.com/nielsAD/gowarcraft3/file/blp"
	"github.com/oov/psd"
	"github.com/samuel/go-pcx/pcx"
	"github.com/xyproto/xpm"
	"golang.org/x/image/bmp"
	"golang.org/x/image/tiff"
	"golang.org/x/image/vp8l"
	"golang.org/x/image/webp"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"lelux.net/x/image/qoi"
	"os"
	"path/filepath"
	"strings"
	"vimagination.zapto.org/limage/xcf"

	megaSD "github.com/bodgit/megasd/image"
	webpEncoder "github.com/chai2010/webp"
	"github.com/leotaku/mobi/jfif"

	selfJfif "github.com/xubiod/img-convert/dedicated-decoder/jfif"
)

// ValidInputTypes
//
// A list of valid decoders by file type.
var ValidInputTypes = []string{
	"png",  // std/image
	"gif",  // std/image
	"jpeg", // std/image
	"jpg",  // std/image

	"bmp",  // golang.org/x/image
	"tiff", // golang.org/x/image
	"vp8l", // golang.org/x/image
	"webp", // golang.org/x/image

	"pbm",      // github.com/jbuchbinder/gopnm
	"pgm",      // github.com/jbuchbinder/gopnm
	"ppm",      // github.com/jbuchbinder/gopnm
	"pcx",      // github.com/samuel/go-pcx/pcx
	"blp",      // github.com/nielsAD/gowarcraft3/file/blp
	"exr",      // github.com/mokiat/goexr/exr
	"megasd",   // github.com/bodgit/megasd/image
	"qoi",      // lelux.net/x/image/qoi
	"tga",      // github.com/blezek/tga
	"xcf",      // vimagination.zapto.org/limage/xcf
	"psd",      // github.com/oov/psd
	"ase",      // github.com/askeladdk/aseprite
	"aseprite", // github.com/askeladdk/aseprite

	"jfif", // self
}

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
		_ = f.Close()
		return fmt.Errorf("%s couldn't be opened, skipping (%s)", filename, err.Error())
	}

	inputValid := false
	for _, item := range ValidInputTypes {
		if filepath.Ext(filename) == "."+item {
			inputValid = true
		}
	}

	if !inputValid {
		_ = f.Close()
		return fmt.Errorf("%s is not a valid input, skipping", filename)
	}

	switch strings.ToLower(filepath.Ext(filename)) {
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
	case ".jfif":
		decodedImage, err = selfJfif.Decode(f)
	case ".pbm", ".pgm", ".ppm":
		decodedImage, err = pnm.Decode(f)
	case ".pcx":
		decodedImage, err = pcx.Decode(f)
	case ".blp":
		decodedImage, err = blp.Decode(f)
	case ".exr":
		decodedImage, err = exr.Decode(f)
	case ".megasd":
		decodedImage, err = megaSD.Decode(f)
	case ".qoi":
		decodedImage, err = qoi.Decode(f)
	case ".tga":
		decodedImage, err = tga.Decode(f)
	case ".xcf":
		decodedImage, err = xcf.Decode(f)
	case ".psd":
		var psdResult *psd.PSD
		psdResult, _, err = psd.Decode(f, nil)
		decodedImage = psdResult.Picker
	case ".ase", ".aseprite":
		decodedImage, err = aseprite.Decode(f)
	}

	if err != nil {
		_ = f.Close()
		return fmt.Errorf("%s couldn't be decoded (%s), skipping", filename, err.Error())
	}

	r, err := os.Create(filename + "." + outputFileType)
	if err != nil {
		_ = f.Close()
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
		err = webpEncoder.Encode(r, decodedImage, &webpEncoder.Options{
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
		_ = f.Close()
		_ = r.Close()
		return fmt.Errorf("%s.%s couldn't be encoded (%s)", filename, outputFileType, err.Error())
	}

	_ = f.Close()
	_ = r.Close()

	return nil
}
