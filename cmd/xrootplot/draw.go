package main

import (
	"encoding/base64"
	"os"
)

func draw(b []c) {

	// iterm2.com/documentation-images.html
	// github.com/mintty/mintty/blob/master/src/termout.c ...1337
	// github.com/mintty/utils/blob/master/showimg
	os.Stdout.Write([]byte{27})
	os.Stdout.Write([]byte("]1337;File=:"))
	enc := base64.NewEncoder(base64.StdEncoding, os.Stdout)
	enc.Write(b)
	enc.Close()
	os.Stdout.Write([]byte{7})
}
