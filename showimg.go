package main

// A simple Gio program. See https://gioui.org for more information.

import (
	"bufio"
	"image"
	"image/color"

	"gioui.org/app"
	"gioui.org/f32"
	"gioui.org/gesture"
	"gioui.org/io/pointer"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/paint"
	"gioui.org/text"
	"gioui.org/text/shape"
	"gioui.org/unit"
	"gioui.org/widget"
	"golang.org/x/exp/shiny/iconvg"
	"golang.org/x/image/draw"

	"golang.org/x/exp/shiny/materialdesign/icons"
	//"golang.org/x/image/font/gofont/goregular"
	//"golang.org/x/image/font/gofont/gomonobold"
	"golang.org/x/image/font/gofont/gomono"
	"golang.org/x/image/font/sfnt"

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

type myDrawState struct {
	w       *app.Window
	regular *sfnt.Font
	faces   shape.Faces
	maroon  color.RGBA
	cream   color.RGBA
	black   color.RGBA
	face    *shape.Face
	message string
	gtx     *layout.Context

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
		m.regular, err = sfnt.Parse(gomono.TTF)
		if err != nil {
			panic("failed to load font")
		}
		m.maroon = color.RGBA{127, 0, 0, 255}
		m.cream = color.RGBA{240, 240, 127, 255}
		m.black = color.RGBA{0, 0, 0, 255}
		m.face = m.faces.For(m.regular, unit.Sp(72))
		m.message = "231412" // Hello, Gio
		m.gtx = &layout.Context{
			Queue: w.Queue(),
		}
		//gtx := m.gtx

		m.pngPlot, _, err = LoadImage("points.png")
		panicOn(err)
		m.pngPlotRect = m.pngPlot.(*image.NRGBA).Rect
		vv("m.pngPlot.Rect = '%#v'", m.pngPlotRect)

		var mainErr error

	mainLoop:
		for {
			e := <-w.Events()
			switch e := e.(type) {
			case app.DestroyEvent:
				mainErr = e.Err
				break mainLoop
			case app.UpdateEvent:
				vv("e is '%#v'", e)
				showImage(e, m)
			}
		}
		panicOn(mainErr)
	}()
	app.Main()
}

func showImage(e app.UpdateEvent, m *myDrawState) {
	gtx := m.gtx
	ops := &op.Ops{}
	m.gtx.Ops = ops
	w := m.w
	gtx.Reset(&e.Config, e.Size)
	m.faces.Reset(gtx.Config)

	width := 1000
	height := 600
	x0 := 300
	x1 := x0 + width
	y0 := 200
	y1 := y0 + height

	// background
	ops.Reset()
	paint.ColorOp{Color: m.cream}.Add(ops)
	paint.PaintOp{Rect: f32.Rectangle{Max: f32.Point{X: float32(e.Size.X), Y: float32(e.Size.Y)}}}.Add(ops)

	// foreground
	paint.ColorOp{Color: m.maroon}.Add(ops) // color
	//margin := int(float64(minInt(e.Size.X, e.Size.Y)) * 0.10)
	box := image.Rectangle{
		Min: image.Point{X: (x0), Y: (y0)},
		Max: image.Point{X: (x1), Y: (y1)},
	}
	boxPos := toRectF(box)
	_ = boxPos
	// image/texture foreground

	ico := (&icon{src: icons.NavigationArrowBack, size: unit.Dp(24)}).image(gtx, rgb(0xffffff))
	vv("ico.Bounds = '%#v'", ico.Bounds()) // 41x41
	// widget.Image is a widget that displays an image.
	//wi := widget.Image{Src: ico, Rect: ico.Bounds(), Scale: 1}
	//wi.Layout(gtx)
	//gtx.Dimensions.Size.X += gtx.Px(unit.Dp(4))
	//pointer.RectAreaOp{Rect: image.Rectangle{Max: gtx.Dimensions.Size}}.Add(gtx.Ops)
	//pointer.RectAreaOp{Rect: box}.Add(gtx.Ops)
	//pointer.RectAreaOp{Rect: ico.Bounds()}.Add(gtx.Ops)

	imgSrc := m.pngPlot
	//imgSrc := ico
	imgPos := image.Rectangle{
		Min: image.Point{X: x0, Y: y0},
		Max: image.Point{X: x1, Y: y1},
	}
	paint.ImageOp{Src: imgSrc, Rect: imgPos}.Add(ops) // or image
	paintOp := paint.PaintOp{
		Rect: boxPos,
	}
	paintOp.Add(ops)

	// add number label
	//labelOps := &op.Ops{}
	var material op.MacroOp
	material.Record(ops)
	paint.ColorOp{Color: m.black}.Add(ops)
	material.Stop()

	// position
	op.TransformOp{}.Offset(f32.Point{X: float32(x0), Y: float32(y0)}).Add(ops)

	// clip to rectangle
	clipBox := image.Rectangle{
		Min: image.Point{X: 0, Y: 0},
		Max: image.Point{X: x1 - x0, Y: y1 - y0},
	}
	_ = clipBox
	paint.RectClip(clipBox).Add(ops)

	//	vv("m.face = '%#v'", m.face)
	//	vv("m.face = '%#v'", m.face.face)
	lab := text.Label{
		Material:  material,
		Face:      m.face,
		Alignment: text.Start,
		Text:      m.message}
	lab.Layout(gtx)
	//paintOp.Add(ops)
	w.Update(ops)
}

func minInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func LoadImage(filename string) (image.Image, string, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, "", err
	}
	defer f.Close()
	return image.Decode(bufio.NewReader(f))
}

func (ic *icon) image(c unit.Converter, col color.RGBA) image.Image {
	sz := c.Px(ic.size)
	if sz == ic.imgSize {
		return ic.img
	}
	m, _ := iconvg.DecodeMetadata(ic.src)
	dx, dy := m.ViewBox.AspectRatio()
	img := image.NewRGBA(image.Rectangle{Max: image.Point{X: sz, Y: int(float32(sz) * dy / dx)}})
	var ico iconvg.Rasterizer
	ico.SetDstImage(img, img.Bounds(), draw.Src)
	m.Palette[0] = col
	iconvg.Decode(&ico, ic.src, &iconvg.DecodeOptions{
		Palette: &m.Palette,
	})
	ic.img = img
	ic.imgSize = sz
	return img
}

type icon struct {
	src  []byte
	size unit.Value

	// Cached values.
	img     image.Image
	imgSize int
}

type IconButton struct {
	Icon  *icon
	Inset layout.Inset
	buttonState
}

type buttonState struct {
	click  gesture.Click
	clicks int
}

func toRectF(r image.Rectangle) f32.Rectangle {
	return f32.Rectangle{
		Min: f32.Point{X: float32(r.Min.X), Y: float32(r.Min.Y)},
		Max: f32.Point{X: float32(r.Max.X), Y: float32(r.Max.Y)},
	}
}

func rgb(c uint32) color.RGBA {
	return argb((0xff << 24) | c)
}

func argb(c uint32) color.RGBA {
	return color.RGBA{
		A: uint8(c >> 24),
		R: uint8(c >> 16),
		G: uint8(c >> 8),
		B: uint8(c),
	}
}
