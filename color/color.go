// Package color defines plot colors, orders and palettes.
// Orders are used for line plots, palettes are used for paletted images, such as heatmaps.
package color

import (
	"fmt"
	syscolor "image/color"
	"sort"
	"strconv"
)

// Color type which implements the image/color.Color interface.
// The string representation may be any of the predifined color names, or html syntax: #RRGGBB
type Color string

type rgba uint32

var colorMap = map[Color]syscolor.RGBA{
	"transparent": syscolor.RGBA{0, 0, 0, 0},
	"black":       syscolor.RGBA{0, 0, 0, 0xFF},
	"white":       syscolor.RGBA{0xFF, 0xFF, 0xFF, 0xFF},
	"red":         syscolor.RGBA{0xFF, 0, 0, 0xFF},
	"green":       syscolor.RGBA{0, 0xFF, 0, 0xFF},
	"blue":        syscolor.RGBA{0, 0, 0xFF, 0xFF},
	"cyan":        syscolor.RGBA{0, 0xFF, 0xFF, 0xFF},
	"magenta":     syscolor.RGBA{0xFF, 0, 0xFF, 0xFF},
	"yellow":      syscolor.RGBA{0xFF, 0xFF, 0, 0xFF},
	"grey":        syscolor.RGBA{0x80, 0x80, 0x80, 0xFF},
}

// DefaultColorNames returns an unsorted list of known colors.
func DefaultColorNames() []string {
	var names []string
	for k, _ := range colorMap {
		names = append(names, string(k))
	}
	sort.Strings(names)
	return names
}

// Verify if color is known.
func (c Color) Verify() error {
	if _, err := parseHtmlColor(string(c)); err == nil {
		return nil
	}
	if _, ok := colorMap[c]; ok {
		return nil
	}
	return fmt.Errorf("unknown color: %s", c)
}

// Color returns an image/color.
func (c Color) Color() syscolor.Color {
	if rgba, ok := colorMap[c]; ok {
		return rgba
	}
	if rgba, err := parseHtmlColor(string(c)); err != nil {
		// panic("unknown color: " + string(c))
		return syscolor.Black
	} else {
		return rgba
	}
}

// parseHtmlColor parses a HTML color string in the form: #RRGGBB.
func parseHtmlColor(s string) (syscolor.RGBA, error) {
	var c syscolor.RGBA
	if len(s) != 7 || s[0] != '#' {
		return c, fmt.Errorf("cannot parse html color: '%s'", s)
	}

	var err error
	var red, green, blue uint64
	ok := true
	red, err = strconv.ParseUint(s[1:3], 16, 8)
	if err != nil {
		ok = false
	}
	green, err = strconv.ParseUint(s[3:5], 16, 8)
	if err != nil {
		ok = false
	}
	blue, err = strconv.ParseUint(s[5:7], 16, 8)
	if err != nil {
		ok = false
	}
	if ok == false {
		return c, fmt.Errorf("cannot parse html color: '%s'", s)
	}
	c.R = uint8(red)
	c.G = uint8(green)
	c.B = uint8(blue)
	c.A = 0xFF
	return c, nil
}
