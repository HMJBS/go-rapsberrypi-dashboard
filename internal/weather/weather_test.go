package weather

import (
	"context"
	"io"
	"net/http"
	"strings"
	"testing"
	"time"
)

func TestClientFetchIncludesDailyAndPopulatesWeather(t *testing.T) {

	var capturedQuery map[string]string
	client := Client{HTTPClient: &http.Client{Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
		capturedQuery = map[string]string{
			"latitude":  req.URL.Query().Get("latitude"),
			"longitude": req.URL.Query().Get("longitude"),
			"current":   req.URL.Query().Get("current"),
			"daily":     req.URL.Query().Get("daily"),
			"timezone":  req.URL.Query().Get("timezone"),
		}
		return jsonResponse(`{
			"current": {
				"temperature_2m": 21.5,
				"weather_code": 3,
				"time": "2026-03-29T11:45"
			},
			"daily": {
				"sunrise": ["2026-03-29T05:32"],
				"sunset": ["2026-03-29T18:01"]
			}
		}`), nil
	})}}

	w, err := client.Fetch(context.Background(), 35.681236, 139.767125, "Asia/Tokyo")
	if err != nil {
		t.Fatalf("Fetch returned error: %v", err)
	}

	if capturedQuery["current"] != "temperature_2m,weather_code" {
		t.Fatalf("current query = %q", capturedQuery["current"])
	}
	if capturedQuery["timezone"] != "Asia/Tokyo" {
		t.Fatalf("timezone query = %q", capturedQuery["timezone"])
	}
	if w.TempC != 21.5 {
		t.Fatalf("TempC = %v, want 21.5", w.TempC)
	}
	if w.Code != 3 {
		t.Fatalf("Code = %d, want 3", w.Code)
	}
}

func TestFetchDailyWeather(t *testing.T) {
	loc := time.FixedZone("JST", 9*60*60)
	client := Client{HTTPClient: &http.Client{Transport: roundTripFunc(func(_ *http.Request) (*http.Response, error) {
		return jsonResponse(`{
			"daily": {
				"sunrise": ["2026-03-29T05:30"],
				"sunset": ["2026-03-29T18:00"]
			}
		}`), nil
	})}}

	daily, err := fetchDailyWeather(context.Background(), client, 35.681236, 139.767125, "Asia/Tokyo")
	if err != nil {
		t.Fatalf("fetchDailyWeather returned error: %v", err)
	}

	if got := daily.Sunrise.In(loc).Format("2006-01-02T15:04"); got != "2026-03-29T05:30" {
		t.Fatalf("Sunrise = %s", got)
	}
	if got := daily.Sunset.In(loc).Format("2006-01-02T15:04"); got != "2026-03-29T18:00" {
		t.Fatalf("Sunset = %s", got)
	}
}

type roundTripFunc func(*http.Request) (*http.Response, error)

func (fn roundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return fn(req)
}

func jsonResponse(body string) *http.Response {
	return &http.Response{
		StatusCode: http.StatusOK,
		Header:     make(http.Header),
		Body:       io.NopCloser(strings.NewReader(body)),
	}
}
