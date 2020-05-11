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

csv FORMAT < data.csv

csv ,fff            # use , as separator(first)  default spaces/tabs
csv h2f5            # skip 2 header lines..

csv fff             # 3 float columns, ignore remaining
csv f*              # * repeat last format for all remaining cols
csv f3              # repeat last format 3 times (f3 is fff)
csv a               # abs angle, e.g 3.0 40 -> 3.0a40
csv @               # abs angle,               3.0@40
csv f-f             # - skip 2nd
csv nf              # n missing/parse error as 0n instead of trap


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
		fatal("wrong number of args\n" + help)
	}
	a := os.Args[1]
	if len(a) > 0 {
		s := ",; \t"
		if i := strings.IndexByte(s, a[0]); i != -1 {
			sep, a = s[i], a[1:]
		}
	}
	if len(a) > 0 && a[0] == 'n' {
		nan, a = true, a[1:]
	}
	for len(a) > 0 {
		shift := func(c byte) {
			f = append(f, c)
			a = a[1:]
		}
		switch a[0] {
		case 'h':
			shift('h')
		case 'f':
			shift('f')
		case 'a', 'z':
			shift('a')
		case '@':
			shift('@')
		case '-':
			shift('-')
		case 's', 'c':
			shift('s')
		case 'q':
			shift('q')
		case '`':
			shift('`')
		case '*':
			if len(f) == 0 {
				fatal("* follows empty format\n" + help)
			}
			if len(a) > 1 {
				fatal("* must be last in format\n" + help)
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
				fatal("number must follow a format\n" + help)
			}
			for i := 0; i < n-1; i++ {
				f = append(f, f[len(f)-1])
			}
			a = a[e:]
		default:
			fatal("unknown token: " + string(a[0]) + "\n" + help)
		}
	}
	headers()
	// show()
	read()
}
func headers() { // remove hhh from f(format)
	var ff []byte
	for _, c := range f {
		if c == 'h' {
			hdr++
		} else {
			ff = append(ff, c)
		}
	}
	f = ff
}
func read() {
	debom := func(s string) string { // remove byte-order-mark
		if len(s) >= 3 && s[0] == 0xEF && s[1] == 0xBB && s[2] == 0xBF {
			return s[3:]
		}
		return s
	}
	scn := bufio.NewScanner(os.Stdin)
	num := -1
	for scn.Scan() {
		linenumber++
		t := scn.Text()
		if linenumber == 1 {
			t = debom(t)
		}
		if len(t) > 0 && t[len(t)-1] == '\r' {
			t = t[:len(t)-1]
		}
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
		} else if len(v) == 0 || (len(v) == 1 && v[0] == "") {
			continue // last line may be empty
		} else if len(v) != num {
			es := fmt.Sprintf("line %d: %d cols (not %d)", linenumber, len(v), num)
			fmt.Printf("t: %q\n", t)
			fatal("csv: input is not rectangular\n" + es)
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
	fmt.Print("\n")
}
func format(v []string, i int, c byte) {
	switch c {
	case 'f':
		fmt.Print(atof(v, i))
	case 'a', '@':
		a, b := atof(v, i), atof(v, i+1)
		n, _ := strconv.ParseFloat(b, 64)
		if n < 0 {
			n += 360.0
		}
		b = strconv.FormatFloat(n, 'f', -1, 64) // no E in angle
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
func show() {
	fmt.Println("ad(@):", ad)
	fmt.Println("hdr:", hdr)
	fmt.Println("nan(continue):", nan)
	fmt.Println("rem(*):", string(rem))
	fmt.Println("sep:", string(sep))
	fmt.Println("f(format):", string(f))
}

var linenumber int
