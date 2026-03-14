// Package photos は、ローカルディレクトリに同期済みの画像を読み込み、画面サイズへ整形します。
package photos

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	_ "image/jpeg" // JPEG デコードを有効化
	_ "image/png"  // PNG デコードを有効化
	"math"
	"os"
	"path/filepath"
	"strings"
)

// ListImages lists supported image files in dir (non-recursive).
func ListImages(dir string) ([]string, error) {
	ents, err := os.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("readdir %s: %w", dir, err)
	}
	out := make([]string, 0, len(ents))
	for _, e := range ents {
		if e.IsDir() {
			continue
		}
		n := strings.ToLower(e.Name())
		if !strings.HasSuffix(n, ".jpg") && !strings.HasSuffix(n, ".jpeg") && !strings.HasSuffix(n, ".png") {
			continue
		}
		out = append(out, filepath.Join(dir, e.Name()))
	}
	return out, nil
}

// LoadScreenImage decodes path, scales it to fit within (screenW, screenH) while
// preserving aspect ratio, and composites it centered on a screen-sized RGBA filled with bg.
func LoadScreenImage(path string, screenW, screenH int, bg color.RGBA) (*image.RGBA, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("open image: %w", err)
	}
	defer func() { _ = f.Close() }()

	src, _, err := image.Decode(f)
	if err != nil {
		return nil, fmt.Errorf("decode image: %w", err)
	}

	srcB := src.Bounds()
	sw := srcB.Dx()
	sh := srcB.Dy()
	if sw <= 0 || sh <= 0 {
		return nil, fmt.Errorf("invalid image bounds")
	}

	scale := math.Min(float64(screenW)/float64(sw), float64(screenH)/float64(sh))
	if scale <= 0 {
		return nil, fmt.Errorf("invalid scale")
	}

	newW := int(math.Round(float64(sw) * scale))
	newH := int(math.Round(float64(sh) * scale))
	if newW < 1 {
		newW = 1
	}
	if newH < 1 {
		newH = 1
	}

	srcRGBA := toRGBA(src)
	scaled := scaleNearestRGBA(srcRGBA, newW, newH)

	screen := image.NewRGBA(image.Rect(0, 0, screenW, screenH))
	fill(screen, bg)

	offX := (screenW - newW) / 2
	offY := (screenH - newH) / 2
	draw.Draw(screen, image.Rect(offX, offY, offX+newW, offY+newH), scaled, image.Point{}, draw.Src)
	return screen, nil
}

func toRGBA(src image.Image) *image.RGBA {
	b := src.Bounds()
	dst := image.NewRGBA(image.Rect(0, 0, b.Dx(), b.Dy()))
	draw.Draw(dst, dst.Bounds(), src, b.Min, draw.Src)
	return dst
}

func scaleNearestRGBA(src *image.RGBA, w, h int) *image.RGBA {
	dst := image.NewRGBA(image.Rect(0, 0, w, h))
	sw := src.Bounds().Dx()
	sh := src.Bounds().Dy()

	for y := 0; y < h; y++ {
		sy := (y * sh) / h
		srcRow := src.Pix[sy*src.Stride : sy*src.Stride+sw*4]
		dstRow := dst.Pix[y*dst.Stride : y*dst.Stride+w*4]
		for x := 0; x < w; x++ {
			sx := (x * sw) / w
			si := sx * 4
			di := x * 4
			dstRow[di+0] = srcRow[si+0]
			dstRow[di+1] = srcRow[si+1]
			dstRow[di+2] = srcRow[si+2]
			dstRow[di+3] = srcRow[si+3]
		}
	}
	return dst
}

func fill(img *image.RGBA, c color.RGBA) {
	b := img.Bounds()
	w := b.Dx()
	h := b.Dy()
	for y := 0; y < h; y++ {
		row := img.Pix[y*img.Stride : y*img.Stride+w*4]
		for x := 0; x < w; x++ {
			i := x * 4
			row[i+0] = c.R
			row[i+1] = c.G
			row[i+2] = c.B
			row[i+3] = c.A
		}
	}
}
