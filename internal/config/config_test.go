package config

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestApplyJSONFileAndFlagOverrides(t *testing.T) {
	defaults := DefaultValues(filepath.Join(t.TempDir(), "dashboard"))
	configPath := filepath.Join(t.TempDir(), "dashboard-config.json")
	if err := os.WriteFile(configPath, []byte(`{
		"latitude": 40.7128,
		"longitude": -74.0060,
		"timezone": "America/New_York",
		"photos_dir": "/srv/dashboard/photos",
		"cache_dir": "/srv/dashboard/cache",
		"photo_interval_seconds": 120,
		"photo_rescan_interval_seconds": 900,
		"weather_interval_minutes": 30,
		"preview_dir": "/tmp/dashboard-preview",
		"preview_every_ms": 5000,
		"screen_w": 1920,
		"screen_h": 1080
	}`), 0o644); err != nil {
		t.Fatalf("write config file: %v", err)
	}

	resolved, err := ApplyJSONFile(defaults, configPath, true)
	if err != nil {
		t.Fatalf("ApplyJSONFile returned error: %v", err)
	}

	flags := resolved
	flags.Timezone = "Asia/Tokyo"
	flags.PhotoInterval = 45 * time.Second
	flags.ScreenWidth = 1280
	visited := map[string]bool{
		"tz":             true,
		"photo_interval": true,
		"screen_w":       true,
	}

	resolved = ApplyFlagOverrides(resolved, flags, visited)

	if resolved.Timezone != "Asia/Tokyo" {
		t.Fatalf("Timezone = %q, want Asia/Tokyo", resolved.Timezone)
	}
	if resolved.PhotoInterval != 45*time.Second {
		t.Fatalf("PhotoInterval = %s, want 45s", resolved.PhotoInterval)
	}
	if resolved.ScreenWidth != 1280 {
		t.Fatalf("ScreenWidth = %d, want 1280", resolved.ScreenWidth)
	}
	if resolved.CacheDir != "/srv/dashboard/cache" {
		t.Fatalf("CacheDir = %q, want /srv/dashboard/cache", resolved.CacheDir)
	}
	if resolved.WeatherInterval != 30*time.Minute {
		t.Fatalf("WeatherInterval = %s, want 30m", resolved.WeatherInterval)
	}
}

func TestApplyJSONFileOptionalMissing(t *testing.T) {
	defaults := DefaultValues(filepath.Join(t.TempDir(), "dashboard"))
	resolved, err := ApplyJSONFile(defaults, filepath.Join(t.TempDir(), "missing.json"), false)
	if err != nil {
		t.Fatalf("ApplyJSONFile returned error: %v", err)
	}
	if resolved != defaults {
		t.Fatalf("resolved config changed for optional missing file")
	}
}

func TestValidate(t *testing.T) {
	valid := DefaultValues(filepath.Join(t.TempDir(), "dashboard"))
	if err := Validate(valid); err != nil {
		t.Fatalf("Validate(valid) returned error: %v", err)
	}

	invalid := valid
	invalid.Latitude = 91
	if err := Validate(invalid); err == nil {
		t.Fatal("Validate should fail for invalid latitude")
	}
}
