package main

// A simple Gio program. See https://gioui.org

import (
	"bufio"
	"image"
	"image/color"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"os"

	"gioui.org/app"
	"gioui.org/f32"
	"gioui.org/layout"
	"gioui.org/op/paint"
)

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
		vv("m.pngPlot.Rect = '%#v'", m.pngPlotRect)

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
	m.gtx.Reset(&e.Config, e.Size)
	ops := m.gtx.Ops

	// choose where to show the image
	width := 1000
	height := 600
	x0 := 300
	x1 := x0 + width
	y0 := 200
	y1 := y0 + height
	imgPos := image.Rectangle{
		Min: image.Point{X: x0, Y: y0},
		Max: image.Point{X: x1, Y: y1},
	}

	// get full window coordinates to paint the background
	fullWindowRect := image.Rectangle{Max: image.Point{X: e.Size.X, Y: e.Size.Y}}
	fullWindowRect32 := toRectF(fullWindowRect)
	if yellowBkg {
		// lets us see easily where the image frame is.
		paint.ColorOp{Color: colors["cream"]}.Add(ops)
		paint.PaintOp{Rect: fullWindowRect32}.Add(ops)
	}

	// show the image/texture.
	paint.ImageOp{Src: m.pngPlot, Rect: m.pngPlotRect}.Add(ops) // display the png, part 1
	paint.PaintOp{Rect: toRectF(imgPos)}.Add(ops)               // display the png, part 2

	m.w.Update(ops)
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
