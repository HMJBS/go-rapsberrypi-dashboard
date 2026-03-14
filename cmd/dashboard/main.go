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
	"syscall"

	"dashboard/internal/app"
	"dashboard/internal/config"
)

func main() {
	flagValues, debug, configPath, configPathExplicit, visited := parseFlags()

	logger := log.New(os.Stderr, "dashboard: ", log.LstdFlags)
	if debug {
		logger.SetFlags(log.LstdFlags | log.Lmicroseconds)
	}

	resolved := config.DefaultValues(config.DefaultDataRoot())
	var err error
	resolved, err = config.ApplyJSONFile(resolved, configPath, configPathExplicit)
	if err != nil {
		logger.Fatalf("load config: %v", err)
	}
	resolved = config.ApplyFlagOverrides(resolved, flagValues, visited)
	if err := config.Validate(resolved); err != nil {
		logger.Fatalf("invalid config: %v", err)
	}

	if err := os.MkdirAll(resolved.PhotosDir, 0o755); err != nil {
		logger.Fatalf("create photos_dir: %v", err)
	}
	if err := os.MkdirAll(resolved.CacheDir, 0o755); err != nil {
		logger.Fatalf("create cache_dir: %v", err)
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	cfg := app.Config{
		FBPath:          resolved.FBPath,
		PhotosDir:       resolved.PhotosDir,
		CacheDir:        resolved.CacheDir,
		PreviewDir:      resolved.PreviewDir,
		PreviewEvery:    resolved.PreviewEvery,
		ScreenWidth:     resolved.ScreenWidth,
		ScreenHeight:    resolved.ScreenHeight,
		Latitude:        resolved.Latitude,
		Longitude:       resolved.Longitude,
		Timezone:        resolved.Timezone,
		PhotoInterval:   resolved.PhotoInterval,
		RescanInterval:  resolved.RescanInterval,
		WeatherInterval: resolved.WeatherInterval,
		Background:      color.RGBA{R: 0, G: 0, B: 0, A: 255},
	}

	if err := app.Run(ctx, logger, cfg); err != nil {
		if errors.Is(err, context.Canceled) {
			return
		}
		// If running on dev machine without framebuffer, provide a helpful hint.
		var pathErr *os.PathError
		if errors.As(err, &pathErr) && pathErr.Path == resolved.FBPath {
			fmt.Fprintln(os.Stderr, "Hint: This program needs Linux framebuffer. On dev machines, you can use -preview_dir to write latest.png without /dev/fb0.")
		}
		logger.Fatalf("fatal: %v", err)
	}
}

// parseFlags はコマンドラインフラグを解釈し、設定値と補助情報を返します。
// 戻り値は順に、フラグ値、debug 指定の有無、設定ファイルパス、config フラグの明示指定有無、
// そして JSON 設定より優先して上書きすべきフラグ名の集合です。
func parseFlags() (config.Values, bool, string, bool, map[string]bool) {
	defaults := config.DefaultValues(config.DefaultDataRoot())
	configPath := config.DefaultConfigPath
	debug := false

	flag.StringVar(&defaults.FBPath, "fb", defaults.FBPath, "framebuffer device path")
	flag.StringVar(&configPath, "config", configPath, "path to JSON config file")
	flag.StringVar(&defaults.PhotosDir, "photos_dir", defaults.PhotosDir, "directory containing synced photos")
	flag.StringVar(&defaults.CacheDir, "cache_dir", defaults.CacheDir, "directory for cached data (weather)")
	flag.Float64Var(&defaults.Latitude, "lat", defaults.Latitude, "latitude for weather")
	flag.Float64Var(&defaults.Longitude, "lon", defaults.Longitude, "longitude for weather")
	flag.StringVar(&defaults.Timezone, "tz", defaults.Timezone, "IANA timezone name")
	flag.DurationVar(&defaults.PhotoInterval, "photo_interval", defaults.PhotoInterval, "interval for changing photo")
	flag.DurationVar(&defaults.RescanInterval, "photo_rescan_interval", defaults.RescanInterval, "interval for rescanning photos_dir")
	flag.DurationVar(&defaults.WeatherInterval, "weather_interval", defaults.WeatherInterval, "interval for updating weather")
	flag.BoolVar(&debug, "debug", debug, "enable debug logs")
	flag.StringVar(&defaults.PreviewDir, "preview_dir", defaults.PreviewDir, "write latest frame as PNG into this directory (dev/debug mode; disables framebuffer)")
	flag.DurationVar(&defaults.PreviewEvery, "preview_every", defaults.PreviewEvery, "interval for updating latest.png when preview_dir is set")
	flag.IntVar(&defaults.ScreenWidth, "screen_w", defaults.ScreenWidth, "screen width for preview mode")
	flag.IntVar(&defaults.ScreenHeight, "screen_h", defaults.ScreenHeight, "screen height for preview mode")
	flag.Parse()

	visited := map[string]bool{}
	flag.CommandLine.Visit(func(f *flag.Flag) {
		visited[f.Name] = true
	})

	return defaults, debug, configPath, visited["config"], visited
}
