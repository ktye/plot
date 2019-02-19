package plot

import "image"

// rotate creates an image which is rotated counter-clockwise by 90 degree.
func rotate(src *image.RGBA) *image.RGBA {
	srcW := src.Bounds().Max.X
	srcH := src.Bounds().Max.Y
	dstW := srcH
	dstH := srcW
	dst := image.NewRGBA(image.Rect(0, 0, dstW, dstH))

	for dstY := 0; dstY < dstH; dstY++ {
		for dstX := 0; dstX < dstW; dstX++ {
			srcX := dstH - dstY - 1
			srcY := dstX

			srcOff := srcY*src.Stride + srcX*4
			dstOff := dstY*dst.Stride + dstX*4

			copy(dst.Pix[dstOff:dstOff+4], src.Pix[srcOff:srcOff+4])
		}
	}

	return dst
}
