module github.com/ktye/plot/serve

go 1.24

toolchain go1.24.0

require (
	github.com/ktye/plot v0.0.0
	github.com/ktye/plot/clipboard v0.0.0
)

require (
	github.com/golang/freetype v0.0.0-20170609003504-e2365dfdc4a0 // indirect
	github.com/lxn/win v0.0.0-20201111105847-2a20daff6a55 // indirect
	golang.org/x/image v0.0.0-20200119044424-58c23975cae1 // indirect
	golang.org/x/sys v0.0.0-20201018230417-eeed37f84f13 // indirect
)

replace github.com/ktye/plot => ../

replace github.com/ktye/plot/clipboard => ../clipboard
