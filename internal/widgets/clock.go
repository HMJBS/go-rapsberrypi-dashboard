package widgets

import (
	"dashboard/internal/assets"
	"dashboard/internal/theme"
	"image"
	"time"

	"golang.org/x/image/font"
	"golang.org/x/image/math/fixed"
)

// DrawClockWidget は、時計ウィジェットを描画します。
// ここでは、時計ウィジェットの描画コードを実装します。
func DrawClockWidget(dst *image.RGBA, now time.Time, locale string) {
	// 日付を描画
	drawDate(dst, now)

	// 時刻を描画
	drawTime(dst, now)

	// ロケール文字列を描画
	drawLocale(dst, locale)
}

func drawDate(dst *image.RGBA, now time.Time) {

	x := theme.DefaultTheme.ClockWidgetDateX
	y := theme.DefaultTheme.ClockWidgetDateY
	face := assets.ClockDateFont
	dateStr := now.Format("2006/01/02")
	gray := theme.DefaultTheme.ClockWidgetDateColor
	drawer := font.Drawer{
		Dst:  dst,
		Src:  image.NewUniform(gray),
		Face: face,
		Dot:  fixed.P(x, y),
	}
	drawer.DrawString(dateStr)
}

func drawTime(dst *image.RGBA, now time.Time) {

	x := theme.DefaultTheme.ClockWidgetTimeX
	y := theme.DefaultTheme.ClockWidgetTimeY
	face := assets.ClockTimeFont
	timeStr := now.Format("15:04:05")
	black := theme.DefaultTheme.ClockWidgetTimeColor
	drawer := font.Drawer{
		Dst:  dst,
		Src:  image.NewUniform(black),
		Face: face,
		Dot:  fixed.P(x, y),
	}
	drawer.DrawString(timeStr)
}

func drawLocale(dst *image.RGBA, locale string) {

	x := theme.DefaultTheme.ClockWidgetLocaleX
	y := theme.DefaultTheme.ClockWidgetLocaleY
	face := assets.ClockDateFont
	gray := theme.DefaultTheme.ClockWidgetLocaleColor
	drawer := font.Drawer{
		Dst:  dst,
		Src:  image.NewUniform(gray),
		Face: face,
		Dot:  fixed.P(x, y),
	}
	drawer.DrawString(locale)
}
