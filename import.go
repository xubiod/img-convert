package main

import (
	"fmt"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"os"
	"path/filepath"
	"strings"

	"github.com/askeladdk/aseprite"
	"github.com/blezek/tga"
	megaSD "github.com/bodgit/megasd/image"
	"github.com/hhrutter/tiff"
	pnm "github.com/jbuchbinder/gopnm"
	"github.com/mat/besticon/ico"
	"github.com/mokiat/goexr/exr"
	"github.com/nielsAD/gowarcraft3/file/blp"
	"github.com/oov/psd"
	"github.com/samuel/go-pcx/pcx"
	selfJfif "github.com/xubiod/img-convert/dedicated-decoder/jfif"
	"golang.org/x/image/bmp"
	"golang.org/x/image/vp8l"
	"lelux.net/x/image/qoi"
	"lelux.net/x/image/webp"
	"vimagination.zapto.org/limage/xcf"
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
	"ico",      // github.com/mat/besticon/ico

	"jfif", // self
}

func Import(filename string) (decodedImage image.Image, err error) {
	f, err := os.Open(filename)
	if err != nil {
		_ = f.Close()
		return decodedImage, fmt.Errorf("%s couldn't be opened, skipping (%s)", filename, err.Error())
	}

	inputValid := false
	for _, item := range ValidInputTypes {
		if filepath.Ext(filename) == "."+item {
			inputValid = true
		}
	}

	if !inputValid {
		_ = f.Close()
		return decodedImage, fmt.Errorf("%s is not a valid input, skipping", filename)
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
	case ".ico":
		decodedImage, err = ico.Decode(f)
	}

	if err != nil {
		_ = f.Close()
		return decodedImage, fmt.Errorf("%s couldn't be decoded (%s), skipping", filename, err.Error())
	}

	_ = f.Close()

	return decodedImage, nil
}
