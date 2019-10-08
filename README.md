![Screenshot](/screenshot.png)

# hello world for Gio graphics

`hello_gio.go` is my hello world program for the
Gio graphics package for Golang by Elias Naur.

https://gioui.org

Eschewing the elegant constraint layout system for
more direct control, `hello_gio.go` demonstrates how to
plot rectangular boxes at specific screen positions of your choosing.

It then demonstrates how to place labels over those boxes, and to clip
the text to stay inside the box.

Finally, we add the display of a pre-rendered png image on the window,
and color the background yellow. This last part is in `showimg.go`.
Technically this is rendered first, but it was added subsequently.

# intro to Gio

For background on Gio, see Elias's talk:
"GopherCon 2019: Elias Naur - Portable, Immediate Mode GUI Programs for Mobile and Desktop in 100% Go"

https://www.youtube.com/watch?v=9D6eWP4peYM

https://go-talks.appspot.com/github.com/eliasnaur/gophercon-2019-talk/gophercon-2019.slide#1

# installation notes

A link to the installation notes for Gio on various platforms. I did not come across
this until Elias mentioned it to me specifically. It talks about installing the DLLs
and how to link on windows to avoid the extra console.
(hint: `$ go build -ldflags="-H windowsgui" gioui.org/apps/hello`)

https://man.sr.ht/~eliasnaur/gio/install.md

# credits

Author: Jason E. Aten, Ph.D.

License: Unlicense or MIT, same as Gio.

