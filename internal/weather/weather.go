// Package weather は、外部 API から天気情報を取得し、ローカルにキャッシュします。
package weather

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"time"
)

// Weather は現在の天気（気温・天気コード）と取得時刻を保持します。
type Weather struct {
	TempC     float64   `json:"temp_c"`
	Code      int       `json:"code"`
	Observed  time.Time `json:"observed"`
	FetchedAt time.Time `json:"fetched_at"`
}

// CachePath は weather キャッシュの保存先パスを返します。
func CachePath(cacheDir string) string {
	return filepath.Join(cacheDir, "weather.json")
}

// LoadCache はキャッシュファイルを読み込みます。
func LoadCache(path string) (Weather, bool) {
	b, err := os.ReadFile(path)
	if err != nil {
		return Weather{}, false
	}
	var w Weather
	if err := json.Unmarshal(b, &w); err != nil {
		return Weather{}, false
	}
	if w.FetchedAt.IsZero() {
		return Weather{}, false
	}
	return w, true
}

// SaveCache はキャッシュファイルへ保存します。
func SaveCache(path string, w Weather) error {
	b, err := json.MarshalIndent(w, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal weather cache: %w", err)
	}
	dir := filepath.Dir(path)
	tmpFile, err := os.CreateTemp(dir, "weather-*.tmp")
	if err != nil {
		return fmt.Errorf("create temp weather cache: %w", err)
	}
	tmpName := tmpFile.Name()
	defer func() {
		// Clean up the temp file on error paths.
		if tmpName != "" {
			_ = os.Remove(tmpName)
		}
		_ = tmpFile.Close()
	}()

	if _, err := tmpFile.Write(b); err != nil {
		return fmt.Errorf("write temp weather cache: %w", err)
	}
	if err := tmpFile.Sync(); err != nil {
		return fmt.Errorf("sync temp weather cache: %w", err)
	}
	if err := tmpFile.Close(); err != nil {
		return fmt.Errorf("close temp weather cache: %w", err)
	}

	// Prevent deferred removal after successful rename.
	name := tmpName
	tmpName = ""
	if err := os.Rename(name, path); err != nil {
		return fmt.Errorf("rename temp weather cache: %w", err)
	}
	if err := os.Chmod(path, 0o644); err != nil {
		return fmt.Errorf("chmod weather cache: %w", err)
	}
	return nil
}

// Client は天気 API クライアントです。
type Client struct {
	HTTPClient *http.Client
}

// Fetch は Open-Meteo から現在の気温と天気コードを取得します。
func (c Client) Fetch(ctx context.Context, lat, lon float64, tz string) (Weather, error) {
	hc := c.HTTPClient
	if hc == nil {
		hc = &http.Client{Timeout: 10 * time.Second}
	}

	q := url.Values{}
	q.Set("latitude", fmt.Sprintf("%.6f", lat))
	q.Set("longitude", fmt.Sprintf("%.6f", lon))
	q.Set("current", "temperature_2m,weather_code")
	q.Set("timezone", tz)

	u := url.URL{Scheme: "https", Host: "api.open-meteo.com", Path: "/v1/forecast", RawQuery: q.Encode()}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return Weather{}, fmt.Errorf("new request: %w", err)
	}

	resp, err := hc.Do(req)
	if err != nil {
		return Weather{}, fmt.Errorf("request: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode/100 != 2 {
		return Weather{}, fmt.Errorf("open-meteo status=%s", resp.Status)
	}

	var decoded struct {
		Current struct {
			Temperature2m float64 `json:"temperature_2m"`
			WeatherCode   int     `json:"weather_code"`
			Time          string  `json:"time"`
		} `json:"current"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&decoded); err != nil {
		return Weather{}, fmt.Errorf("decode: %w", err)
	}

	observed := time.Now()
	if decoded.Current.Time != "" {
		if t, err := time.Parse(time.RFC3339, decoded.Current.Time); err == nil {
			observed = t
		} else if t, err := time.Parse("2006-01-02T15:04", decoded.Current.Time); err == nil {
			observed = t
		}
	}

	return Weather{
		TempC:     decoded.Current.Temperature2m,
		Code:      decoded.Current.WeatherCode,
		Observed:  observed,
		FetchedAt: time.Now(),
	}, nil
}

// CodeLabel は Open-Meteo の weather_code を簡易ラベルへ変換します。
func CodeLabel(code int) string {
	switch code {
	case 0:
		return "CLEAR"
	case 1, 2, 3:
		return "CLOUDS"
	case 45, 48:
		return "FOG"
	case 51, 53, 55, 56, 57:
		return "DRIZZLE"
	case 61, 63, 65, 66, 67:
		return "RAIN"
	case 71, 73, 75, 77:
		return "SNOW"
	case 80, 81, 82:
		return "SHOWERS"
	case 85, 86:
		return "SNOWSHWR"
	case 95, 96, 99:
		return "THUNDER"
	default:
		return fmt.Sprintf("CODE%d", code)
	}
}
