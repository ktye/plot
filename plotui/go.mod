module github.com/ktye/plot/plotui

go 1.13

require (
	github.com/golang/freetype v0.0.0-20170609003504-e2365dfdc4a0
	github.com/ktye/plot v0.0.0
	github.com/lxn/walk v0.0.0-20191128110447-55ccb3a9f5c1
	github.com/lxn/win v0.0.0-20191128105842-2da648fda5b4
	gopkg.in/Knetic/govaluate.v3 v3.0.0 // indirect
)

replace github.com/ktye/plot => ../

replace github.com/ktye/pptx => ../../pptx

replace github.com/ktye/pptx/pptxt => ../../pptx/pptxt
