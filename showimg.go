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

func setupDrawState(w *app.Window) *myDrawState {
	m := &myDrawState{
		w: w,
	}
	m.gtx = &layout.Context{
		Queue: w.Queue(),
	}
	var err error
	m.pngPlot, _, err = LoadImage("points.png")
	panicOn(err)
	m.pngPlotRect = m.pngPlot.(*image.NRGBA).Rect
	vv("m.pngPlot.Rect = '%#v'", m.pngPlotRect)
	return m
}

func showImageMain() {

	go func() {
		w := app.NewWindow()

		var err error
		m := setupDrawState(w)

		yellowBkg := false // show image on background field of yellow?

	mainLoop:
		for {
			e := <-w.Events()
			switch e := e.(type) {
			case app.DestroyEvent:
				err = e.Err
				break mainLoop
			case app.UpdateEvent:
				//vv("e is '%#v'", e)
				showImage(e, m, yellowBkg)
			}
		}
		panicOn(err)
	}()
	app.Main()
}

func showImage(e app.UpdateEvent, m *myDrawState, yellowBkg bool) {
	m.gtx.Reset(&e.Config, e.Size)
	ops := m.gtx.Ops

	// choose how big to show the png
	width := 1000
	// choose height to maintain aspect ratio
	yx := float64(m.pngPlotRect.Max.Y) / float64(m.pngPlotRect.Max.X)
	height := int(float64(width) * yx)

	// choose where to place the png
	x0 := 300
	x1 := x0 + width
	y0 := 200
	y1 := y0 + height
	imgPos := image.Rectangle{
		Min: image.Point{X: x0, Y: y0},
		Max: image.Point{X: x1, Y: y1},
	}
	borderPx := 5 // pixel width of border
	borderRect := image.Rectangle{
		Min: image.Point{X: x0 - borderPx, Y: y0 - borderPx},
		Max: image.Point{X: x1 + borderPx, Y: y1 + borderPx},
	}

	// Get full window rectangle in order to paint the background.
	fullWindowRect := image.Rectangle{Max: image.Point{X: e.Size.X, Y: e.Size.Y}}
	fullWindowRect32 := toRectF(fullWindowRect)
	if yellowBkg {
		// lets us see easily where the png limits are.
		paint.ColorOp{Color: colors["cream"]}.Add(ops)
		paint.PaintOp{Rect: fullWindowRect32}.Add(ops)
	} else {
		// just paint a border
		paint.ColorOp{Color: colors["cream"]}.Add(ops)
		paint.PaintOp{Rect: toRectF(borderRect)}.Add(ops)
	}

	// Show the png image.
	// Here we show all of it, but we could just show a subset by
	// changing Imageop.Rect.

	// The ImageOp.Rect specifies the source rectangle.
	// The PaintOp.Rect field specifies the destination rectangle.
	// Scale the PaintOp.Rect to change the size of the rendered png.
	paint.ImageOp{Src: m.pngPlot, Rect: m.pngPlotRect}.Add(ops) // set the source for the png.
	paint.PaintOp{Rect: toRectF(imgPos)}.Add(ops)               // set the destination.

	//m.w.Update(ops)
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
