package serve

import (
	"bytes"
	"fmt"
	"log"
	"runtime"
	"syscall"
	"unsafe"

	"github.com/ktye/plot"
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
func WriteClipboard(plots plot.Plots, columns int, hi []plot.HighlightID, f CopyFormat) {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	fmt.Println("WriteClipboard", len(plots), "plot", columns, "columns, formats:", f)
	if len(plots) == 0 {
		return
	}
	w, h, font, f1, f2, captionfont, f3 := f.Width, f.Height, f.PlotFont, f.F1, f.F2, f.CaptionFont, f.F3

	//ui.iplots.SetLimitsTo(ui.plots) // may have changed interactively

	if w*h == 0 {
		return
	}
	b, err := plots.Emf(w, h, columns, hi, font, f1, f2)
	if err != nil {
		fmt.Println(err)
		return
	}

	r, _, err := openClipboard.Call()
	fmt.Println("openClipboad", r, err)

	defer closeClipboard.Call()
	emptyClipboard.Call()

	//win.OpenClipboard(ui.canvas.Handle())
	//win.EmptyClipboard()
	fmt.Println("metafile length:", len(b))
	r, _, err = setEnhMetaFileBits.Call(uintptr(len(b)), uintptr(unsafe.Pointer(&b[0])))
	if r == 0 {
		log.Print(err)
	} else {
		r, _, err = setClipboardData.Call(14, r)
		fmt.Println("set clipboard data(metafile):", r, err)
	}

	caption := plots[0].Caption
	if caption != nil {
		clipCaptionRtf(caption, captionfont, f3)
		clipCaptionTxt(caption)
	}
}
func clipCaptionTxt(caption *plot.Caption) {
	var buf bytes.Buffer
	if _, e := caption.WriteTable(&buf, plot.Numbers); e != nil {
		return
	}
	b := buf.Bytes()
	s, _ := syscall.UTF16FromString(string(b))

	/*

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
	*/

	hMem, _, err := gAlloc.Call(2, uintptr(len(s)*int(unsafe.Sizeof(s[0]))))
	if hMem == 0 {
		fmt.Println("failed to alloc global memory:", err)
		return
	}

	p, _, err := gLock.Call(hMem)
	if p == 0 {
		fmt.Println("failed to lock global memory:", err)
		return
	}
	defer gUnlock.Call(hMem)

	memMove.Call(p, uintptr(unsafe.Pointer(&s[0])), uintptr(len(s)*int(unsafe.Sizeof(s[0]))))

	v, _, err := setClipboardData.Call(13, hMem)
	if v == 0 {
		gFree.Call(hMem)
		fmt.Println("failed to set text to clipboard:", err)
		return
	}
}

func clipCaptionRtf(caption *plot.Caption, font string, fs int) { // caption to clipboard as rtf
	return
}

/*
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
*/
func regFormat(s string) uintptr {
	u, _ := syscall.UTF16FromString(s)
	r, _, _ := registerClipboardFormat.Call(uintptr(unsafe.Pointer(&u[0])))
	return r //0 error
}

var (
	gdi32              = syscall.MustLoadDLL("gdi32.dll")
	setEnhMetaFileBits = gdi32.MustFindProc("SetEnhMetaFileBits")
	//getEnhMetaFileBits = gdi32.MustFindProc("GetEnhMetaFileBits")

	kernel32 = syscall.MustLoadDLL("kernel32")
	gAlloc   = kernel32.MustFindProc("GlobalAlloc")
	gFree    = kernel32.MustFindProc("GlobalFree")
	gLock    = kernel32.MustFindProc("GlobalLock")
	gUnlock  = kernel32.MustFindProc("GlobalUnlock")
	memMove  = kernel32.MustFindProc("RtlMoveMemory")

	user32                  = syscall.MustLoadDLL("user32")
	openClipboard           = user32.MustFindProc("OpenClipboard")
	closeClipboard          = user32.MustFindProc("CloseClipboard")
	emptyClipboard          = user32.MustFindProc("EmptyClipboard")
	setClipboardData        = user32.MustFindProc("SetClipboardData")
	registerClipboardFormat = user32.MustFindProc("RegisterClipboardFormatW")
)
