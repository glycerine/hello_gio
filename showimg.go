package main

// A simple Gio program. See https://gioui.org for more information.

import (
	"bufio"
	"image"
	"image/color"

	"gioui.org/app"
	"gioui.org/f32"
	//"gioui.org/gesture"
	"gioui.org/io/pointer"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/paint"
	//"gioui.org/text"
	//"gioui.org/text/shape"
	//"gioui.org/unit"
	"gioui.org/widget"
	//"golang.org/x/exp/shiny/iconvg"
	//"golang.org/x/image/draw"

	"golang.org/x/exp/shiny/materialdesign/icons"
	//"golang.org/x/image/font/gofont/goregular"
	//"golang.org/x/image/font/gofont/gomonobold"
	//"golang.org/x/image/font/gofont/gomono"
	//"golang.org/x/image/font/sfnt"

	//"fmt"
	//"io"
	"os"
	//"path"
	//"runtime"
	//"sync"
	//"time"

	"4d63.com/tz"

	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
)

var _ = pointer.RectAreaOp{}
var _ = widget.Image{}
var _ = tz.LoadLocation
var _ = icons.NavigationArrowBack

var colors = make(map[string]color.RGBA)

func init() {
	colors["maroon"] = color.RGBA{127, 0, 0, 255}
	colors["cream"] = color.RGBA{240, 240, 127, 255}
	colors["black"] = color.RGBA{0, 0, 0, 255}
}

type myDrawState struct {
	w   *app.Window
	gtx *layout.Context

	pngPlot     image.Image
	pngPlotRect image.Rectangle
}

func showImageMain() {

	go func() {
		w := app.NewWindow()

		var err error
		m := &myDrawState{
			w: w,
		}
		m.gtx = &layout.Context{
			Queue: w.Queue(),
		}

		m.pngPlot, _, err = LoadImage("points.png")
		panicOn(err)
		m.pngPlotRect = m.pngPlot.(*image.NRGBA).Rect
		//vv("m.pngPlot.Rect = '%#v'", m.pngPlotRect)

		var mainErr error
		yellowBkg := true // show image on background field of yellow

	mainLoop:
		for {
			e := <-w.Events()
			switch e := e.(type) {
			case app.DestroyEvent:
				mainErr = e.Err
				break mainLoop
			case app.UpdateEvent:
				//vv("e is '%#v'", e)
				showImage(e, m, yellowBkg)
			}
		}
		panicOn(mainErr)
	}()
	app.Main()
}

func showImage(e app.UpdateEvent, m *myDrawState, yellowBkg bool) {
	gtx := m.gtx
	ops := &op.Ops{}
	m.gtx.Ops = ops
	w := m.w
	gtx.Reset(&e.Config, e.Size)

	width := 1000
	height := 600
	x0 := 300
	x1 := x0 + width
	y0 := 200
	y1 := y0 + height
	_ = x1
	_ = y1

	ops.Reset()
	if yellowBkg {
		// lets us see easily where the image frame is.
		paint.ColorOp{Color: colors["cream"]}.Add(ops)
		paint.PaintOp{Rect: f32.Rectangle{Max: f32.Point{X: float32(e.Size.X), Y: float32(e.Size.Y)}}}.Add(ops)
	}

	// show the image/texture.
	// there might be a more efficient way? Ask around.
	// How do we proportionately scale down the image to fit in a given box?

	imgSrc := m.pngPlot
	imgPos := image.Rectangle{
		Min: image.Point{X: x0, Y: y0},
		Max: image.Point{X: x1, Y: y1},
	}
	paint.ImageOp{Src: imgSrc, Rect: imgPos}.Add(ops) // display the png, part 1
	paint.PaintOp{Rect: toRectF(imgPos)}.Add(ops)     // display the png, part 2

	w.Update(ops)
}

func LoadImage(filename string) (image.Image, string, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, "", err
	}
	defer f.Close()
	return image.Decode(bufio.NewReader(f))
}

func toRectF(r image.Rectangle) f32.Rectangle {
	return f32.Rectangle{
		Min: f32.Point{X: float32(r.Min.X), Y: float32(r.Min.Y)},
		Max: f32.Point{X: float32(r.Max.X), Y: float32(r.Max.Y)},
	}
}
