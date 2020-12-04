// +build !windows

package main

func consoleSize() (int, int) { return 0, 0 }
func drawConsole(w, h int, c []c) {
	panic("only on windows")
}
