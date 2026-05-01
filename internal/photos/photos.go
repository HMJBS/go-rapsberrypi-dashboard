// Package photos は、ローカルディレクトリに同期済みの画像を読み込み、画面サイズへ整形します。
package photos

import (
	"fmt"
	"image"
	"image/draw"
	_ "image/jpeg" // JPEG デコードを有効化
	_ "image/png"  // PNG デコードを有効化
	"os"
	"path/filepath"
	"strings"
)

// ListImages lists supported image files in dir (non-recursive).
func ListImages(dir string) ([]string, error) {
	photos, err := os.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("readdir %s: %w", dir, err)
	}
	out := make([]string, 0, len(photos))
	for _, e := range photos {
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
func LoadScreenImage(path string) (*image.RGBA, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("open image: %w", err)
	}
	defer func() { _ = f.Close() }()

	src, _, err := image.Decode(f)
	if err != nil {
		return nil, fmt.Errorf("decode image: %w", err)
	}

	return toRGBA(src), nil
}

func toRGBA(src image.Image) *image.RGBA {
	b := src.Bounds()
	dst := image.NewRGBA(image.Rect(0, 0, b.Dx(), b.Dy()))
	draw.Draw(dst, dst.Bounds(), src, b.Min, draw.Src)
	return dst
}
