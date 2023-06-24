package main

//
//import (
//	"fmt"
//	"os"
//)
//
//func cli() {
//	if len(os.Args) < 2 {
//		help()
//		return
//	}
//
//	output_type := os.Args[1]
//	//var quality = ""
//	startIndex := 2
//
//	if !(output_type == "png" || output_type == "bmp") {
//		quality = os.Args[2]
//		startIndex++
//	}
//
//	valid := false
//	for _, item := range ValidOutputTypes {
//		if output_type == item {
//			valid = true
//		}
//	}
//
//	if !valid {
//		fmt.Printf("output format not supported\nvalid types: %#v", ValidOutputTypes)
//		return
//	}
//
//	//for _, arg := range os.Args[startIndex:] {
//	//	i, _ := strconv.Atoi(quality)
//	//	// ConvertTo(arg, output_type, i)
//	//}
//	fmt.Println("job completed")
//}
//
//func help() {
//	fmt.Printf("img-convert [output-type] [quality?] path1 path2...\n"+
//		"a simple utility to mass convert images into one type\n\n"+
//		"supported inputs: %#v\n"+
//		"NOTE: GIFS ARE TREATED AS IMAGES, NOT ANIMATIONS\n\n"+
//		"output-type:\tone of the following: %#v\n"+
//		"quality?:\tif tiff, gif, or jpeg, must be set, otherwise omitted\n"+
//		"\t\ttiff: follows golang.org/x/image/tiff system:\n"+
//		"\t\t\t0 -> uncompressed\n"+
//		"\t\t\t1 -> deflate\n"+
//		"\t\t\t2 -> LZW\n"+
//		"\t\t\t3 -> CCITTGroup3\n"+
//		"\t\t\t4 -> CCITTGroup4\n"+
//		"\t\tgif: number of colors, 1 to 256, Plan9 is always the quantizer, FloydSteinberg is always the drawer\n"+
//		"\t\tjpeg: quality%% (0%% worst, 100%% best)\n\n", ValidInputTypes, ValidOutputTypes)
//}
