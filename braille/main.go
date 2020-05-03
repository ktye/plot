package main

/*
export xmin xmax ymin ymax rows cols
cat EOF | ./braille
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
*/

import (
	"fmt"
	"image"
	"image/color"
	"io"
	"os"
	"strconv"
	"strings"
)

func main() {
	// r := []rune("⠀⠁⠂⠃⠄⠅⠆⠇⠈⠉⠊⠋⠌⠍⠎⠏⠐⠑⠒⠓⠔⠕⠖⠗⠘⠙⠚⠛⠜⠝⠞⠟⠠⠡⠢⠣⠤⠥⠦⠧⠨⠩⠪⠫⠬⠭⠮⠯⠰⠱⠲⠳⠴⠵⠶⠷⠸⠹⠺⠻⠼⠽⠾⠿⡀⡁⡂⡃⡄⡅⡆⡇⡈⡉⡊⡋⡌⡍⡎⡏⡐⡑⡒⡓⡔⡕⡖⡗⡘⡙⡚⡛⡜⡝⡞⡟⡠⡡⡢⡣⡤⡥⡦⡧⡨⡩⡪⡫⡬⡭⡮⡯⡰⡱⡲⡳⡴⡵⡶⡷⡸⡹⡺⡻⡼⡽⡾⡿⢀⢁⢂⢃⢄⢅⢆⢇⢈⢉⢊⢋⢌⢍⢎⢏⢐⢑⢒⢓⢔⢕⢖⢗⢘⢙⢚⢛⢜⢝⢞⢟⢠⢡⢢⢣⢤⢥⢦⢧⢨⢩⢪⢫⢬⢭⢮⢯⢰⢱⢲⢳⢴⢵⢶⢷⢸⢹⢺⢻⢼⢽⢾⢿⣀⣁⣂⣃⣄⣅⣆⣇⣈⣉⣊⣋⣌⣍⣎⣏⣐⣑⣒⣓⣔⣕⣖⣗⣘⣙⣚⣛⣜⣝⣞⣟⣠⣡⣢⣣⣤⣥⣦⣧⣨⣩⣪⣫⣬⣭⣮⣯⣰⣱⣲⣳⣴⣵⣶⣷⣸⣹⣺⣻⣼⣽⣾⣿")

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

	w := cols * 2
	h := rows * 4
	m := image.NewGray(image.Rectangle{Max: image.Point{w, h}})

	sx := func(x float64) int { return int(scale(x, xmin, xmax, 0, float64(w)-1)) }
	sy := func(y float64) int { return int(scale(y, ymin, ymax, float64(h)-1, 0)) }

	var x, y, xx, yy int
	var X, Y float64
	for i := 0; ; i++ {
		if n, err := fmt.Scanf("%f %f", &X, &Y); n != 2 || err != nil {
			break
		}
		x, y = sx(X), sy(Y)
		if i > 0 {
			line(xx, yy, x, y, m)
		}
		xx, yy = x, y
	}

	fmt.Printf("[%v, %v]\n", ymin, ymax)
	br := BrailleFlusher{}
	br.Flush(os.Stdout, m)

	s := fmt.Sprintf("[%v, %v]", xmin, xmax)
	fmt.Printf("%s%s\n", strings.Repeat(" ", cols-len(s)), s)
}

func abs(x int) int {
	if x < 0 {
		return -x
	} else {
		return x
	}
}
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

// from: github.com/kevin-cantwell/dotmatrix
type Braille [2][4]int
type BrailleFlusher struct{}

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
func (BrailleFlusher) Flush(w io.Writer, m *image.Gray) error {
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
			if _, err := w.Write([]byte(b.String())); err != nil {
				return err
			}
		}
		if _, err := w.Write([]byte{'\n'}); err != nil {
			return err
		}
	}
	return nil
}
