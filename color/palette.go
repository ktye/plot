package color

import (
	"fmt"
	imcolor "image/color"
	"sort"
)

// Palette defines a color palette.
type Palette string

// generator is a function that can produce a palette.
type generator func() [256]imcolor.Color

// All known palette generators are stored in this map.
var palettes = map[Palette]generator{
	"":     inferno, // default Palette if unset.
	"grey": grey,
}

// grey make a grey-scale palette.
func grey() (g [256]imcolor.Color) {
	for i := range g {
		g[i] = imcolor.Gray{Y: uint8(i)}
	}
	return g
}

func (p Palette) Verify() error {
	if len(p) > 0 && p[0] == '-' {
		p = Palette(string(p)[1:])
	}
	if _, ok := palettes[p]; !ok {
		return fmt.Errorf("unknown palette: '%s'", p)
	}
	return nil
}

// Palette returns an array of colors.
// If the palette name starts with "-" it is inverted.
func (p Palette) Palette() [256]imcolor.Color {
	inv := false
	if len(p) > 0 && p[0] == '-' {
		inv = true
		p = p[1:]
	}
	if f, ok := palettes[p]; !ok {
		return grey()
	} else {
		pal := f()
		if inv {
			var invpal [256]imcolor.Color
			for i := 0; i < 256; i++ {
				invpal[i] = pal[255-i]
			}
			return invpal
		} else {
			return pal
		}
	}
}

func GetAllPalettes() []string {
	var all []string
	for key := range palettes {
		all = append(all, string(key), "-"+string(key))
	}
	sort.Strings(all)
	return all
}
