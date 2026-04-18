// Package app は、描画ループとデータ取得（天気/写真）を統合して実行します。
package app

import (
	"context"
	"fmt"
	"image"
	"log"
	"math/rand"
	"sync"
	"time"

	"dashboard/internal/assets"
	"dashboard/internal/fb"
	"dashboard/internal/gfx"
	"dashboard/internal/photos"
	"dashboard/internal/previewpng"
	"dashboard/internal/theme"
	"dashboard/internal/weather"
	"dashboard/internal/widgets"
)

// Config はダッシュボードの実行設定です。
type Config struct {
	FBPath    string
	PhotosDir string
	CacheDir  string

	// PreviewDir が空でない場合、FrameBuffer を使わずに描画結果を PNG に出力します。
	// 開発環境向けの簡易デバッグ用途です。
	PreviewDir   string
	PreviewEvery time.Duration
	ScreenWidth  int
	ScreenHeight int

	Latitude  float64
	Longitude float64
	Timezone  string

	PhotoInterval   time.Duration
	RescanInterval  time.Duration
	WeatherInterval time.Duration
}

type appState struct {
	mu sync.RWMutex

	photo     *image.RGBA
	photoPath string
	photoAt   time.Time
	photoErr  string

	photoFiles []string
	lastScan   time.Time

	weather     weather.Weather
	weatherIcon image.Image
	weatherOK   bool
	weatherErr  string
}

// Run はダッシュボードアプリを起動し、ctx がキャンセルされるまで描画を継続します。
func Run(ctx context.Context, logger *log.Logger, cfg Config) error {
	var (
		fbdev *fb.Framebuffer
		sz    image.Point
	)

	if cfg.PreviewDir != "" {
		w := cfg.ScreenWidth
		h := cfg.ScreenHeight
		if w <= 0 {
			w = 1280
		}
		if h <= 0 {
			h = 1024
		}
		sz = image.Point{X: w, Y: h}
	} else {
		var err error
		fbdev, err = fb.Open(cfg.FBPath)
		if err != nil {
			return err
		}
		defer func() { _ = fbdev.Close() }()
		sz = fbdev.Size()
	}

	frame := image.NewRGBA(image.Rect(0, 0, sz.X, sz.Y))

	state := &appState{}
	loadInitialWeatherCache(logger, cfg, state)

	go weatherLoop(ctx, logger, cfg, state)

	rescanPhotos(logger, cfg, state)
	changePhoto(logger, state, sz)

	loc := time.Local
	if cfg.Timezone != "" {
		l, err := time.LoadLocation(cfg.Timezone)
		if err != nil {
			return fmt.Errorf("invalid timezone %q: %w", cfg.Timezone, err)
		}
		loc = l
	}

	nextPhoto := time.Now().In(loc).Add(cfg.PhotoInterval)
	nextScan := time.Now().In(loc).Add(cfg.RescanInterval)

	tick := time.NewTicker(1 * time.Second)
	defer tick.Stop()

	nextPreview := time.Now().In(loc)
	if cfg.PreviewEvery <= 0 {
		cfg.PreviewEvery = 1 * time.Second
	}

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-tick.C:
			now := time.Now().In(loc)
			if cfg.RescanInterval > 0 && now.After(nextScan) {
				rescanPhotos(logger, cfg, state)
				nextScan = now.Add(cfg.RescanInterval)
			}
			if cfg.PhotoInterval > 0 && now.After(nextPhoto) {
				changePhoto(logger, state, sz)
				nextPhoto = now.Add(cfg.PhotoInterval)
			}

			render(frame, now, state)
			if cfg.PreviewDir != "" {
				if now.After(nextPreview) {
					if err := previewpng.WriteLatestPNG(cfg.PreviewDir, frame); err != nil {
						return err
					}
					nextPreview = now.Add(cfg.PreviewEvery)
				}
				continue
			}

			if err := fbdev.BlitRGBA(frame); err != nil {
				return err
			}
		}
	}
}

func loadInitialWeatherCache(logger *log.Logger, cfg Config, state *appState) {
	p := weather.CachePath(cfg.CacheDir)
	w, ok := weather.LoadCache(p)
	wIcon, err := w.GetIcon(context.Background(), weather.Client{}, cfg.Latitude, cfg.Longitude, cfg.Timezone)
	if err != nil {
		logger.Printf("weather icon load failed: %v", err)
	}

	if ok && err == nil {
		state.mu.Lock()
		state.weather = w
		state.weatherIcon = wIcon
		state.weatherOK = true
		state.mu.Unlock()
		logger.Printf("loaded weather cache: temp=%.1f code=%d", w.TempC, w.Code)
	}
}

func weatherLoop(ctx context.Context, logger *log.Logger, cfg Config, state *appState) {
	if cfg.WeatherInterval <= 0 {
		return
	}

	client := weather.Client{}
	cachePath := weather.CachePath(cfg.CacheDir)

	updateWeather(ctx, logger, client, cfg, state, cachePath)

	ticker := time.NewTicker(cfg.WeatherInterval)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			updateWeather(ctx, logger, client, cfg, state, cachePath)
		}
	}
}

func updateWeather(ctx context.Context, logger *log.Logger, client weather.Client, cfg Config, state *appState, cachePath string) {
	w, err1 := client.Fetch(ctx, cfg.Latitude, cfg.Longitude, cfg.Timezone)
	icon, err2 := w.GetIcon(ctx, client, cfg.Latitude, cfg.Longitude, cfg.Timezone)
	state.mu.Lock()
	defer state.mu.Unlock()
	if err1 != nil {
		state.weatherOK = false
		state.weatherIcon = nil
		state.weatherErr = err1.Error()
		logger.Printf("weather fetch failed: %v", err1)
		return
	} else if err2 != nil {
		state.weatherOK = false
		state.weatherIcon = nil
		state.weatherErr = err2.Error()
		logger.Printf("weather icon fetch failed: %v", err2)
		return
	}
	state.weather = w
	state.weatherIcon = icon
	state.weatherOK = true
	state.weatherErr = ""
	if err := weather.SaveCache(cachePath, w); err != nil {
		logger.Printf("weather cache save failed: %v", err)
	}
}

func rescanPhotos(logger *log.Logger, cfg Config, state *appState) {
	files, err := photos.ListImages(cfg.PhotosDir)
	state.mu.Lock()
	defer state.mu.Unlock()
	state.lastScan = time.Now()
	if err != nil {
		state.photoErr = err.Error()
		logger.Printf("photo scan failed: %v", err)
		return
	}
	state.photoFiles = files
	if len(files) == 0 {
		state.photoErr = "no images"
		return
	}
	if state.photoErr == "no images" {
		state.photoErr = ""
	}
}

func changePhoto(logger *log.Logger, state *appState, sz image.Point) {
	state.mu.RLock()
	files := append([]string(nil), state.photoFiles...)
	prev := state.photoPath
	state.mu.RUnlock()

	if len(files) == 0 {
		return
	}

	pick := files[rand.Intn(len(files))]
	if len(files) > 1 {
		for i := 0; i < 3 && pick == prev; i++ {
			pick = files[rand.Intn(len(files))]
		}
	}

	img, err := photos.LoadScreenImage(pick, sz.X, sz.Y, theme.DefaultTheme.BackgroundColor)
	state.mu.Lock()
	defer state.mu.Unlock()
	state.photoAt = time.Now()
	state.photoPath = pick
	if err != nil {
		state.photoErr = err.Error()
		logger.Printf("photo load failed (%s): %v", pick, err)
		return
	}
	state.photo = img
	state.photoErr = ""
	logger.Printf("photo changed: %s", pick)
}

func render(dst *image.RGBA, now time.Time, state *appState) {
	gfx.FillRGBA(dst, theme.DefaultTheme.BackgroundColor)

	state.mu.RLock()
	// photo := state.photo
	w := state.weather
	wIcon := state.weatherIcon
	wOK := state.weatherOK
	wErr := state.weatherErr
	// pErr := state.photoErr
	state.mu.RUnlock()

	// TODO: 画像ファイルを描画する

	// オーバーレイを描画
	gfx.DrawImage(dst, assets.Overlay, 0, 0)

	widgets.DrawWeatherForecastWidget(dst, w, wIcon, wOK, wErr)

	widgets.DrawClockWidget(dst, now)
	// TODO: 時刻を描画
}
