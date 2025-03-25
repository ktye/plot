// Package fonts returns ttf data from a string.
package font

import (
	"archive/zip"
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"
)

func MakeFontSizes(ttf []byte, large, small int) (font.Face, font.Face) {
	font, err := truetype.Parse(ttf)
	if err != nil {
		panic(err)
	}
	f1 := truetype.NewFace(font, &truetype.Options{Size: float64(large), DPI: 72})
	f2 := truetype.NewFace(font, &truetype.Options{Size: float64(small), DPI: 72})
	return f1, f2
}

// TTF returns ttf data from a string that is a zipped version of a ttf file.
// See fonts_test for generating strings.
// TTF return nil on error.
func TTF() []byte {
	r := strings.NewReader(Data)
	zr, err := zip.NewReader(r, r.Size())
	if err != nil {
		panic(err)
	}
	for _, f := range zr.File {
		if f.Name == "font.ttf" {
			rc, err := f.Open()
			if err != nil {
				panic(err)
			}
			defer rc.Close()

			b, err := ioutil.ReadAll(rc)
			if err != nil {
				panic(err)
			}
			return b
		}
	}
	return nil
}

func Face(ttf []byte, size int) (font.Face, error) {
	if ttf == nil {
		return nil, fmt.Errorf("ttf is unset")
	}
	font, err := truetype.Parse(ttf)
	if err != nil {
		return nil, err
	}
	return truetype.NewFace(font, &truetype.Options{Size: float64(size), DPI: 72}), nil
}
