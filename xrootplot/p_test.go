package main

import (
	"io"
	"strings"
	"testing"
)

const tcX = `1.1
2.2
-1.5
`
const tcXY12 = `0.1 1.1 3.2
0.2 2.2 0
0.3 -1.5 1.5

0 3.6 2.1
0.1 -1.5 3.234
0.2 2.22 3.33
0.4 1 2
`
const tcXZ12 = `10 3.21a20
20 3.21a22
30 3.21a25

10 1.22a5
20 1.20a5
30 1.30a6
`
const tcZ12 = `3.21a20
3.21a22
3.21a25

1.22a5
1.20a5
1.30a6
`

func TestX(t *testing.T) {
	do(sr(tcX))
	do(sr(tcXY12))
	do(sr(tcXZ12))
	do(sr(tcZ12))
}
func sr(s string) io.Reader { return strings.NewReader(s) }
