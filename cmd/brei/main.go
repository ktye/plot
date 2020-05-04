package main

/*
export xmin xmax ymin ymax rows cols

# xy-plot
cat << EOF | ./brei
0   -0.1
0.1 0
0.2 0.1
0.3 0.2
0.4 0.1
0.5 0.8
0.6 0.6
0.7 0.3
0.8 -0.1
0.9 -0.2
1   0
EOF

# polar plot (r=ymax, cols=2*rows)
ymax=3; echo -e "2.9 0\n2.5 30\n2a60" | ./brei.exe p

# 2x3 blocks
ymax=3; echo -e "2.9 0\n2.5 30\n2a60" | ./brei.exe 6p
*/

import (
	"fmt"
	"image"
	"image/color"
	"math"
	"os"
	"strconv"
	"strings"
)

func main() {
	// r := []rune("⠀⠁⠂⠃⠄⠅⠆⠇⠈⠉⠊⠋⠌⠍⠎⠏⠐⠑⠒⠓⠔⠕⠖⠗⠘⠙⠚⠛⠜⠝⠞⠟⠠⠡⠢⠣⠤⠥⠦⠧⠨⠩⠪⠫⠬⠭⠮⠯⠰⠱⠲⠳⠴⠵⠶⠷⠸⠹⠺⠻⠼⠽⠾⠿⡀⡁⡂⡃⡄⡅⡆⡇⡈⡉⡊⡋⡌⡍⡎⡏⡐⡑⡒⡓⡔⡕⡖⡗⡘⡙⡚⡛⡜⡝⡞⡟⡠⡡⡢⡣⡤⡥⡦⡧⡨⡩⡪⡫⡬⡭⡮⡯⡰⡱⡲⡳⡴⡵⡶⡷⡸⡹⡺⡻⡼⡽⡾⡿⢀⢁⢂⢃⢄⢅⢆⢇⢈⢉⢊⢋⢌⢍⢎⢏⢐⢑⢒⢓⢔⢕⢖⢗⢘⢙⢚⢛⢜⢝⢞⢟⢠⢡⢢⢣⢤⢥⢦⢧⢨⢩⢪⢫⢬⢭⢮⢯⢰⢱⢲⢳⢴⢵⢶⢷⢸⢹⢺⢻⢼⢽⢾⢿⣀⣁⣂⣃⣄⣅⣆⣇⣈⣉⣊⣋⣌⣍⣎⣏⣐⣑⣒⣓⣔⣕⣖⣗⣘⣙⣚⣛⣜⣝⣞⣟⣠⣡⣢⣣⣤⣥⣦⣧⣨⣩⣪⣫⣬⣭⣮⣯⣰⣱⣲⣳⣴⣵⣶⣷⣸⣹⣺⣻⣼⣽⣾⣿")

	four, two := 4, 2 // 2x8 braille default or 2x3 legacy computing blocks
	polar := false
	if len(os.Args) > 1 {
		s := os.Args[1]
		if len(s) > 0 && s[0] == '6' {
			four = 3
			s = s[1:]
		} else if len(s) > 0 && s[0] == '1' {
			four = 1
			two = 1
			s = s[1:]
		} else if len(s) > 0 && s[0] == '8' {
			s = s[1:]
		}
		if s == "p" {
			polar = true
		}
	}

	flt := func(s string, def float64) float64 {
		f, err := strconv.ParseFloat(s, 64)
		if err != nil {
			return def
		}
		return f
	}
	xmin, xmax := flt(os.Getenv("xmin"), 0), flt(os.Getenv("xmax"), 1)
	ymin, ymax := flt(os.Getenv("ymin"), -1), flt(os.Getenv("ymax"), 1)
	cols, rows := int(flt(os.Getenv("cols"), 72)), int(flt(os.Getenv("rows"), 10))
	scale := func(x, a, b, c, d float64) float64 { return c + (x-a)*(d-c)/(b-a) }

	w := cols * two
	h := rows * four
	if polar {
		cols = 2 * rows
		w = h
		xmin, ymin, xmax = -ymax, -ymax, ymax
	}
	m := image.NewGray(image.Rectangle{Max: image.Point{w, h}})
	if polar {
		circle(w/2-1, m)
	}

	sx := func(x float64) int { return int(scale(x, xmin, xmax, 0, float64(w)-1)) }
	sy := func(y float64) int { return int(scale(y, ymin, ymax, float64(h)-1, 0)) }

	var x, y, xx, yy int
	var X, Y float64
	var dummy rune // allows "1 30" or "1a30"
	for i := 0; ; i++ {
		if n, err := fmt.Scanf("%f%c%f\n", &X, &dummy, &Y); n != 3 || err != nil {
			break
		}
		if polar {
			Y, X = X*math.Cos(Y*math.Pi/180.0), X*math.Sin(Y*math.Pi/180.0) // clockwise
		}
		x, y = sx(X), sy(Y)
		if polar {
			m.Set(x, y, color.Gray{255})
		} else {
			if i > 0 {
				line(xx, yy, x, y, m)
			}
			xx, yy = x, y
		}
	}

	var out flusher
	switch four {
	case 1:
		out = &asciiFlusher{}
	case 3:
		out = &legacyFlusher{}
	case 4:
		out = &brailleFlusher{}
	default:
		panic("?")
	}
	if polar {
		out.label(fmt.Sprintf("%v", ymax))
	} else {
		fmt.Printf("[%v, %v]\n", ymin, ymax)
	}
	out.flush(m)

	if polar {
	} else {
		s := fmt.Sprintf("[%v, %v]", xmin, xmax)
		fmt.Printf("%s%s\n", strings.Repeat(" ", cols-len(s)), s)
	}
}

/*
func line(x0, y0, x1, y1 int, m *image.Gray) {
	var dx, dy, sx, sy, e, e2 int
	dx = abs(x1 - x0)
	dy = -abs(y1 - y0)
	if sx = -1; x0 < x1 {
		sx = 1
	}
	if sy = -1; y0 < y1 {
		sy = 1
	}
	e = dx + dy
	for {
		m.Set(x0, y0, color.Gray{255})
		if x0 == x1 && y0 == y1 {
			break
		}
		if e2 = 2 * e; e2 >= dy {
			e += dy
			x0 += sx
		} else if e2 <= dx {
			e += dx
			y0 += sy
		}
	}
}
*/
func abs(x int) int {
	if x < 0 {
		return -x
	} else {
		return x
	}
}
func line(x0, y0, x1, y1 int, m *image.Gray) {
	if abs(y1-y0) < abs(x1-x0) {
		if x0 > x1 {
			line1(x1, y1, x0, y0, m)
		} else {
			line1(x0, y0, x1, y1, m)
		}
	} else {
		if y0 > y1 {
			line2(x1, y1, x0, y0, m)
		} else {
			line2(x0, y0, x1, y1, m)
		}
	}
}
func line1(x0, y0, x1, y1 int, m *image.Gray) {
	dx, dy, yi := x1-x0, y1-y0, 1
	if dy < 0 {
		yi, dy = -1, -dy
	}
	d, y := 2*dy-dx, y0
	for x := x0; x <= x1; x++ {
		m.Set(x, y, color.Gray{255})
		if d > 0 {
			y += yi
			d -= 2 * dx
		}
		d += 2 * dy
	}
}
func line2(x0, y0, x1, y1 int, m *image.Gray) {
	dx, dy, xi := x1-x0, y1-y0, 1
	if dx < 0 {
		xi, dx = -1, -dx
	}
	d, x := 2*dx-dy, x0
	for y := y0; y <= y1; y++ {
		m.Set(x, y, color.Gray{255})
		if d > 0 {
			x += xi
			d -= 2 * dy
		}
		d += 2 * dx
	}
}
func circle(r int, m *image.Gray) {
	var x, y, e, c int
	x = -r
	e = 2 - 2*r
	c = r
	for x < 0 {
		m.Set(c-x, c+y, color.Gray{255})
		m.Set(c-y, c-x, color.Gray{255})
		m.Set(c+x, c-y, color.Gray{255})
		m.Set(c+y, c+x, color.Gray{255})
		r = e
		if r <= y {
			y++
			e += 2*y + 1
		}
		if r > x || e > y {
			x++
			e += 2*x + 1
		}
	}
}

type flusher interface {
	label(string)
	flush(*image.Gray)
}

type Braille [2][4]int // modified from: github.com/kevin-cantwell/dotmatrix
type brailleFlusher struct{ s string }

func (b Braille) Rune() rune {
	lowEndian := [8]int{b[0][0], b[0][1], b[0][2], b[1][0], b[1][1], b[1][2], b[0][3], b[1][3]}
	var v int
	for i, x := range lowEndian {
		v += int(x) << uint(i)
	}
	return rune(v) + '\u2800'
}
func (b Braille) String() string {
	return string(b.Rune())
}

func (b *brailleFlusher) label(s string) { b.s = s }
func (b *brailleFlusher) flush(m *image.Gray) {
	w := os.Stdout
	max := m.Bounds().Max
	for py := 0; py < max.Y; py += 4 {
		for px := 0; px < max.X; px += 2 {
			var b Braille
			for y := 0; y < 4; y++ {
				for x := 0; x < 2; x++ {
					if px+x >= max.X || py+y >= max.Y {
						continue
					}
					if m.GrayAt(px+x, py+y).Y == 255 {
						b[x][y] = 1
					}
				}
			}
			w.Write([]byte(b.String()))
		}
		if py+4 >= max.Y {
			w.Write([]byte(b.s))
		}
		w.Write([]byte{'\n'})
	}
}

type legacyFlusher struct{ s string }

func (b *legacyFlusher) label(s string) { b.s = s }
func (b *legacyFlusher) flush(m *image.Gray) {
	t := make([]rune, 64) // 2x3 block character from legacy computing range 0x1fb00..
	// 4 are missing in the legacy range: blank, full and the left/right halfs
	t[0] = rune(32)      // blank
	t[21] = rune(0x258c) // left half block 0b010101
	t[42] = rune(0x2590) // right half 0b101010
	t[63] = rune(0x2588) // full block
	i := 1
	for r := rune(0x1fb00); r < 0x1fb00+60; r++ {
		if i == 21 || i == 42 {
			i++
		}
		t[i] = r
		i++
	}

	max := m.Bounds().Max
	w, h := max.X, max.Y
	cols := w / 2
	s := make([]rune, cols)
	for k := 0; k < h; k += 3 {
		l := make([]byte, cols)
		for n := 0; n < 3; n++ {
			for i := 0; i < w; i += 2 {
				if m.GrayAt(i, k+n).Y == 255 {
					l[(i / 2)] += 1 << uint(2*n)
				}
				if m.GrayAt(i+1, k+n).Y == 255 {
					l[(i / 2)] += 1 << uint(2*n+1)
				}
			}
		}
		for i := 0; i < cols; i++ {
			s[i] = t[l[i]]
		}
		c := ""
		if k+3 >= h {
			c = b.s
		}
		fmt.Printf("%s%s\n", string(s), c)
	}
}

type asciiFlusher struct{ s string }

func (b *asciiFlusher) label(s string) { b.s = s }
func (b *asciiFlusher) flush(m *image.Gray) {
	max := m.Bounds().Max
	w, h := max.X, max.Y
	for k := 0; k < h; k++ {
		for i := 0; i < w; i++ {
			if m.GrayAt(i, k).Y == 255 {
				fmt.Printf("X")
			} else {
				fmt.Printf(" ")
			}
		}
		if k == h-1 {
			fmt.Printf("%s", b.s)
		}
		fmt.Println()
	}
}
