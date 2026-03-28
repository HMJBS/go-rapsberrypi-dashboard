// Package assets は、ダッシュボードで使用する画像などのアセットを提供します。
package assets

import (
	"bytes"
	_ "embed"
	"image"
	_ "image/png" // PNG デコードを有効にするためのインポート

	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
)

//go:embed overlay.png
var overlay []byte

// Overlay はダッシュボードのオーバーレイ画像です。
var Overlay = mustLoadOverlay()

func mustLoadOverlay() image.Image {
	img, _, err := image.Decode(bytes.NewReader(overlay))
	if err != nil {
		panic(err)
	}
	return img
}

// InterDisplayMedium は、Google Fonts の Inter Display Medium フォントの TTF データです。
//
//go:embed InterDisplay-Medium.ttf
var interDisplayMedium []byte

// ClockDateFont は、時計ウィジェットの日付描画に使用するフォントです。
var ClockDateFont = mustLoadClockDateFont()

func mustLoadClockDateFont() font.Face {
	tt, err := opentype.Parse(interDisplayMedium)
	if err != nil {
		panic(err)
	}
	face, err := opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    24,
		DPI:     72,
		Hinting: font.HintingNone,
	})
	if err != nil {
		panic(err)
	}
	return face
}

//go:embed Inter-Bold.ttf
var interBold []byte

// ClockTimeFont は、時計ウィジェットの時刻描画に使用するフォントです。
var ClockTimeFont = mustLoadInterBold()

func mustLoadInterBold() font.Face {
	tt, err := opentype.Parse(interBold)
	if err != nil {
		panic(err)
	}
	face, err := opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    66,
		DPI:     72,
		Hinting: font.HintingNone,
	})
	if err != nil {
		panic(err)
	}
	return face
}
