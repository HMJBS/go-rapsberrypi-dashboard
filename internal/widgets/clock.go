package widgets

import (
	"dashboard/internal/assets"
	"dashboard/internal/gfx"
	"dashboard/internal/theme"
	"image"
	"time"
)

// DrawClockWidget は、時計ウィジェットを描画します。
// ここでは、時計ウィジェットの描画コードを実装します。
func DrawClockWidget(dst *image.RGBA, now time.Time) {
	// 日付を描画
	drawDate(dst, now)

	// 時刻を描画
	drawTime(dst, now)

	// ロケール文字列を描画
	drawLocation(dst, now.Location().String())
}

func drawDate(dst *image.RGBA, now time.Time) {

	x := theme.DefaultTheme.ClockWidgetDateX
	y := theme.DefaultTheme.ClockWidgetDateY
	face := assets.ClockDateFont
	dateStr := now.Format("2006/01/02")
	gray := theme.DefaultTheme.ClockWidgetDateColor
	gfx.DrawText(dst, dateStr, x, y, face, gray, gfx.Normal)
}

func drawTime(dst *image.RGBA, now time.Time) {

	x := theme.DefaultTheme.ClockWidgetTimeX
	y := theme.DefaultTheme.ClockWidgetTimeY
	face := assets.ClockTimeFont
	timeStr := now.Format("15:04:05")
	black := theme.DefaultTheme.ClockWidgetTimeColor
	gfx.DrawText(dst, timeStr, x, y, face, black, gfx.Normal)
}

func drawLocation(dst *image.RGBA, location string) {

	x := theme.DefaultTheme.ClockWidgetLocaleX
	y := theme.DefaultTheme.ClockWidgetLocaleY
	face := assets.ClockDateFont
	gray := theme.DefaultTheme.ClockWidgetLocaleColor
	gfx.DrawText(dst, location, x, y, face, gray, gfx.Normal)
}
