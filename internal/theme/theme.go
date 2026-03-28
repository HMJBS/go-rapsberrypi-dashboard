// Package theme は、デザインに関係するコードや座標情報を提供
package theme

import "image/color"

// Theme は、ダッシュボードのテーマ設定を表します。
type Theme struct {
	BackgroundColor            color.RGBA
	ClockWidgetBackgroundColor color.RGBA
	ClockWidgetDateX           int
	ClockWidgetDateY           int
	ClockWidgetDateColor       color.RGBA
	ClockWidgetTimeX           int
	ClockWidgetTimeY           int
	ClockWidgetTimeColor       color.RGBA
	ClockWidgetLocaleX         int
	ClockWidgetLocaleY         int
	ClockWidgetLocaleColor     color.RGBA
}

// DefaultTheme は既定のテーマ
var DefaultTheme = Theme{
	BackgroundColor:            color.RGBA{R: 241, G: 245, B: 249, A: 255},
	ClockWidgetBackgroundColor: color.RGBA{R: 255, G: 255, B: 255, A: 255},
	ClockWidgetDateX:           880,
	ClockWidgetDateY:           202,
	ClockWidgetDateColor:       color.RGBA{R: 100, G: 116, B: 139, A: 255},
	ClockWidgetTimeX:           880,
	ClockWidgetTimeY:           280,
	ClockWidgetTimeColor:       color.RGBA{R: 0, G: 0, B: 0, A: 255},
	ClockWidgetLocaleX:         880,
	ClockWidgetLocaleY:         330,
	ClockWidgetLocaleColor:     color.RGBA{R: 100, G: 116, B: 139, A: 255},
}
