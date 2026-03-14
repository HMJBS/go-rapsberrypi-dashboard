// Package config はダッシュボードの起動設定を扱います。
package config

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"
)

// DefaultConfigPath は既定の設定ファイルパスです。
const DefaultConfigPath = "/etc/dashboard-config.json"

// Values は起動時に利用する設定値一式です。
type Values struct {
	FBPath          string
	PhotosDir       string
	CacheDir        string
	Latitude        float64
	Longitude       float64
	Timezone        string
	PhotoInterval   time.Duration
	RescanInterval  time.Duration
	WeatherInterval time.Duration
	PreviewDir      string
	PreviewEvery    time.Duration
	ScreenWidth     int
	ScreenHeight    int
}

type fileValues struct {
	Latitude                   *float64 `json:"latitude"`
	Longitude                  *float64 `json:"longitude"`
	Timezone                   *string  `json:"timezone"`
	PhotosDir                  *string  `json:"photos_dir"`
	CacheDir                   *string  `json:"cache_dir"`
	PhotoIntervalSeconds       *int64   `json:"photo_interval_seconds"`
	PhotoRescanIntervalSeconds *int64   `json:"photo_rescan_interval_seconds"`
	WeatherIntervalMinutes     *int64   `json:"weather_interval_minutes"`
	PreviewDir                 *string  `json:"preview_dir"`
	PreviewEveryMS             *int64   `json:"preview_every_ms"`
	ScreenWidth                *int     `json:"screen_w"`
	ScreenHeight               *int     `json:"screen_h"`
}

// DefaultDataRoot は既定のデータディレクトリのベースパスを返します。
func DefaultDataRoot() string {
	base, err := os.UserCacheDir()
	if err != nil {
		base = "."
	}
	return filepath.Join(base, "dashboard")
}

// DefaultValues は既定値を返します。
func DefaultValues(dataRoot string) Values {
	return Values{
		FBPath:          "/dev/fb0",
		PhotosDir:       filepath.Join(dataRoot, "photos"),
		CacheDir:        filepath.Join(dataRoot, "cache"),
		Latitude:        35.681236,
		Longitude:       139.767125,
		Timezone:        "Asia/Tokyo",
		PhotoInterval:   1 * time.Minute,
		RescanInterval:  5 * time.Minute,
		WeatherInterval: 10 * time.Minute,
		PreviewEvery:    1 * time.Second,
		ScreenWidth:     1280,
		ScreenHeight:    1024,
	}
}

// ApplyJSONFile は JSON 設定ファイルを読み込んで既存設定へ反映します。
// required が false の場合、ファイル未存在は無視します。
func ApplyJSONFile(current Values, path string, required bool) (Values, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if !required && errors.Is(err, os.ErrNotExist) {
			return current, nil
		}
		return Values{}, fmt.Errorf("read config file %q: %w", path, err)
	}

	var cfg fileValues
	dec := json.NewDecoder(bytes.NewReader(data))
	dec.DisallowUnknownFields()
	if err := dec.Decode(&cfg); err != nil {
		return Values{}, fmt.Errorf("decode config file %q: %w", path, err)
	}
	// 設定ファイル末尾に余分な JSON データがないことを確認する。
	var extra json.RawMessage
	if err := dec.Decode(&extra); err != io.EOF {
		if err == nil {
			return Values{}, fmt.Errorf("decode config file %q: unexpected extra JSON values", path)
		}
		return Values{}, fmt.Errorf("decode config file %q: %w", path, err)
	}

	if cfg.Latitude != nil {
		current.Latitude = *cfg.Latitude
	}
	if cfg.Longitude != nil {
		current.Longitude = *cfg.Longitude
	}
	if cfg.Timezone != nil {
		current.Timezone = *cfg.Timezone
	}
	if cfg.PhotosDir != nil {
		current.PhotosDir = *cfg.PhotosDir
	}
	if cfg.CacheDir != nil {
		current.CacheDir = *cfg.CacheDir
	}
	if cfg.PhotoIntervalSeconds != nil {
		current.PhotoInterval = time.Duration(*cfg.PhotoIntervalSeconds) * time.Second
	}
	if cfg.PhotoRescanIntervalSeconds != nil {
		current.RescanInterval = time.Duration(*cfg.PhotoRescanIntervalSeconds) * time.Second
	}
	if cfg.WeatherIntervalMinutes != nil {
		current.WeatherInterval = time.Duration(*cfg.WeatherIntervalMinutes) * time.Minute
	}
	if cfg.PreviewDir != nil {
		current.PreviewDir = *cfg.PreviewDir
	}
	if cfg.PreviewEveryMS != nil {
		current.PreviewEvery = time.Duration(*cfg.PreviewEveryMS) * time.Millisecond
	}
	if cfg.ScreenWidth != nil {
		current.ScreenWidth = *cfg.ScreenWidth
	}
	if cfg.ScreenHeight != nil {
		current.ScreenHeight = *cfg.ScreenHeight
	}

	return current, nil
}

// ApplyFlagOverrides は明示的に指定されたフラグだけを既存設定へ反映します。
func ApplyFlagOverrides(current, flags Values, visited map[string]bool) Values {
	if visited["fb"] {
		current.FBPath = flags.FBPath
	}
	if visited["photos_dir"] {
		current.PhotosDir = flags.PhotosDir
	}
	if visited["cache_dir"] {
		current.CacheDir = flags.CacheDir
	}
	if visited["lat"] {
		current.Latitude = flags.Latitude
	}
	if visited["lon"] {
		current.Longitude = flags.Longitude
	}
	if visited["tz"] {
		current.Timezone = flags.Timezone
	}
	if visited["photo_interval"] {
		current.PhotoInterval = flags.PhotoInterval
	}
	if visited["photo_rescan_interval"] {
		current.RescanInterval = flags.RescanInterval
	}
	if visited["weather_interval"] {
		current.WeatherInterval = flags.WeatherInterval
	}
	if visited["preview_dir"] {
		current.PreviewDir = flags.PreviewDir
	}
	if visited["preview_every"] {
		current.PreviewEvery = flags.PreviewEvery
	}
	if visited["screen_w"] {
		current.ScreenWidth = flags.ScreenWidth
	}
	if visited["screen_h"] {
		current.ScreenHeight = flags.ScreenHeight
	}
	return current
}

// Validate は設定値の妥当性を検証します。
func Validate(cfg Values) error {
	if cfg.PhotosDir == "" {
		return errors.New("photos_dir must not be empty")
	}
	if cfg.CacheDir == "" {
		return errors.New("cache_dir must not be empty")
	}
	if cfg.Latitude < -90 || cfg.Latitude > 90 {
		return fmt.Errorf("latitude must be between -90 and 90: %v", cfg.Latitude)
	}
	if cfg.Longitude < -180 || cfg.Longitude > 180 {
		return fmt.Errorf("longitude must be between -180 and 180: %v", cfg.Longitude)
	}
	if cfg.PhotoInterval < 0 {
		return fmt.Errorf("photo_interval must be >= 0: %s", cfg.PhotoInterval)
	}
	if cfg.RescanInterval < 0 {
		return fmt.Errorf("photo_rescan_interval must be >= 0: %s", cfg.RescanInterval)
	}
	if cfg.WeatherInterval < 0 {
		return fmt.Errorf("weather_interval must be >= 0: %s", cfg.WeatherInterval)
	}
	if cfg.PreviewEvery < 0 {
		return fmt.Errorf("preview_every must be >= 0: %s", cfg.PreviewEvery)
	}
	if cfg.ScreenWidth <= 0 {
		return fmt.Errorf("screen_w must be > 0: %d", cfg.ScreenWidth)
	}
	if cfg.ScreenHeight <= 0 {
		return fmt.Errorf("screen_h must be > 0: %d", cfg.ScreenHeight)
	}
	return nil
}
