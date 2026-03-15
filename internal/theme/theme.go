// Package theme は、デザインに関係するコードや座標情報を提供
package theme

import "image/color"

type theme struct {
	BackgroundColor color.RGBA
}

// DefaultTheme は既定のテーマ
var DefaultTheme = theme{
	BackgroundColor: color.RGBA{R: 241, G: 245, B: 249, A: 255},
}
