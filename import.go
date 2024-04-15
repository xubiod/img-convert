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
	"github.com/sugyan/ttygif/image/xwd"
	selfJfif "github.com/xubiod/img-convert/dedicated-decoder/jfif"
	"golang.org/x/image/bmp"
	"golang.org/x/image/vp8"
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

	switch strings.ToLower(filepath.Ext(filename)) {
	case ".bmp":
		imported, err = bmp.Decode(f)

	case ".tiff":
		imported, err = tiff.Decode(f)

	case ".vp8l":
		imported, err = vp8l.Decode(f)

	case ".vp8":
		var decoder *vp8.Decoder = vp8.NewDecoder()
		fi, _ := f.Stat()
		decoder.Init(f, int(fi.Size()))
		imported, err = decoder.DecodeFrame()

	case ".webp":
		imported, err = webp.Decode(f)

	case ".gif":
		imported, err = gif.Decode(f)

	case ".jpg", ".jpeg":
		imported, err = jpeg.Decode(f)

	case ".png":
		imported, err = png.Decode(f)

	case ".jfif":
		imported, err = selfJfif.Decode(f)

	case ".pbm", ".pgm", ".ppm":
		imported, err = pnm.Decode(f)

	case ".pcx":
		imported, err = pcx.Decode(f)

	case ".blp":
		imported, err = blp.Decode(f)

	case ".exr":
		imported, err = exr.Decode(f)

	case ".megasd":
		imported, err = megaSD.Decode(f)

	case ".qoi":
		imported, err = qoi.Decode(f)

	case ".tga":
		imported, err = tga.Decode(f)

	case ".xcf":
		imported, err = xcf.Decode(f)

	case ".psd":
		var psdResult *psd.PSD
		psdResult, _, err = psd.Decode(f, nil)
		imported = psdResult.Picker

	case ".ase", ".aseprite":
		imported, err = aseprite.Decode(f)

	case ".ico":
		imported, err = ico.Decode(f)

	case ".xwd":
		imported, err = xwd.Decode(f)
	}

	if err != nil {
		_ = f.Close()
		return imported, fmt.Errorf("%s couldn't be decoded (%s), skipping", filename, err.Error())
	}

	_ = f.Close()

	return imported, nil
}
