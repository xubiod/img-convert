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
	"github.com/sugyan/ttygif/image/xwd"
	selfJfif "github.com/xubiod/img-convert/dedicated-decoder/jfif"
	"golang.org/x/image/bmp"
	"golang.org/x/image/vp8"
	"golang.org/x/image/vp8l"
	"lelux.net/x/image/qoi"
	"lelux.net/x/image/webp"
	"vimagination.zapto.org/limage/xcf"

	"github.com/gen2brain/jpegxl"
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
	"vp8",  // golang.org/x/image
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
	"xwd",      // github.com/sugyan/ttygif/image/xwd

	"jfif", // self
}

var genericImporters = map[string]func(io.Reader) (image.Image, error){
	".bmp":      bmp.Decode,
	".tiff":     tiff.Decode,
	".vp8l":     vp8l.Decode,
	".webp":     webp.Decode,
	".gif":      gif.Decode,
	".jpg":      jpeg.Decode,
	".jpeg":     jpeg.Decode,
	".jxl":      jpegxl.Decode,
	".png":      png.Decode,
	".jfif":     selfJfif.Decode,
	".pbm":      pnm.Decode,
	".pgm":      pnm.Decode,
	".ppm":      pnm.Decode,
	".pcx":      pcx.Decode,
	".blp":      blp.Decode,
	".exr":      exr.Decode,
	".megasd":   megaSD.Decode,
	".qoi":      qoi.Decode,
	".tga":      tga.Decode,
	".ase":      aseprite.Decode,
	".aseprite": aseprite.Decode,
	".ico":      ico.Decode,
	".xwd":      xwd.Decode,
}

func Import(filename string) (imported image.Image, err error) {
	f, err := os.Open(filename)
	if err != nil {
		_ = f.Close()
		return imported, fmt.Errorf("%s couldn't be opened, skipping (%s)", filename, err.Error())
	}

	inputValid := false
	for _, item := range ValidInputTypes {
		if filepath.Ext(filename) == "."+item {
			inputValid = true
		}
	}

	if !inputValid {
		_ = f.Close()
		return imported, fmt.Errorf("%s is not a valid input, skipping", filename)
	}

	ext := strings.ToLower(filepath.Ext(filename))

	switch ext {
	case ".vp8":
		var decoder *vp8.Decoder = vp8.NewDecoder()
		fi, _ := f.Stat()
		decoder.Init(f, int(fi.Size()))
		imported, err = decoder.DecodeFrame()

	case ".xcf":
		imported, err = xcf.Decode(f)

	case ".psd":
		var psdResult *psd.PSD
		psdResult, _, err = psd.Decode(f, nil)
		imported = psdResult.Picker

	_:
		imported, err = genericImporters[ext](f)
	}

	if err != nil {
		_ = f.Close()
		return imported, fmt.Errorf("%s couldn't be decoded (%s), skipping", filename, err.Error())
	}

	_ = f.Close()

	return imported, nil
}
