// +build ignore

package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

const help = `(ktye/plot/cmd/csv.go)

csv < data.csv

csv fff             # 3 float columns, ignore remaining
csv f*              # all remaining floats
csv f3              # same
csv fz              # float abs angle, e.g 1.2 3.0 40 -> 1.2 3.0a40
csv f@              # float abs angle,                   1.2 3.0@40
csv f-f             # skip
csv h2f5            # skip 2 header lines..
csv nf              # missing/parse as 0n instead of trap

csv ,fff            # separator is , (default spaces/tabs)


f(float) a|@(polar) s(string) q("quoted") ` + "`" + `(` + "`" + `symbol)
`

var ad bool // 3@40
var hdr int
var nan bool // 0n and continue
var rem byte // *
var sep byte
var f []byte

func main() {
	if len(os.Args) != 2 {
		fatal(help)
	}
	a := os.Args[1]
	if len(a) > 0 {
		s := ",; \t"
		if i := strings.Index(a, s); i != -1 {
			sep, a = s[i], a[1:]
		}
	}
	if len(a) > 0 && a[0] == 'n' {
		nan, a = true, a[1:]
	}
	for len(a) > 0 {
		switch a[0] {
		case 'f':
			f = append(f, 'f')
		case 'z', 'a':
			f = append(f, 'z')
		case '@':
			f = append(f, '@')
		case '-':
			f = append(f, '-')
		case 's', 'c':
			f = append(f, 's')
		case 'q':
			f = append(f, 'q')
		case '`':
			f = append(f, '`')
		case '*':
			if len(f) == 0 {
				fatal(help)
			}
			rem = f[len(f)-1]
		case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
			e := 1
			for i := range a {
				if a[i] < '0' || a[i] > '9' {
					e = i
					break
				}
			}
			n, err := strconv.Atoi(a[:e])
			if err != nil {
				fatal("?")
			}
			if len(f) == 0 {
				fatal(help)
			}
			for i := 0; i < n; i++ {
				f = append(f, f[len(f)-1])
			}
		default:
			fatal(help)
		}
	}
	read()
}
func read() {
	scn := bufio.NewScanner(os.Stdin)
	num := -1
	for scn.Scan() {
		linenumber++
		t := scn.Text()
		if hdr > 0 {
			hdr--
			continue
		}
		var v []string
		if sep == 0 {
			v = strings.Fields(t)
		} else {
			v = strings.Split(t, string(sep))
		}
		if num < 0 {
			num = len(v)
		} else if len(v) == 0 {
			continue // last line may be empty
		} else if len(v) != num {
			fatal("csv: input is not rectangular")
		}
		write(v)
	}
}
func write(v []string) {
	for i := range f {
		if i > 0 {
			fmt.Print(" ") // OFS
		}
		format(v, i, f[i])
	}
	if rem != 0 {
		for i := len(f); i < len(v); i++ {
			fmt.Print(" ") // OFS
			format(v, i, rem)
		}
	}
}
func format(v []string, i int, c byte) {
	switch c {
	case 'f':
		fmt.Print(atof(v, i))
	case 'z', '@':
		a, b := atof(v, i), atof(v, i+1)
		if b[0] == '-' {
			n, _ := strconv.ParseFloat(b, 64)
			b = strconv.FormatFloat(n+360, 'f', -1, 64)
		}
		fmt.Printf("%v%c%v", a, c, b)
	case '-':
	case 's':
		fmt.Print(v[i])
	case 'q':
		fmt.Printf("%q", v[i])
	case '`':
		fmt.Print("`" + v[i])
	}
}
func atof(v []string, i int) string {
	if i >= len(v) {
		return na()
	}
	s := strings.Replace(v[i], ",", ".", -1)
	if _, err := strconv.ParseFloat(s, 64); err != nil {
		return na()
	}
	return s
}
func na() string {
	if nan {
		return "0n"
	}
	panic("not enough columns for format. line " + strconv.Itoa(linenumber))
}
func fatal(s string) { fmt.Fprintln(os.Stderr, s); os.Exit(1) }

var linenumber int
