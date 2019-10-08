package main

import (
	"fmt"
	"io"
	"os"
	"path"
	"runtime"
	"sync"
	"time"

	"4d63.com/tz"
)

var NYC *time.Location

func init() {
	var err error
	NYC, err = tz.LoadLocation("America/New_York")
	panicOn(err)

}

// for tons of debug output
var Verbose bool = false
var VerboseVerbose bool = false

func P(format string, a ...interface{}) {
	if Verbose {
		TSPrintf(format, a...)
	}
}

func PP(format string, a ...interface{}) {
	if VerboseVerbose {
		TSPrintf(format, a...)
	}
}

func VV(format string, a ...interface{}) {
	TSPrintf(format, a...)
}

func AlwaysPrintf(format string, a ...interface{}) {
	TSPrintf(format, a...)
}

var vv = VV

// without the file/line, otherwise the same as PP
func PPP(format string, a ...interface{}) {
	if VerboseVerbose {
		Printf("\n%s ", ts())
		Printf(format+"\n", a...)
	}
}

func PB(w io.Writer, format string, a ...interface{}) {
	if Verbose {
		fmt.Fprintf(w, "\n"+format+"\n", a...)
	}
}

var tsPrintfMut sync.Mutex

// time-stamped printf
func TSPrintf(format string, a ...interface{}) {
	tsPrintfMut.Lock()
	Printf("\n%s %s ", FileLine(3), ts())
	Printf(format+"\n", a...)
	tsPrintfMut.Unlock()
}

// get timestamp for logging purposes
func ts() string {
	return time.Now().In(NYC).Format("2006-01-02 15:04:05.999 -0700 MST")
}

// so we can multi write easily, use our own printf
var OurStdout io.Writer = os.Stdout

// Printf formats according to a format specifier and writes to standard output.
// It returns the number of bytes written and any write error encountered.
func Printf(format string, a ...interface{}) (n int, err error) {
	return fmt.Fprintf(OurStdout, format, a...)
}

func FileLine(depth int) string {
	_, fileName, fileLine, ok := runtime.Caller(depth)
	var s string
	if ok {
		s = fmt.Sprintf("%s:%d", path.Base(fileName), fileLine)
	} else {
		s = ""
	}
	return s
}

func p(format string, a ...interface{}) {
	if Verbose {
		TSPrintf(format, a...)
	}
}

var pp = PP

func pbb(w io.Writer, format string, a ...interface{}) {
	if Verbose {
		fmt.Fprintf(w, "\n"+format+"\n", a...)
	}
}

// quieted for now, uncomment below to display
func VPrintf(format string, a ...interface{}) (n int, err error) {
	//return fmt.Fprintf(OurStdout, format, a...)
	return
}

func QPrintf(format string, a ...interface{}) (n int, err error) {
	//return fmt.Fprintf(OurStdout, format, a...)
	return
}

func panicOn(err error) {
	if err != nil {
		panic(err)
	}
}

func stopOn(err error) {
	if err == nil {
		return
	}
	fmt.Fprintf(os.Stderr, "%s: %v\n", FileLine(2), err.Error())
	os.Exit(1)
}

// abort the program with error code 1 after printing msg to Stderr.
func stop(msg interface{}) {
	switch e := msg.(type) {
	case error:
		fmt.Fprintf(os.Stderr, "%s: %s\n", FileLine(2), e.Error())
		os.Exit(1)
	default:
		fmt.Fprintf(os.Stderr, "%s: %v\n", FileLine(2), msg)
		os.Exit(1)
	}
}
