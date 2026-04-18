// Package gfx は RGBA バッファ向けの最小描画ユーティリティです。
package gfx

import (
	"image"
	"image/color"

	"golang.org/x/image/draw"
	"golang.org/x/image/font"
	"golang.org/x/image/math/fixed"
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

// DrawImage は img 上に src を描画します。
func DrawImage(dst *image.RGBA, src image.Image, x, y int) {
	srcBounds := src.Bounds()
	drawRect := image.Rect(x, y, x+srcBounds.Dx(), y+srcBounds.Dy()).Intersect(dst.Bounds())
	if drawRect.Empty() {
		return
	}

	srcStartX := srcBounds.Min.X + (drawRect.Min.X - x)
	srcStartY := srcBounds.Min.Y + (drawRect.Min.Y - y)
	draw.Draw(dst, drawRect, src, image.Pt(srcStartX, srcStartY), draw.Over)
}

// ScaleImage は src を scale 倍して返します。
func ScaleImage(src image.Image, scale float64) image.Image {
	srcBounds := src.Bounds()
	dstWidth := int(float64(srcBounds.Dx()) * scale)
	dstHeight := int(float64(srcBounds.Dy()) * scale)
	dst := image.NewRGBA(image.Rect(0, 0, dstWidth, dstHeight))
	draw.ApproxBiLinear.Scale(dst, dst.Bounds(), src, srcBounds, draw.Over, nil)
	return dst
}

// TextDrawMode は DrawText の描画モードを表します。
type TextDrawMode int

// DrawText の描画モード定数です。
const (
	Normal      TextDrawMode = iota // Normal は、指定した座標をテキストの左下とする描画モードです。
	Centralized                     // Centralized は、指定した座標をテキストの水平中央とする描画モードです。
)

// DrawText は img 上に text を描画します。
func DrawText(
	dst *image.RGBA, text string, x, y int, face font.Face,
	c color.Color, mode TextDrawMode,
) {
	switch mode {
	case Normal:
		drawer := font.Drawer{
			Dst:  dst,
			Src:  image.NewUniform(c),
			Face: face,
			Dot:  fixed.P(x, y),
		}
		drawer.DrawString(text)
	case Centralized:
		textWidth := font.MeasureString(face, text).Round()
		drawer := font.Drawer{
			Dst:  dst,
			Src:  image.NewUniform(c),
			Face: face,
			Dot:  fixed.P(x-textWidth/2, y),
		}
		drawer.DrawString(text)
	default:
		panic("unknown TextDrawMode")
	}
}
