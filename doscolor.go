package doscolor

import (
	"errors"
	"github.com/anschelsc/w32"
	"os"
	"syscall"
)

// A color is a color. The names are the same as the ANSI colors, but the
// actual hues are pretty different. Wikipedia has a nice chart on the "ANSI
// Colors" page.
type Color uint16

// These are the foreground colors; use BG() to get the equivalent background
// colors.
const (
	Blue Color = 1 << iota
	Green
	Red
	Bright

	Cyan    = Blue | Green
	Magenta = Blue | Red
	Yellow  = Green | Red
	White   = Blue | Green | Red
	Black   = 0
)

// BG turns a foreground color into a background color.
func BG(c Color) Color { return c << 4 }

// Some useful masks for use with SetMask.
const (
	Foreground = 0x0F
	Background = 0xF0
)

// A wrapper lets you change the colors of text written to a DOS console.
type Wrapper struct {
	*os.File
	h     w32.HANDLE
	saved *w32.CONSOLE_SCREEN_BUFFER_INFO
}

// Don't call this unless f is a console (i.e. os.Stdout).
func NewWrapper(f *os.File) *Wrapper {
	return &Wrapper{f, w32.HANDLE(f.Fd()), nil}
}

// Save saves the current coloring so that it can be Restored later.
func (w *Wrapper) Save() error {
	w.saved = w32.GetConsoleScreenBufferInfo(w.h)
	if w.saved == nil {
		return syscall.Errno(w32.GetLastError())
	}
	return nil
}

// Restore restores a previously Saved coloring. It is an error to call Restore
// without ever having called Save.
func (w *Wrapper) Restore() error {
	if w.saved == nil {
		return errors.New("attempted to restore without saving.")
	}
	if !w32.SetConsoleTextAttribute(w.h, w.saved.WAttributes) {
		return syscall.Errno(w32.GetLastError())
	}
	return nil
}

// Change the coloring to c.
func (w *Wrapper) Set(c Color) error {
	return w.SetMask(c, 0xFF)
}

// Change the coloring to c, but only in those bits where m is 1.
func (w *Wrapper) SetMask(c Color, m uint16) error {
	current := w32.GetConsoleScreenBufferInfo(w.h)
	if current == nil {
		return syscall.Errno(w32.GetLastError())
	}
	if !w32.SetConsoleTextAttribute(w.h, apply(current.WAttributes, uint16(c), uint16(m))) {
		return syscall.Errno(w32.GetLastError())
	}
	return nil
}

func apply(old, rep, mask uint16) uint16 {
	return (mask & rep) | (^mask & rep)
}
