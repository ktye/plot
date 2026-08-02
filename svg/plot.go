package svg

import "github.com/ktye/plot"
import "bytes"
import "encoding/json"
import _ "embed"

//go:embed plot.js
var Plot []byte

func HtmlPlot(p plot.Plots, width, height, columns int, idx []plot.HighlightID) ([]byte, error) {
	var buf bytes.Buffer
	buf.Write([]byte(head))
	buf.Write(Plot)
	buf.Write([]byte(t1))
	buf.Write([]byte("plot("))
	if e := json.NewEncoder(&buf).Encode(p); e != nil {
		return nil, e
	}
	buf.Write([]byte(",0,svg)"))
	buf.Write([]byte(t2))
	return buf.Bytes(), nil
}

const head = `
<!DOCTYPE html>
<head><meta charset="utf-8">
<title>plots</title>
<style>*{font-family:monospace;box-sizing:border-box}
html,body{margin:0;padding:0}
body{width:100vw;height:100vh}
.link{color:#0000ee}
.link:hover{cursor:pointer}
.b{color:#aaa}
.a:hover{cursor:pointer}
.b:hover{cursor:pointer}
</style></head><script>
`

const t1 = `
</script>
<body>
<svg id="svg" style="flex-grow:2;border:1px solid #ddd"></svg>
<script>
`
const t2 = `
</script>
</body></html>
`
