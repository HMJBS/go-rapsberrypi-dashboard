package gfx

import (
	"image"
	"image/color"
	"testing"
)

func TestDrawImageAlphaBlendsSource(t *testing.T) {
	dst := image.NewRGBA(image.Rect(0, 0, 1, 1))
	dst.SetRGBA(0, 0, color.RGBA{R: 10, G: 20, B: 30, A: 255})

	src := image.NewNRGBA(image.Rect(0, 0, 1, 1))
	src.SetNRGBA(0, 0, color.NRGBA{R: 200, G: 100, B: 0, A: 128})

	DrawImage(dst, src, 0, 0)

	got := dst.RGBAAt(0, 0)
	want := color.RGBA{R: 105, G: 60, B: 14, A: 255}
	if got != want {
		t.Fatalf("RGBAAt(0, 0) = %#v, want %#v", got, want)
	}
}

func TestDrawImageKeepsDestinationForTransparentSource(t *testing.T) {
	dst := image.NewRGBA(image.Rect(0, 0, 1, 1))
	dst.SetRGBA(0, 0, color.RGBA{R: 10, G: 20, B: 30, A: 255})

	src := image.NewNRGBA(image.Rect(0, 0, 1, 1))
	src.SetNRGBA(0, 0, color.NRGBA{R: 200, G: 100, B: 0, A: 0})

	DrawImage(dst, src, 0, 0)

	got := dst.RGBAAt(0, 0)
	want := color.RGBA{R: 10, G: 20, B: 30, A: 255}
	if got != want {
		t.Fatalf("RGBAAt(0, 0) = %#v, want %#v", got, want)
	}
}

func TestDrawImageRespectsDestinationOffset(t *testing.T) {
	dst := image.NewRGBA(image.Rect(0, 0, 4, 4))
	src := image.NewRGBA(image.Rect(0, 0, 2, 2))

	src.SetRGBA(0, 0, color.RGBA{R: 1, G: 2, B: 3, A: 255})
	src.SetRGBA(1, 0, color.RGBA{R: 4, G: 5, B: 6, A: 255})
	src.SetRGBA(0, 1, color.RGBA{R: 7, G: 8, B: 9, A: 255})
	src.SetRGBA(1, 1, color.RGBA{R: 10, G: 11, B: 12, A: 255})

	DrawImage(dst, src, 1, 1)

	if got := dst.RGBAAt(0, 0); got != (color.RGBA{}) {
		t.Fatalf("RGBAAt(0, 0) = %#v, want zero", got)
	}
	if got := dst.RGBAAt(1, 1); got != (color.RGBA{R: 1, G: 2, B: 3, A: 255}) {
		t.Fatalf("RGBAAt(1, 1) = %#v", got)
	}
	if got := dst.RGBAAt(2, 1); got != (color.RGBA{R: 4, G: 5, B: 6, A: 255}) {
		t.Fatalf("RGBAAt(2, 1) = %#v", got)
	}
	if got := dst.RGBAAt(1, 2); got != (color.RGBA{R: 7, G: 8, B: 9, A: 255}) {
		t.Fatalf("RGBAAt(1, 2) = %#v", got)
	}
	if got := dst.RGBAAt(2, 2); got != (color.RGBA{R: 10, G: 11, B: 12, A: 255}) {
		t.Fatalf("RGBAAt(2, 2) = %#v", got)
	}
}

func TestDrawImageClipsToDestinationBounds(t *testing.T) {
	dst := image.NewRGBA(image.Rect(10, 20, 13, 23))
	src := image.NewRGBA(image.Rect(0, 0, 3, 3))

	for y := 0; y < 3; y++ {
		for x := 0; x < 3; x++ {
			src.SetRGBA(x, y, color.RGBA{R: uint8(10 + x), G: uint8(20 + y), B: 30, A: 255})
		}
	}

	DrawImage(dst, src, 12, 22)

	if got := dst.RGBAAt(12, 22); got != (color.RGBA{R: 10, G: 20, B: 30, A: 255}) {
		t.Fatalf("RGBAAt(12, 22) = %#v", got)
	}
	if got := dst.RGBAAt(10, 20); got != (color.RGBA{}) {
		t.Fatalf("RGBAAt(10, 20) = %#v, want zero", got)
	}
}
