package widgets

import (
	"dashboard/internal/gfx"
	"image"
)

// DrawPhotoViewerWidget は、写真表示ウィジェットを描画します。
func DrawPhotoViewerWidget(dst *image.RGBA, photo image.Image, pErr string) {

	PhotoAreaOriginX := 51
	PhotoAreaOriginY := 51
	PhotoAreaEndX := 805
	PhotoAreaEndY := 973

	if photo != nil {
		rect := image.Rect(
			PhotoAreaOriginX, PhotoAreaOriginY,
			PhotoAreaEndX, PhotoAreaEndY,
		)
		gfx.DrawImage(dst, photo, rect, gfx.ImageFitContain)
	} else if pErr != "" {
		drawErrorText(dst, pErr)
	}
}
