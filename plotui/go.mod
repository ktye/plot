module github.com/ktye/plot/plotui

go 1.24

toolchain go1.24.0

require (
	github.com/ktye/plot v0.0.0
	github.com/ktye/plot/clipboard v0.0.0
	github.com/lxn/walk v0.0.0-20210112085537-c389da54e794
	github.com/lxn/win v0.0.0-20201111105847-2a20daff6a55
)

require (
	github.com/golang/freetype v0.0.0-20170609003504-e2365dfdc4a0 // indirect
	golang.org/x/image v0.0.0-20201208152932-35266b937fa6 // indirect
	golang.org/x/sys v0.0.0-20210124154548-22da62e12c0c // indirect
	gopkg.in/Knetic/govaluate.v3 v3.0.0 // indirect
)

replace github.com/ktye/plot => ../

replace github.com/ktye/plot/clipboard => ../clipboard

replace github.com/ktye/pptx => ../../pptx

replace github.com/ktye/pptx/pptxt => ../../pptx/pptxt
