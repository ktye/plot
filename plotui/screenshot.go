package plotui

import (
	"bytes"
	"fmt"
	"log"
	"strconv"
	"syscall"
	"unsafe"

	"github.com/ktye/plot"
	"github.com/lxn/walk"
	"github.com/lxn/win"
)

type CopyFormat struct {
	Width       int    // image width (px)
	Height      int    // image height (px)
	PlotFont    string // plot font family
	F1          int    // plot font title size
	F2          int    // plot font number label size
	CaptionFont string // caption font family
	F3          int    // caption font size
}

// WriteClipboard writes the plot image as EMFPLUS to clipboard as CF_ENHMETAFILE,
// the caption table as "Rich Text Format" with line color markers and as plain text.
// func (ui *Plot) WriteClipboard(w, h int, font string, f1, f2 int, captionfont string, f3 int) {
func (ui *Plot) WriteClipboard() {
	if ui.GetCopyFormat == nil || ui.plots == nil {
		return
	}
	n := len(*ui.plots)
	f := ui.GetCopyFormat(n)
	w, h, font, f1, f2, captionfont, f3 := f.Width, f.Height, f.PlotFont, f.F1, f.F2, f.CaptionFont, f.F3

	ui.iplots.SetLimitsTo(ui.plots) // may have changed interactively

	if w*h == 0 {
		return
	}
	b, err := ui.plots.Emf(w, h, ui.Columns, ui.hi, font, f1, f2)
	if err != nil {
		log.Print(err)
		return
	}

	win.OpenClipboard(ui.canvas.Handle())
	win.EmptyClipboard()
	r, _, err := setEnhMetaFileBits.Call(uintptr(len(b)), uintptr(unsafe.Pointer(&b[0])))
	if r == 0 {
		log.Print(err)
	} else {
		win.SetClipboardData(14, win.HANDLE(r))
	}
	ui.clipCaptionRtf(captionfont, f3)
	ui.clipCaptionTxt()
	win.CloseClipboard()
}
func (ui *Plot) clipCaptionTxt() {
	if ui.caption == nil {
		return
	}
	var buf bytes.Buffer
	if _, e := ui.caption.WriteTable(&buf, plot.Numbers); e != nil {
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
func (ui *Plot) clipCaptionRtf(font string, fs int) { // caption to clipboard as rtf
	if ui.caption == nil {
		return
	}
	var buf bytes.Buffer
	if _, e := ui.caption.WriteTable(&buf, plot.Rtf); e != nil {
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
	//getEnhMetaFileBits = gdi32.MustFindProc("GetEnhMetaFileBits")

	user32                  = syscall.MustLoadDLL("user32")
	registerClipboardFormat = user32.MustFindProc("RegisterClipboardFormatW")
)

// Screenshot copies the plot image to the clipboard.
func (ui *Plot) Screenshot() {
	ui.canvas.SetPaintMode(walk.PaintNormal)
	defer func() {
		ui.canvas.SetPaintMode(walk.PaintBuffered)
	}()
	if hbmp, err := hBitmapFromWindow(ui.canvas); err != nil {
		return
	} else {
		// TODO: do we need the mainwindow?
		hBitmapToClipboard(ui.canvas, hbmp)
	}
}

func hBitmapFromWindow(window walk.Window) (win.HBITMAP, error) {
	hdcMem := win.CreateCompatibleDC(0)
	if hdcMem == 0 {
		return 0, fmt.Errorf("CreateCompatibleDC failed")
	}
	defer win.DeleteDC(hdcMem)

	var r win.RECT
	if !win.GetWindowRect(window.Handle(), &r) {
		return 0, fmt.Errorf("GetWindowRect failed")
	}

	hdc := win.GetDC(window.Handle())
	width, height := r.Right-r.Left, r.Bottom-r.Top
	hBmp := win.CreateCompatibleBitmap(hdc, width, height)
	win.ReleaseDC(window.Handle(), hdc)

	hOld := win.SelectObject(hdcMem, win.HGDIOBJ(hBmp))
	flags := win.PRF_CHILDREN | win.PRF_CLIENT | win.PRF_ERASEBKGND | win.PRF_NONCLIENT | win.PRF_OWNED
	window.SendMessage(win.WM_PRINT, uintptr(hdcMem), uintptr(flags))

	win.SelectObject(hdcMem, hOld)
	return hBmp, nil
}

func hBitmapToClipboard(window walk.Window, hBmp win.HBITMAP) {
	win.OpenClipboard(window.Handle())
	win.EmptyClipboard()
	win.SetClipboardData(win.CF_BITMAP, win.HANDLE(hBmp))
	win.CloseClipboard()
}
