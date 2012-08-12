package doscolor

import (
	"errors"
	"github.com/anschelsc/w32"
	"os"
	"syscall"
)

type Color uint16

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

func BG(c Color) Color { return c << 4 }

// Some useful masks
const (
	Foreground = 0x0F
	Background = 0xF0
)

type Wrapper struct {
	*os.File
	h     w32.HANDLE
	saved *w32.CONSOLE_SCREEN_BUFFER_INFO
}

func NewWrapper(f *os.File) *Wrapper {
	return &Wrapper{f, w32.HANDLE(f.Fd()), nil}
}

func (w *Wrapper) Save() error {
	w.saved = w32.GetConsoleScreenBufferInfo(w.h)
	if w.saved == nil {
		return syscall.Errno(w32.GetLastError())
	}
	return nil
}

func (w *Wrapper) Restore() error {
	if w.saved == nil {
		return errors.New("attempted to restore without saving.")
	}
	if !w32.SetConsoleTextAttribute(w.h, w.saved.WAttributes) {
		return syscall.Errno(w32.GetLastError())
	}
	return nil
}

// Set c with the default mask (change both foreground and background)
func (w *Wrapper) Set(c Color) error {
	return w.SetMask(c, 0xFF)
}

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
