package clipboard

import (
	"bytes"
	"log"
	"strconv"
	"syscall"
	"unsafe"

	"github.com/ktye/plot"
	"github.com/lxn/win"
)

// CopyToClipboard writes the plot image as EMFPLUS to clipboard as CF_ENHMETAFILE,
// the caption table as "Rich Text Format" with line color markers and as plain text.
func CopyToClipboard(plots plot.Plots, columns int, hi []plot.HighlightID, f CopyFormat) {
	w, h, font, f1, f2, captionfont, f3 := f.Width, f.Height, f.PlotFont, f.F1, f.F2, f.CaptionFont, f.F3
	//todo: ui.iplots.SetLimitsTo(ui.plots) // may have changed interactively

	if len(plots) == 0 || w*h == 0 {
		return
	}
	b, err := plots.Emf(w, h, columns, hi, font, f1, f2)
	if err != nil {
		log.Print(err)
		return
	}

	win.OpenClipboard(0) //ui.canvas.Handle())
	win.EmptyClipboard()
	r, _, err := setEnhMetaFileBits.Call(uintptr(len(b)), uintptr(unsafe.Pointer(&b[0])))
	if r == 0 {
		log.Print(err)
	} else {
		win.SetClipboardData(14, win.HANDLE(r))
	}
	clipCaptionRtf(plots[0].Caption, captionfont, f3)
	clipCaptionTxt(plots[0].Caption)
	win.CloseClipboard()

}
func clipCaptionTxt(caption *plot.Caption) {
	if caption == nil {
		return
	}
	var buf bytes.Buffer
	if _, e := caption.WriteTable(&buf, plot.Numbers); e != nil {
		return
	}
	b := buf.Bytes()
	u, _ := syscall.UTF16FromString(string(b))
	h := win.GlobalAlloc(win.GMEM_MOVEABLE, uintptr(len(u)*2))
	if h == 0 {
		return
	}
	p := win.GlobalLock(h)
	if p == nil {
		return
	}
	win.MoveMemory(p, unsafe.Pointer(&u[0]), uintptr(len(u)*2))
	win.GlobalUnlock(h)
	if 0 == win.SetClipboardData(win.CF_UNICODETEXT, win.HANDLE(h)) {
		win.GlobalFree(h)
	}
}
func clipCaptionRtf(caption *plot.Caption, font string, fs int) { // caption to clipboard as rtf
	if caption == nil {
		return
	}
	var buf bytes.Buffer
	if _, e := caption.WriteTable(&buf, plot.Rtf); e != nil {
		return
	}
	b := buf.Bytes()
	if font == "" {
		font = "Consolas"
	}
	if fs == 0 {
		fs = 12
	}
	if font != "Consolas" {
		b = bytes.Replace(b, []byte("Consolas"), []byte(font), 1)
	}
	if fs != 12 {
		b = bytes.Replace(b, []byte("\\fs24"), []byte(`\fs`+strconv.Itoa(2*fs)), 1)
	}

	rtf := regFormat("Rich Text Format")
	if rtf == 0 {
		return
	}

	h := win.GlobalAlloc(win.GMEM_MOVEABLE, uintptr(len(b)))
	if h == 0 {
		return //GlobalAlloc failed
	}
	p := win.GlobalLock(h)
	if p == nil {
		return
	}
	win.MoveMemory(p, unsafe.Pointer(&b[0]), uintptr(len(b)))
	win.GlobalUnlock(h)
	if 0 == win.SetClipboardData(uint32(rtf), win.HANDLE(h)) {
		win.GlobalFree(h)
	}
}
func regFormat(s string) uintptr {
	u, _ := syscall.UTF16FromString(s)
	r, _, _ := registerClipboardFormat.Call(uintptr(unsafe.Pointer(&u[0])))
	return r //0 error
}

var (
	gdi32              = syscall.MustLoadDLL("gdi32.dll")
	setEnhMetaFileBits = gdi32.MustFindProc("SetEnhMetaFileBits")

	user32                  = syscall.MustLoadDLL("user32")
	registerClipboardFormat = user32.MustFindProc("RegisterClipboardFormatW")
)
