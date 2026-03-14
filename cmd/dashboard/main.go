// Package main provides the dashboard executable.
package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"image/color"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"dashboard/internal/app"
)

func main() {
	var (
		fbPath          = flag.String("fb", "/dev/fb0", "framebuffer device path")
		photosDir       = flag.String("photos_dir", "", "directory containing synced photos")
		cacheDir        = flag.String("cache_dir", "", "directory for cached data (weather)")
		latitude        = flag.Float64("lat", 35.681236, "latitude for weather")
		longitude       = flag.Float64("lon", 139.767125, "longitude for weather")
		timezone        = flag.String("tz", "Asia/Tokyo", "IANA timezone name")
		photoInterval   = flag.Duration("photo_interval", 1*time.Minute, "interval for changing photo")
		rescanInterval  = flag.Duration("photo_rescan_interval", 5*time.Minute, "interval for rescanning photos_dir")
		weatherInterval = flag.Duration("weather_interval", 10*time.Minute, "interval for updating weather")
		debug           = flag.Bool("debug", false, "enable debug logs")
		previewDir      = flag.String("preview_dir", "", "write latest frame as PNG into this directory (dev/debug mode; disables framebuffer)")
		previewEvery    = flag.Duration("preview_every", 1*time.Second, "interval for updating latest.png when preview_dir is set")
		screenW         = flag.Int("screen_w", 1280, "screen width for preview mode")
		screenH         = flag.Int("screen_h", 1024, "screen height for preview mode")
	)
	flag.Parse()

	logger := log.New(os.Stderr, "dashboard: ", log.LstdFlags|log.Lmsgprefix)
	if *debug {
		logger.SetFlags(log.LstdFlags | log.Lmicroseconds | log.Lmsgprefix)
	}

	if *photosDir == "" {
		// Default to user cache dir if not specified.
		base, err := os.UserCacheDir()
		if err != nil {
			base = "."
		}
		*photosDir = filepath.Join(base, "dashboard", "photos")
	}
	if *cacheDir == "" {
		base, err := os.UserCacheDir()
		if err != nil {
			base = "."
		}
		*cacheDir = filepath.Join(base, "dashboard", "cache")
	}

	if err := os.MkdirAll(*photosDir, 0o755); err != nil {
		logger.Fatalf("create photos_dir: %v", err)
	}
	if err := os.MkdirAll(*cacheDir, 0o755); err != nil {
		logger.Fatalf("create cache_dir: %v", err)
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	cfg := app.Config{
		FBPath:          *fbPath,
		PhotosDir:       *photosDir,
		CacheDir:        *cacheDir,
		PreviewDir:      *previewDir,
		PreviewEvery:    *previewEvery,
		ScreenWidth:     *screenW,
		ScreenHeight:    *screenH,
		Latitude:        *latitude,
		Longitude:       *longitude,
		Timezone:        *timezone,
		PhotoInterval:   *photoInterval,
		RescanInterval:  *rescanInterval,
		WeatherInterval: *weatherInterval,
		Background:      color.RGBA{R: 0, G: 0, B: 0, A: 255},
	}

	if err := app.Run(ctx, logger, cfg); err != nil {
		if errors.Is(err, context.Canceled) {
			return
		}
		// If running on dev machine without framebuffer, provide a helpful hint.
		var pathErr *os.PathError
		if errors.As(err, &pathErr) && pathErr.Path == *fbPath {
			fmt.Fprintln(os.Stderr, "Hint: This program needs Linux framebuffer. On dev machines, you can use -preview_dir to write latest.png without /dev/fb0.")
		}
		logger.Fatalf("fatal: %v", err)
	}
}
