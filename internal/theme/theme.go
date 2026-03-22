// Package theme は、デザインに関係するコードや座標情報を提供
package theme

import "image/color"

type theme struct {
	BackgroundColor            color.RGBA
	ClockWidgetBackgroundColor color.RGBA
	ClockWidgetX               int
	ClockWidgetY               int
	ClockWidgetWidth           int
	ClockWidgetHeight          int
	ClockWidgetRadius          int
}

// DefaultTheme は既定のテーマ
var DefaultTheme = theme{
	BackgroundColor:            color.RGBA{R: 241, G: 245, B: 249, A: 255},
	ClockWidgetBackgroundColor: color.RGBA{R: 255, G: 255, B: 255, A: 255},
	ClockWidgetWidth:           400,
	ClockWidgetHeight:          376,
	ClockWidgetRadius:          24,
}
