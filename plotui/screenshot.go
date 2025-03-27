package plotui

import (
	"log"
	"syscall"
	"unsafe"

	"github.com/lxn/win"
)

func (ui *Plot) Screenshot() {
	bounds := ui.canvas.ClientBoundsPixels()
	w, h := bounds.Width, bounds.Height
	if ui.plots == nil || w*h == 0 {
		return
	}
	b, err := ui.plots.Emf(w, h, ui.Columns, ui.hi)
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
	win.CloseClipboard()
}

/*
todo: RegisterClipboardFormat(CF_RTF) where CF_RTF is "Rich Text Format"
0000000: 7b5c 7274 6631 5c61 6e73 695c 616e 7369  {\rtf1\ansi\ansi
0000010: 6370 6731 3235 325c 6465 6666 305c 6e6f  cpg1252\deff0\no
0000020: 7569 636f 6d70 6174 5c64 6566 6c61 6e67  uicompat\deflang
0000030: 3130 3331 7b5c 666f 6e74 7462 6c7b 5c66  1031{\fonttbl{\f
0000040: 305c 666e 696c 5c66 6368 6172 7365 7430  0\fnil\fcharset0
0000050: 2043 6f6e 736f 6c61 733b 7d7b 5c66 315c   Consolas;}{\f1\
0000060: 666e 696c 5c66 6368 6172 7365 7430 2043  fnil\fcharset0 C
0000070: 616c 6962 7269 3b7d 7d0d 0a7b 5c63 6f6c  alibri;}}..{\col
0000080: 6f72 7462 6c20 3b5c 7265 6432 3535 5c67  ortbl ;\red255\g
0000090: 7265 656e 305c 626c 7565 303b 5c72 6564  reen0\blue0;\red
00000a0: 3235 355c 6772 6565 6e31 3932 5c62 6c75  255\green192\blu
00000b0: 6530 3b5c 7265 6430 5c67 7265 656e 3137  e0;\red0\green17
00000c0: 365c 626c 7565 3830 3b7d 0d0a 7b5c 2a5c  6\blue80;}..{\*\
00000d0: 6765 6e65 7261 746f 7220 5269 6368 6564  generator Riched
00000e0: 3230 2031 302e 302e 3232 3632 317d 5c76  20 10.0.22621}\v
00000f0: 6965 776b 696e 6434 5c75 6331 200d 0a5c  iewkind4\uc1 ..\
0000100: 7061 7264 5c73 6c32 3430 5c73 6c6d 756c  pard\sl240\slmul
0000110: 7431 5c63 6631 5c66 305c 6673 3234 5c6c  t1\cf1\f0\fs24\l
0000120: 616e 6737 206f 205c 6366 3020 7468 6973  ang7 o \cf0 this
0000130: 2069 7320 6c69 6e65 2031 5c70 6172 0d0a   is line 1\par..
0000140: 5c63 6632 206f 5c63 6631 2020 5c63 6630  \cf2 o\cf1  \cf0
0000150: 2074 6869 7320 6973 206c 696e 6520 325c   this is line 2\
0000160: 7061 720d 0a5c 6366 3320 6f5c 6366 3120  par..\cf3 o\cf1
0000170: 205c 6366 3020 7468 6973 2069 7320 6c69   \cf0 this is li
0000180: 6e65 2033 5c66 315c 6673 3232 5c70 6172  ne 3\f1\fs22\par
0000190: 0d0a 7d0d 0a00                           ..}...
*/

var (
	gdi32              = syscall.MustLoadDLL("gdi32.dll")
	setEnhMetaFileBits = gdi32.MustFindProc("SetEnhMetaFileBits")
	//getEnhMetaFileBits = gdi32.MustFindProc("GetEnhMetaFileBits")

)

/*
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
*/
