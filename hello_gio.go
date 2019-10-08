// SPDX-License-Identifier: Unlicense OR MIT
//
// hello_gio.go is my hello world program for the
// Gio graphics package for Golang by Elias Naur.
// https://gioui.org
//
// Eschewing the elegant constraint layout system for
// more direct control, hello_gio.go demonstrates how to
// plot rectangular boxes at specific screen positions of your choosing.
//
// It then demonstrates how to place labels over those boxes, and to clip
// those labels to stay within their boxes.
//
package main

import (
	"fmt"
	"image"
	"image/color"
	"log"

	"gioui.org/app"
	"gioui.org/f32"
	"gioui.org/io/pointer"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/paint"
	"gioui.org/text"
	"gioui.org/text/shape"
	"gioui.org/unit"
	"golang.org/x/image/font/gofont/goregular"
	"golang.org/x/image/font/sfnt"
)

var _ = paint.ImageOp{}
var _ = op.TransformOp{}
var _ image.Image
var _ = pointer.Event{}
var _ = fmt.Printf

func main() {
	// temporary: take over main() to get texture/png display working. See showimg.go.
	showImageMain()
	return

	go func() {
		w := app.NewWindow()
		if err := loop(w); err != nil {
			log.Fatal(err)
		}
	}()
	app.Main()
}

func loop(w *app.Window) error {
	regular, err := sfnt.Parse(goregular.TTF)
	if err != nil {
		panic("failed to load font")
	}
	_ = regular
	var faces shape.Faces
	face := faces.For(regular, unit.Sp(20))
	gtx := &layout.Context{
		Queue: w.Queue(),
	}

	// load image once
	//	m.pngPlot, _, err = LoadImage("points.png")
	//	panicOn(err)
	//	m.pngPlotRect = m.pngPlot.(*image.NRGBA).Rect
	//	vv("m.pngPlot.Rect = '%#v'", m.pngPlotRect)

	for {
		e := <-w.Events()
		switch e := e.(type) {
		case app.DestroyEvent:
			return e.Err
		case app.UpdateEvent:
			gtx.Reset(&e.Config, e.Size)
			faces.Reset(gtx.Config)
			direct(gtx, w, e, face)
		}
	}
}

func direct(gtx *layout.Context, w *app.Window, e app.UpdateEvent, face text.Face) {
	ops := gtx.Ops
	ops.Reset()
	aqua := color.RGBA{A: 0xff, G: 0xcc, B: 200}
	_ = aqua
	const borderPix = 5

	m := &box{
		h:         50,
		w:         50,
		color:     aqua,
		borderPix: borderPix,
		face:      face,
	}

	// draws 5 squares
	for i := 0; i < 5; i++ {
		x := 100 + i*100
		y := 50
		if i%2 == 0 {
			y = 0
		}
		ci := 50 * byte(i) // color increment

		// add _0123 to the end so we can see the clipping in action.
		m.drawc(gtx, x, y, color.RGBA{A: 0xff, G: 0xcc, B: ci, R: 255 - ci}, fmt.Sprintf("%v_0123", i))
	}

	// Submit operations to the window.
	w.Update(ops)
}

type box struct {
	h         int        //height
	w         int        //width
	color     color.RGBA // default
	borderPix int        // number of pixels to inset the label from the box edge.
	face      text.Face  // for label
}

// draw a rectangle at x0,y0 with given color; adding it to the gtx.Ops chain.
func (e *box) drawc(gtx *layout.Context, x0, y0 int, color color.RGBA, label string) {
	ops := gtx.Ops
	paint.ColorOp{Color: color}.Add(ops)
	re := f32.Rectangle{
		Min: f32.Point{X: float32(x0), Y: float32(y0)},
		Max: f32.Point{X: float32(x0 + e.w), Y: float32(y0 + e.h)},
	}
	paint.PaintOp{Rect: re}.Add(ops)

	// add a label:
	// To position the label, use a TransformOp. This is what
	// layout.Inset does internally too.
	var stack op.StackOp
	stack.Push(ops)
	op.TransformOp{}.Offset(f32.Point{
		X: float32(x0 + e.borderPix),
		Y: float32(y0 + e.borderPix),
	}).Add(ops)

	// this next line clips the label to the box edges, so it does not
	// extend beyond or overflow outside the box.
	e.inPlaceClip(e.borderPix).Add(ops)

	text.Label{Face: e.face, Text: label}.Layout(gtx)
	stack.Pop() // ops)

	// Elias comments:
	//
	//The StackOp is for undoing the effect of the transformation.
	//
	//	If you want rectangular clipping in general, use gioui.org/op/paint.RectClip (as above).
	//	If you want path based clipping, use paint.PathBuilder.
	//	Clip operations are also undone by the StackOps.
	//
	//	Note that text.Label already clips itself to honor the constraints set
	// 	in gtx.Constraints.

}

// The returned ClipOp is at the (0,0) origin, since the
// offset/location is typically already on the stack.
// Thus it is an "in-place" clip, relative to the current
// stack's position, rather than at an arbitrary screen position.
//
func (e *box) inPlaceClip(border int) paint.ClipOp {
	return paint.RectClip(image.Rectangle{
		Min: image.Point{X: 0, Y: 0},
		Max: image.Point{X: e.w - border, Y: e.h - border},
	})
}
