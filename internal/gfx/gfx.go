// Package gfx は RGBA バッファ向けの最小描画ユーティリティです。
package gfx

import (
	"image"
	"image/color"
)

// FillRGBA は img 全体を単色で塗りつぶします。
func FillRGBA(img *image.RGBA, c color.RGBA) {
	r := img.Bounds()
	w := r.Dx()
	h := r.Dy()
	for y := 0; y < h; y++ {
		row := img.Pix[y*img.Stride : y*img.Stride+w*4]
		for x := 0; x < w; x++ {
			i := x * 4
			row[i+0] = c.R
			row[i+1] = c.G
			row[i+2] = c.B
			row[i+3] = c.A
		}
	}
}

// RectRGBA は img 上に塗りつぶし矩形を描画します。
func RectRGBA(img *image.RGBA, x0, y0, x1, y1 int, c color.RGBA) {
	b := img.Bounds()
	if x0 < b.Min.X {
		x0 = b.Min.X
	}
	if y0 < b.Min.Y {
		y0 = b.Min.Y
	}
	if x1 > b.Max.X {
		x1 = b.Max.X
	}
	if y1 > b.Max.Y {
		y1 = b.Max.Y
	}
	if x0 >= x1 || y0 >= y1 {
		return
	}

	for y := y0; y < y1; y++ {
		row := img.Pix[y*img.Stride : y*img.Stride+x1*4]
		for x := x0; x < x1; x++ {
			i := x * 4
			row[i+0] = c.R
			row[i+1] = c.G
			row[i+2] = c.B
			row[i+3] = c.A
		}
	}
}
