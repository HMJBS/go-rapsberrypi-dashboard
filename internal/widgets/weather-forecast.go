package widgets

import (
	"dashboard/internal/assets"
	"dashboard/internal/gfx"
	"dashboard/internal/theme"
	"dashboard/internal/weather"
	"fmt"
	"image"

	"golang.org/x/image/font"
	"golang.org/x/image/math/fixed"
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
	x := 900
	y := 700
	face := assets.ClockTimeFont
	gray := theme.DefaultTheme.ClockWidgetDateColor
	drawer := font.Drawer{
		Dst:  dst,
		Src:  image.NewUniform(gray),
		Face: face,
		Dot:  fixed.P(x, y),
	}
	drawer.DrawString(err)
}

func drawTemperatureText(dst *image.RGBA, tempC float64) {
	x := 930
	y := 800

	face := assets.ClockTimeFont
	gray := theme.DefaultTheme.ClockWidgetTimeColor
	drawer := font.Drawer{
		Dst:  dst,
		Src:  image.NewUniform(gray),
		Face: face,
		Dot:  fixed.P(x, y),
	}
	drawer.DrawString(fmt.Sprintf("%.1f°C", tempC))
}

func drawWeatherText(dst *image.RGBA, text string) {
	x := 930
	y := 880

	face := assets.ClockDateFont
	gray := theme.DefaultTheme.ClockWidgetDateColor
	drawer := font.Drawer{
		Dst:  dst,
		Src:  image.NewUniform(gray),
		Face: face,
		Dot:  fixed.P(x, y),
	}
	drawer.DrawString(text)
}
