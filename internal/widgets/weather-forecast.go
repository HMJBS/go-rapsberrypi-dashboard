package widgets

import (
	"dashboard/internal/assets"
	"dashboard/internal/gfx"
	"dashboard/internal/theme"
	"dashboard/internal/weather"
	"fmt"
	"image"
)

// DrawWeatherForecastWidget は、天気予報ウィジェットを描画します。
func DrawWeatherForecastWidget(dst *image.RGBA, forecast weather.Weather, icon image.Image, wOK bool, wErr string) {

	if !wOK {
		drawErrorText(dst, wErr)
		return
	}

	drawIcon(dst, icon)
	drawTemperatureText(dst, forecast.TempC)
	drawWeatherText(dst, weather.CodeLabel(forecast.Code))
}

func drawIcon(dst *image.RGBA, icon image.Image) {

	scaledIcon := gfx.ScaleImage(icon, 1.5)
	gfx.DrawImage(dst, scaledIcon, 930, 530)
}

func drawErrorText(dst *image.RGBA, err string) {
	x := 1030
	y := 850
	face := assets.ClockDateFont
	gray := theme.DefaultTheme.ClockWidgetDateColor
	gfx.DrawText(dst, err, x, y, face, gray, gfx.Centralized)
}

func drawTemperatureText(dst *image.RGBA, tempC float64) {
	x := 1030
	y := 800

	face := assets.ClockTimeFont
	gray := theme.DefaultTheme.ClockWidgetTimeColor
	gfx.DrawText(dst, fmt.Sprintf("%.1f°C", tempC), x, y, face, gray, gfx.Centralized)
}

func drawWeatherText(dst *image.RGBA, text string) {
	x := 1030
	y := 850

	face := assets.ClockDateFont
	gray := theme.DefaultTheme.ClockWidgetDateColor
	gfx.DrawText(dst, text, x, y, face, gray, gfx.Centralized)
}
