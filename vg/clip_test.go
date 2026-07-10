package vg

import (
	"strconv"
	"strings"
	"testing"
)

func TestClipRanges(t *testing.T) {
	testCases := [][3]string{
		//         0123456
		[3]string{"", "", ""},
		[3]string{"0", "", ""},
		[3]string{"1", "0:1", "1"},
		[3]string{"00", "", ""},
		[3]string{"01", "0:2", "01"},
		[3]string{"10", "0:2", "10"},
		[3]string{"11", "0:2", "11"},
		[3]string{"0000000", "", ""},
		[3]string{"abcdefg", "0:7", "abcdefg"},
		[3]string{"ab000fg", "0:3,4:7", "ab00fg"},
		[3]string{"00c000g", "1:4,5:7", "0c00g"},
		[3]string{"00def00", "1:6", "0def0"},
		[3]string{"abcde00", "0:6", "abcde0"},
		[3]string{"00cdefg", "1:7", "0cdefg"},
		[3]string{"000d000", "2:5", "0d0"},
		[3]string{"0bc0ef0", "0:7", "0bc0ef0"},
		[3]string{"abc0efg", "0:7", "abc0efg"},
		[3]string{"0b00ef0", "0:7", "0b00ef0"},
		[3]string{"ab00efg", "0:7", "ab00efg"},
		[3]string{"a0b0e0g", "0:7", "a0b0e0g"},
		[3]string{"a00c00g", "0:7", "a00c00g"},
	}
	bs := func(s string) (b []bool) {
		for _, c := range s {
			b = append(b, c != '0')
		}
		return b
	}
	sr := func(a, b []int) string {
		if len(a) != len(b) {
			return "length"
		}
		var r []string
		for i := range a {
			r = append(r, strconv.Itoa(a[i])+":"+strconv.Itoa(b[i]))
		}
		return strings.Join(r, ",")
	}
	tr := func(s string, a, b []int) string {
		r := ""
		for i := range a {
			r += s[a[i]:b[i]]
		}
		return r
	}
	for _, tc := range testCases {
		a, b := clipRanges(bs(tc[0]))
		s := sr(a, b)
		if tc[1] != s {
			t.Fatalf("%s: expected %s got %s", tc[0], tc[1], s)
		}
		q := tr(tc[0], a, b)
		if q != tc[2] {
			t.Fatalf("%s: expected %s got %s", tc[0], tc[2], q)
		}
	}
}
