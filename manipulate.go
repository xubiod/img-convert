package main

import (
	"image"

	"github.com/disintegration/imaging"
)

type FilterType int

const (
	FilterUpscale FilterType = iota
	FilterFlipH
	FilterFlipV
	FilterRot90
	FilterRot180
	FilterRot270
)

var FilterNames []string = []string{
	"1:1 upscale",
	"flip horizontal",
	"flip vertical",
	"rotate 90 deg (CW)",
	"rotate 180 deg (CW)",
	"rotate 270 deg (CW)",
}

type Filter struct {
	What FilterType

	IntFactor int32
	Resample  imaging.ResampleFilter
}

var Filters []*Filter = make([]*Filter, 0)

func ApplyFilters(img image.Image) image.Image {
	var dst image.Image = img

	for _, value := range Filters {
		switch value.What {
		case FilterUpscale:
			dst = imaging.Resize(dst, dst.Bounds().Dx()*int(value.IntFactor), dst.Bounds().Dy()*int(value.IntFactor), value.Resample)
		case FilterFlipH:
			dst = imaging.FlipH(img)
		case FilterFlipV:
			dst = imaging.FlipV(img)
		case FilterRot90:
			dst = imaging.Rotate90(img)
		case FilterRot180:
			dst = imaging.Rotate180(img)
		case FilterRot270:
			dst = imaging.Rotate270(img)
		}
	}

	return dst
}
