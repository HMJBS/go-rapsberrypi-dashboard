// Package widgets はダッシュボード表示用の簡易ウィジェットを提供します。
package widgets

import (
	"image"
	"image/color"
	"time"

	"dashboard/internal/gfx"
)

const (
	segA = 1 << iota
	segB
	segC
	segD
	segE
	segF
	segG
)

var digitSegments = [10]uint8{
	0: segA | segB | segC | segD | segE | segF,
	1: segB | segC,
	2: segA | segB | segG | segE | segD,
	3: segA | segB | segG | segC | segD,
	4: segF | segG | segB | segC,
	5: segA | segF | segG | segC | segD,
	6: segA | segF | segG | segE | segC | segD,
	7: segA | segB | segC,
	8: segA | segB | segC | segD | segE | segF | segG,
	9: segA | segB | segC | segD | segF | segG,
}

// Draw7SegClock draws HH:MM:SS in a seven-segment style.
func Draw7SegClock(img *image.RGBA, x, y int, digitW, digitH, thickness, gap int, c color.RGBA, now time.Time) {
	h := now.Hour()
	m := now.Minute()
	s := now.Second()

	digits := [6]int{h / 10, h % 10, m / 10, m % 10, s / 10, s % 10}

	cx := x
	for i := 0; i < 6; i++ {
		drawDigit(img, cx, y, digitW, digitH, thickness, c, digits[i])
		cx += digitW + gap
		if i == 1 || i == 3 {
			drawColon(img, cx-gap/2, y, digitW/5, digitH, thickness, c)
			cx += (digitW / 3)
		}
	}
}

func drawDigit(img *image.RGBA, x, y, w, h, t int, c color.RGBA, d int) {
	if d < 0 || d > 9 {
		return
	}
	mask := digitSegments[d]
	midY0 := y + (h/2 - t/2)
	midY1 := midY0 + t

	if mask&segA != 0 {
		gfx.RectRGBA(img, x+t, y, x+w-t, y+t, c)
	}
	if mask&segG != 0 {
		gfx.RectRGBA(img, x+t, midY0, x+w-t, midY1, c)
	}
	if mask&segD != 0 {
		gfx.RectRGBA(img, x+t, y+h-t, x+w-t, y+h, c)
	}

	upperY0 := y + t
	upperY1 := y + h/2
	lowerY0 := y + h/2
	lowerY1 := y + h - t

	if mask&segF != 0 {
		gfx.RectRGBA(img, x, upperY0, x+t, upperY1, c)
	}
	if mask&segB != 0 {
		gfx.RectRGBA(img, x+w-t, upperY0, x+w, upperY1, c)
	}
	if mask&segE != 0 {
		gfx.RectRGBA(img, x, lowerY0, x+t, lowerY1, c)
	}
	if mask&segC != 0 {
		gfx.RectRGBA(img, x+w-t, lowerY0, x+w, lowerY1, c)
	}
}

func drawColon(img *image.RGBA, x, y, w, h, t int, c color.RGBA) {
	s := t
	if s < 2 {
		s = 2
	}
	cx0 := x
	cx1 := x + w
	cy1 := y + h/3
	cy2 := y + 2*h/3

	gfx.RectRGBA(img, cx0, cy1, cx1, cy1+s, c)
	gfx.RectRGBA(img, cx0, cy2, cx1, cy2+s, c)
}
