// Package previewpng は、描画結果を PNG としてディスクへ出力します。
package previewpng

import (
	"fmt"
	"image"
	"image/png"
	"os"
	"path/filepath"
)

// WriteLatestPNG は dir/latest.png を更新します。
// 一時ファイルへ書き出してからリネームすることで、読み取り側が途中状態を掴みにくくします。
func WriteLatestPNG(dir string, img image.Image) error {
	if dir == "" {
		return fmt.Errorf("preview dir is empty")
	}
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return fmt.Errorf("create preview dir: %w", err)
	}

	tmp, err := os.CreateTemp(dir, "latest-*.png")
	if err != nil {
		return fmt.Errorf("create temp file: %w", err)
	}
	tmpName := tmp.Name()

	encodeErr := png.Encode(tmp, img)
	closeErr := tmp.Close()
	if encodeErr != nil {
		_ = os.Remove(tmpName)
		return fmt.Errorf("encode png: %w", encodeErr)
	}
	if closeErr != nil {
		_ = os.Remove(tmpName)
		return fmt.Errorf("close temp file: %w", closeErr)
	}

	final := filepath.Join(dir, "latest.png")
	if err := os.Rename(tmpName, final); err != nil {
		_ = os.Remove(tmpName)
		return fmt.Errorf("rename to latest.png: %w", err)
	}
	return nil
}
