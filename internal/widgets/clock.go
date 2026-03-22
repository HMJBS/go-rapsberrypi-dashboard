package widgets

import (
	"dashboard/internal/assets"
	"image"
	"image/color"
	"time"

	"golang.org/x/image/font"
	"golang.org/x/image/math/fixed"
)

// DrawClockWidget は、時計ウィジェットを描画します。
// ここでは、時計ウィジェットの描画コードを実装します。
func DrawClockWidget(dst *image.RGBA, now time.Time) {
	// 日付を描画
	drawDate(dst, now)

	// 時刻を描画
	drawTime(dst, now)
}

func drawDate(dst *image.RGBA, now time.Time) {

	x := 860
	y := 202
	face := assets.ClockDateFont
	dateStr := now.Format("2006/01/02")
	gray := color.RGBA{R: 100, G: 116, B: 139, A: 255}
	drawer := font.Drawer{
		Dst:  dst,
		Src:  image.NewUniform(gray),
		Face: face,
		Dot:  fixed.P(x, y),
	}
	drawer.DrawString(dateStr)
}

func drawTime(dst *image.RGBA, now time.Time) {

	x := 860
	y := 300
	face := assets.ClockTimeFont
	timeStr := now.Format("15:04:05")
	black := color.RGBA{R: 0, G: 0, B: 0, A: 255}
	drawer := font.Drawer{
		Dst:  dst,
		Src:  image.NewUniform(black),
		Face: face,
		Dot:  fixed.P(x, y),
	}
	drawer.DrawString(timeStr)
}
