package main

import (
	"fmt"
	imgui "github.com/AllenDang/cimgui-go"
	"sync"
)

var (
	exporterWaitGroup sync.WaitGroup

	backend         imgui.Backend
	specificBackend imgui.GLFWBackend
	windowFlags     imgui.GLFWWindowFlags

	skipSameType   = true
	overwriteFiles = true

	errorPopupName = "Whoops!"

	compiledErrors = ""
	toldErrors     = false

	showAbout bool
	showDemo  bool

	selectedFileType = 0

	lossless        = true
	losslessTooltip = "supported by: webp"

	qualityInt     int32 = 100
	qualityTooltip       = "supported by: jpeg, webp\n\n0%% is worst, 100%% is best"

	gifColors       int32 = 256
	tiffCompression int32 = 0

	//progress float32 = 0.0
)

const credit string = `credits n stuff
overall program by xubiod 2023
made with go 1.20.5

libraries used are from imports

------------ decoder libraries ------------
---------- (allowed file inputs) ----------
gif         => image/gif
jpeg/jpg    => image/jpeg
png         => image/png

bmp         => golang.org/x/image/bmp
tiff        => golang.org/x/image/tiff

webp        => github.com/chai2010/webp
pbm         => github.com/jbuchbinder/gopnm
pgm         => github.com/jbuchbinder/gopnm
ppm         => github.com/jbuchbinder/gopnm
pcx         => github.com/samuel/go-pcx/pcx

jfif        => internally (i wrote it!)

------------ encoder libraries ------------
---------- (allowed file output) ----------
gif         => image/gif
jpeg/jpg    => image/jpeg
png         => image/png

bmp         => golang.org/x/image/bmp
tiff        => golang.org/x/image/tiff
vp8l        => golang.org/x/image/vp8l
webp        => golang.org/x/image/webp

jfif        => github.com/leotaku/mobi/jfif
pbm         => github.com/jbuchbinder/gopnm
pgm         => github.com/jbuchbinder/gopnm
ppm         => github.com/jbuchbinder/gopnm
pcx         => github.com/samuel/go-pcx/pcx

------------- other libraries -------------
imgui       => github.com/AllenDang/cimgui-go`

var showCredit = false

var tiffCompressionNames = []string{
	"uncompressed",
	"deflate",
	"lzw",
	"ccittgroup3",
	"ccittgroup4",
}

func uiLoop() {
	showMiniWindow()
	showConfigurationWindow()
	if showCredit {
		showCreditWindow()
	}
}

var windowSize = imgui.Vec2{X: 800, Y: 800}
var configWindowSize = imgui.Vec2{X: 650, Y: 700}

func showConfigurationWindow() {
	// imgui.SetNextWindowPosV(imgui.NewVec2(0, 0), imgui.CondOnce, imgui.NewVec2(0, 0))

	imgui.SetNextWindowSizeV(configWindowSize, imgui.CondOnce)
	imgui.Begin("config")

	imgui.TextUnformatted("to convert, just drop it on the main window, files will appear in the same directory")
	imgui.TextUnformatted(fmt.Sprintf("currently supports decoding: %v", ValidInputTypes))

	imgui.NewLine()

	imgui.BeginListBoxV("convert to", imgui.Vec2{Y: 200})
	for i := 0; i < len(ValidOutputTypes); i++ {
		isSelected := selectedFileType == i
		if imgui.SelectableBoolPtr(ValidOutputTypes[i], &isSelected) {
			selectedFileType = i
		}

		if isSelected {
			imgui.SetItemDefaultFocus()
		}
	}
	imgui.EndListBox()

	imgui.NewLine()

	imgui.Checkbox("skip same types", &skipSameType)
	if imgui.IsItemHovered() {
		imgui.SetTooltip("if checked, the converter will skip loading and converting files\n" +
			"that are the same type as the output, i.e. skipping png => png\n\n" +
			"this can be unchecked if the behaviour is desired, but sizes and quality might change,\n" +
			"but the input files will not be overwritten due to how files are named")
	}

	imgui.SameLine()

	imgui.Checkbox("overwrite files", &overwriteFiles)
	if imgui.IsItemHovered() {
		imgui.SetTooltip("if checked, the converter will overwrite files that already exist,\n" +
			"given that they share the output filename of *.[old type].[new type]\n" +
			"i.e. overwriting \"file.png.jpg\" if converting \"file.png\" to jpg.\n\n" +
			"if unchecked, it will skip the file and report the skip when finished")
	}

	//if floating {
	//	windowFlags = windowFlags ^ imgui.GLFWWindowFlagsFloating
	//}

	if compiledErrors != "" && !toldErrors {
		imgui.OpenPopupStr(errorPopupName)
		toldErrors = true
	}

	if imgui.BeginPopupModalV(errorPopupName, nil, imgui.WindowFlagsNoResize|imgui.WindowFlagsAlwaysAutoResize) {
		imgui.TextUnformatted("There were some problems. Please review:")
		imgui.Separator()
		imgui.TextUnformatted(compiledErrors)

		if imgui.ButtonV("Acknowledge", imgui.Vec2{X: 120}) {
			compiledErrors = ""
			imgui.CloseCurrentPopup()
		}

		imgui.EndPopup()
	}

	imgui.NewLine()
	imgui.TextUnformatted("generic options (for those that support it):")

	imgui.Checkbox("lossless", &lossless)
	if imgui.IsItemHovered() {
		imgui.SetTooltip(losslessTooltip)
	}

	imgui.SliderInt("quality", &qualityInt, 0, 100)
	if imgui.IsItemHovered() {
		imgui.SetTooltip(qualityTooltip)
	}

	imgui.NewLine()
	imgui.TextUnformatted("specific options:")
	imgui.SliderInt("gif colors", &gifColors, 1, 256)
	if imgui.IsItemHovered() {
		imgui.SetTooltip("changes the amount of allowed colors for gif\n\nthere's only a maximum of 256 available on gif's\npalette table, but it can be less")
	}

	imgui.BeginListBoxV("tiff compression type", imgui.Vec2{Y: 87})
	var i int32 = 0
	for i = 0; i < int32(len(tiffCompressionNames)); i++ {
		isSelected := tiffCompression == i
		if imgui.SelectableBoolPtr(tiffCompressionNames[i], &isSelected) {
			tiffCompression = i
		}

		if isSelected {
			imgui.SetItemDefaultFocus()
		}
	}
	imgui.EndListBox()
	//imgui.ProgressBar(progress)
	if imgui.Button("credits") {
		showCredit = !showCredit
	}

	imgui.End()
}

func showCreditWindow() {
	imgui.SetNextWindowSizeV(imgui.NewVec2(500, 400), imgui.CondOnce)
	imgui.Begin("credit")
	imgui.TextWrapped(credit)
	imgui.End()
}

func showMiniWindow() {
	imgui.Begin("deboog")
	imgui.Checkbox("imgui about", &showAbout)
	imgui.Checkbox("imgui demo", &showDemo)

	if showAbout {
		imgui.ShowAboutWindow()
	}

	if showDemo {
		imgui.ShowDemoWindow()
	}

	imgui.End()
}

func ui() {
	//var err error

	specificBackend = *imgui.NewGLFWBackend()
	backend = imgui.CreateBackend(&specificBackend)

	backend.SetBgColor(imgui.NewVec4(0.45, .55, .6, 1.0))
	backend.CreateWindow("title", int(windowSize.X), int(windowSize.Y), windowFlags)

	backend.SetDropCallback(func(p []string) {
		fmt.Printf("drop: %v", p)

		var genericQuality = 0

		switch ValidOutputTypes[selectedFileType] {
		case "tiff":
			genericQuality = int(tiffCompression)
		case "gif":
			genericQuality = int(gifColors)
		case "jpeg", "jpg", "jfif":
			genericQuality = int(qualityInt)
		}

		compiledErrors = ""
		toldErrors = true
		//progress = 0.0
		for i := range p {
			exporterWaitGroup.Add(1)
			go func(idx int) {
				defer exporterWaitGroup.Done()
				err := ConvertTo(p[idx], ValidOutputTypes[selectedFileType], QualityInformation{
					QualityInt:   genericQuality,
					QualityFloat: float32(qualityInt),
				}, !skipSameType, overwriteFiles)
				//progress = float32(i+1) / float32(len(p))

				if err != nil {
					compiledErrors += err.Error() + "\n"
				}
			}(i)
		}

		exporterWaitGroup.Wait()

		if compiledErrors != "" {
			toldErrors = false
		}
	})

	imgui.StyleColorsClassic()

	backend.Run(uiLoop)
}

func main() {
	ui()
}
