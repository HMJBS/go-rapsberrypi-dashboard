package widgets

import (
	"dashboard/internal/gfx"
	"dashboard/internal/theme"
	"image"
)

// DrawPhotoViewerWidget は、写真表示ウィジェットを描画します。
func DrawPhotoViewerWidget(dst *image.RGBA, photo *image.RGBA, pErr string) {

	if photo != nil {
		rect := image.Rect(
			theme.DefaultTheme.PhotoViewerAreaOriginX, theme.DefaultTheme.PhotoViewerAreaOriginY,
			theme.DefaultTheme.PhotoViewerAreaEndX, theme.DefaultTheme.PhotoViewerAreaEndY,
		)
		gfx.DrawImage(dst, photo, rect, gfx.ImageFitContain)
	} else if pErr != "" {
		drawErrorText(dst, pErr)
	}
}
