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
	fps             uint = 45

	skipSameType   = true
	overwriteFiles = true

	errorPopupName = "Whoops!"

	compiledErrors = ""
	toldErrors     = false

	showAbout  bool
	showDemo   bool
	showCredit bool
	showMini   bool

	selectedFileType = 0

	opts = SettingsSavable{
		SkipSameType:   true,
		OverwriteFiles: true,

		Lossless: true,
		Quality:  100,

		WebpExact:       true,
		GifColors:       256,
		TiffCompression: 0,
		TiffPredictor:   false,
	}

	losslessTooltip = "supported by: webp"

	exact        = true
	exactTooltip = "preserves RGB values in transparent areas\nif off, rgba(255, 127, 0, 0.0) would (probably) become rgba(0, 0, 0, 0.0)"

	qualityInt     int32 = 100
	qualityTooltip       = "supported by: jpeg, webp, jfif\n\n0%% is worst, 100%% is best"

	gifColors       int32 = 256
	tiffCompression int32 = 0
	tiffPredictor         = false
)

const credit string = `credits n stuff
overall program by xubiod 2023
made with go 1.20.5

libraries used are from imports

------- encoder + decoder libraries -------
----- (does both in the same library) -----
gif         => image/gif
jpeg/jpg    => image/jpeg
png         => image/png

bmp         => golang.org/x/image/bmp
tiff        => golang.org/x/image/tiff

pbm         => github.com/jbuchbinder/gopnm
pgm         => github.com/jbuchbinder/gopnm
ppm         => github.com/jbuchbinder/gopnm
pcx         => github.com/samuel/go-pcx/pcx
megasd      => github.com/bodgit/megasd/image
qoi         => lelux.net/x/image/qoi
tga         => github.com/blezek/tga
xcf         => vimagination.zapto.org/limage/xcf

------------ encoder libraries ------------
---------- (allowed file output) ----------
webp        => github.com/chai2010/webp
blp         => github.com/nielsAD/gowarcraft3/file/blp
exr         => github.com/mokiat/goexr/exr
xpm         => github.com/xyproto/xpm

jfif        => internally (i wrote it!)

------------ decoder libraries ------------
---------- (allowed file inputs) ----------
vp8l        => golang.org/x/image/vp8l
webp        => golang.org/x/image/webp

jfif        => github.com/leotaku/mobi/jfif
psd         => github.com/oov/psd
ase         => github.com/askeladdk/aseprite
ico         => github.com/mat/besticon/ico

------------- other libraries -------------
imgui       => github.com/AllenDang/cimgui-go`

var typeExplainer = map[string]string{
	"png": "as specified by: https://www.w3.org/TR/PNG/",
	"gif": "as specified by: https://www.w3.org/Graphics/GIF/spec-gif89a.txt\n\n" +
		"quantizes with go's palette.Plan9 palette, and draws with go's default\n" +
		"draw.FloydSteinburg. rest of the options are given to you to modify",
	"jpeg": "as specified by: https://www.w3.org/Graphics/JPEG/itu-t81.pdf\n\n" +
		"all available options are given to you to modify",

	"bmp": "as specified by: https://www.digicamsoft.com/bmp/bmp.html",
	"tiff": "as specified by: https://partners.adobe.com/public/developer/en/tiff/TIFF6.pdf\n\n" +
		"all available options are given to you to modify",

	"jfif": "package doesn't list a spec\n\n" +
		"the jfif header is baked into the package used, and can't be modified\n" +
		"all available options are given to you to modify (compat. w/ jpeg)",
	"webp": "as specified by: https://developers.google.com/speed/webp/docs/riff_container\n\n" +
		"all available options are given to you to modify",
	"pbm": "package doesn't list a spec\n\n" +
		"PBM makes 1-bit bitmaps (only black at 0 and white at 255)",
	"pgm": "package doesn't list a spec\n\n" +
		"PGM makes grayscale bitmaps",
	"ppm": "package doesn't list a spec\n\n" +
		"PPM makes color bitmaps",
	"pcx": "as specified by: https://web.archive.org/web/20030111010058/http://www.nist.fss.ru/hr/doc/spec/pcx.htm",
	"megasd": "package defines its own spec:\n\n" +
		"\"The format is defined as 64 by 40 pixels exactly which is split into forty 8 by 8 tiles.\n" +
		"Up to three 16 color palettes can be defined and each tile can use only one of these palettes.\n" +
		"The first color in each palette is reserved for transparency.\"",
	"qoi": "as specified by: https://qoiformat.org/qoi-specification.pdf",
	"tga": "package doesn't list a spec, but passes TGA 2.0 conformance\n" +
		"which is available at https://googlesites.inequation.org/tgautilities",
	"xpm": "package says \"X PixMap (XPM3)\"",
	"xcf": "package says \"GIMP's XCF format\"",
}

var tiffCompressionNames = []string{
	"uncompressed",
	"deflate",
	"lzw",
	"ccittgroup3",
	"ccittgroup4",
}

func uiLoop() {
	if showMini {
		showMiniWindow()
	}

	showConfigurationWindow()
	if showCredit {
		showCreditWindow()
	}
}

var windowSize = imgui.Vec2{X: 800, Y: 800}
var configWindowSize = imgui.Vec2{X: 650, Y: 700}

func showConfigurationWindow() {
	imgui.SetNextWindowSizeV(configWindowSize, imgui.CondOnce)
	imgui.Begin("config")

	imgui.TextUnformatted("to convert, just drop it on the main window, files will appear in the same directory")
	imgui.TextWrapped(fmt.Sprintf("currently supports decoding: %v", ValidInputTypes))
	imgui.TextWrapped("layered formats will get flattened, and animations *might* become texture aliases, but it depends on the format and decoder")

	imgui.NewLine()

	imgui.BeginListBoxV("convert to", imgui.Vec2{Y: 200})
	for i := 0; i < len(ValidOutputTypes); i++ {
		isSelected := selectedFileType == i
		if imgui.SelectableBoolPtr(ValidOutputTypes[i], &isSelected) {
			selectedFileType = i
		}

		if _, ok := typeExplainer[ValidOutputTypes[i]]; ok && imgui.IsItemHovered() {
			imgui.SetTooltip(typeExplainer[ValidOutputTypes[i]])
		}

		if isSelected {
			imgui.SetItemDefaultFocus()
		}
	}
	imgui.EndListBox()

	imgui.NewLine()

	imgui.Checkbox("skip same types", &opts.SkipSameType)
	if imgui.IsItemHovered() {
		imgui.SetTooltip("if checked, the converter will skip loading and converting files\n" +
			"that are the same type as the output, i.e. skipping png => png\n\n" +
			"this can be unchecked if the behaviour is desired, but sizes and quality might change,\n" +
			"(as the raw pixel data is getting unpacked and repacked) but the input files will not\n" +
			"be overwritten due to how files are named")
	}

	imgui.SameLine()

	imgui.Checkbox("overwrite files", &opts.OverwriteFiles)
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

	if imgui.CollapsingHeaderTreeNodeFlags("export options (defaults to highest quality per method)") {
		imgui.TextUnformatted("generic options (for those that support it):")

		imgui.Checkbox("lossless", &opts.Lossless)
		if imgui.IsItemHovered() {
			imgui.SetTooltip(losslessTooltip)
		}

		imgui.SliderInt("quality", &opts.Quality, 0, 100)
		if imgui.IsItemHovered() {
			imgui.SetTooltip(qualityTooltip)
		}

		imgui.NewLine()
		imgui.TextUnformatted("specific options:")

		switch ValidOutputTypes[selectedFileType] {
		case "gif":
			imgui.SliderInt("gif colors", &opts.GifColors, 1, 256)
			if imgui.IsItemHovered() {
				imgui.SetTooltip("changes the amount of allowed colors for gif\n\nthere's only a maximum of 256 available on gif's\npalette table, but it can be less")
			}

		case "tiff":
			imgui.BeginListBoxV("tiff compression type", imgui.Vec2{Y: 87})
			var i int32 = 0
			for i = 0; i < int32(len(tiffCompressionNames)); i++ {
				isSelected := opts.TiffCompression == i
				if imgui.SelectableBoolPtr(tiffCompressionNames[i], &isSelected) {
					opts.TiffCompression = i
				}

				if isSelected {
					imgui.SetItemDefaultFocus()
				}
			}
			imgui.EndListBox()

			imgui.Checkbox("tiff predictor", &opts.TiffPredictor)
			if imgui.IsItemHovered() {
				imgui.SetTooltip("determines whether a differencing predictor is used\nit can improve the compression in certain situations")
			}

		case "webp":
			imgui.Checkbox("webp exact", &opts.WebpExact)
			if imgui.IsItemHovered() {
				imgui.SetTooltip(exactTooltip)
			}

		default:
			imgui.TextUnformatted(ValidOutputTypes[selectedFileType] + " doesn't have specific options")
		}
	}

	imgui.NewLine()

	//imgui.ProgressBar(progress)
	if imgui.Button("credits") {
		showCredit = !showCredit
	}
	imgui.SameLine()
	if imgui.Button("imgui builtin") {
		showMini = !showMini
	}

	imgui.End()
}

func showCreditWindow() {
	imgui.SetNextWindowSizeV(imgui.NewVec2(500, 400), imgui.CondOnce)
	imgui.BeginV("credit", &showCredit, imgui.WindowFlagsNone)
	imgui.TextWrapped(credit)
	imgui.End()
}

var showMetrics = false
var showUserGuide = false
var showDebugLog = false
var showStackTool = false
var showStyleEdit = false

func showMiniWindow() {
	imgui.BeginV("imgui builtin", &showMini, imgui.WindowFlagsNone)
	imgui.Checkbox("about", &showAbout)
	imgui.Checkbox("demo", &showDemo)
	imgui.Checkbox("debuglog", &showDebugLog)
	imgui.Checkbox("metrics", &showMetrics)
	imgui.Checkbox("stacktool", &showStackTool)
	imgui.Checkbox("styleedit", &showStyleEdit)
	imgui.Checkbox("userguide", &showUserGuide)

	if showAbout {
		imgui.ShowAboutWindowV(&showAbout)
	}

	if showDemo {
		imgui.ShowDemoWindowV(&showDemo)
	}

	if showDebugLog {
		imgui.ShowDebugLogWindowV(&showDebugLog)
	}

	if showMetrics {
		imgui.ShowMetricsWindowV(&showMetrics)
	}

	if showStackTool {
		imgui.ShowStackToolWindowV(&showStackTool)
	}

	if showStyleEdit {
		imgui.ShowStyleEditor()
	}

	if showUserGuide {
		imgui.ShowUserGuide()
	}

	imgui.End()
}

func dropOn(p []string) {
	fmt.Printf("drop: %v\n", p)

	var genericQuality = 0

	switch ValidOutputTypes[selectedFileType] {
	case "tiff":
		genericQuality = int(opts.TiffCompression)
	case "gif":
		genericQuality = int(opts.GifColors)
	case "jpeg", "jpg", "jfif":
		genericQuality = int(opts.Quality)
	}

	compiledErrors = ""
	toldErrors = true
	for i := range p {
		exporterWaitGroup.Add(1)
		go func(idx int) {
			defer exporterWaitGroup.Done()
			err := ConvertTo(p[idx], ValidOutputTypes[selectedFileType], QualityInformation{
				QualityInt:    genericQuality,
				QualityFloat:  float32(qualityInt),
				TiffPredictor: tiffPredictor,
				WebpExact:     exact,
			}, !skipSameType, overwriteFiles)

			if err != nil {
				compiledErrors += err.Error() + "\n"
			}
		}(i)
	}

	exporterWaitGroup.Wait()

	if compiledErrors != "" {
		toldErrors = false
	}
}

func ui() {
	//var err error

	specificBackend = *imgui.NewGLFWBackend()
	backend = imgui.CreateBackend(&specificBackend)

	backend.SetBgColor(imgui.NewVec4(0.45, .55, .6, 1.0))
	backend.CreateWindow("img-convert - dropzone", int(windowSize.X), int(windowSize.Y), windowFlags)

	backend.SetDropCallback(dropOn)

	imgui.StyleColorsClassic()
	backend.SetTargetFPS(fps)

	backend.Run(uiLoop)
}

func main() {
	ui()
}
