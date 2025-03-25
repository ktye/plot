package font

import (
	"archive/zip"
	"bytes"
	"fmt"
	"io/ioutil"
	"strings"
	"testing"
)

func TestFonts(t *testing.T) {
	// File base must fit both, package name and ttf file name.
	files := []string{} // {"file.ttf"}

	for _, name := range files {
		b, err := ioutil.ReadFile(name)
		if err != nil {
			t.Fatal(err)
		}
		b, err = pack(b)
		if err != nil {
			t.Fatal(err)
		}
		writeSrc(strings.TrimSuffix(name, ".ttf"), string(b))
	}
}

func pack(b []byte) ([]byte, error) {
	buf := new(bytes.Buffer)
	w := zip.NewWriter(buf)
	f, err := w.Create("font.ttf")
	if err != nil {
		return nil, err
	}
	_, err = f.Write(b)
	if err != nil {
		return nil, err
	}
	if err = w.Close(); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func writeSrc(name string, s string) {
	fmt.Printf("package %s\n\nconst Data = %q\n", name, s)
}
