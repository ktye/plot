package plotui

import (
	"fmt"

	"github.com/ktye/plot/clipboard"
	"github.com/lxn/walk"
	"github.com/lxn/win"
)

func (ui *Plot) WriteClipboard() {
	if ui.GetCopyFormat == nil || ui.plots == nil {
		return
	}
	n := len(*ui.plots)
	f := ui.GetCopyFormat(n)

	ui.iplots.SetLimitsTo(ui.plots)

	//win.OpenClipboard used ui.canvas.Handle(), seems to work with 0 too.
	clipboard.CopyToClipboard(*ui.plots, ui.Columns, ui.hi, f)
}

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
