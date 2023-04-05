package handlers

import (
	"encoding/json"
	"github.com/go-chi/chi/v5/middleware"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/dobb2/go-musthave-devops-tpl/internal/storage/metrics"
	"github.com/dobb2/go-musthave-devops-tpl/internal/storage/metrics/cache"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func testRequest(t *testing.T, ts *httptest.Server, path, method string, body io.Reader) (int, string) {
	req, err := http.NewRequest(method, ts.URL+path, body)
	require.NoError(t, err)

	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)

	respBody, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	defer resp.Body.Close()
	return resp.StatusCode, string(respBody)
}

func TestMetricsHandler_PostUpdateMetric(t *testing.T) {
	type want struct {
		code int
	}
	tests := []struct {
		name   string
		url    string
		json   string
		method string
		want   want
	}{
		{
			name:   "positive update test #1",
			url:    "/update/",
			json:   `{"id":"HeapInuse","type":"gauge","value":933888.43, "hash":""}`,
			method: "POST",
			want: want{
				code: http.StatusOK,
			},
		},
		{
			name:   "negative update test #2",
			url:    "/update/",
			json:   `{"id":"HeapInuse","type":"gauge","value":933888.43, "hash":""}`,
			method: "GET",
			want: want{
				code: http.StatusMethodNotAllowed,
			},
		},
		{
			name:   "negative update test #3",
			url:    "/update/",
			json:   `{"id":"HeapInuse","type":"gauge","value":933888fdfd, "hash":""}`,
			method: "POST",
			want: want{
				code: http.StatusBadRequest,
			},
		},
		{
			name:   "negative update test #4",
			url:    "/update/",
			json:   `{"id":"HeapInuse","type":"gauge"}`,
			method: "POST",
			want: want{
				code: http.StatusBadRequest,
			},
		},
		{
			name:   "positive update test #5",
			url:    "/update/",
			json:   `{"id":"PollCount","type":"counter","delta":13, "hash":""}`,
			method: "POST",
			want: want{
				code: http.StatusOK,
			},
		},
		{
			name:   "negative update test #6",
			url:    "/update/",
			json:   `{"id":"PollCount","type":"counter","delta":13, "hash":""}`,
			method: "GET",
			want: want{
				code: http.StatusMethodNotAllowed,
			},
		},
		{
			name:   "negative update test #7",
			url:    "/update/",
			json:   `{"id":"PollCount","type":"counter","delta":13cd, "hash":""}`,
			method: "POST",
			want: want{
				code: http.StatusBadRequest,
			},
		},
		{
			name:   "negative update test #8",
			url:    "/update/",
			json:   `{"id":"PollCount","type":"counter","delta":13.33, "hash":""}`,
			method: "POST",
			want: want{
				code: http.StatusBadRequest,
			},
		},
		{
			name:   "negative update test #9",
			url:    "/update/",
			json:   `{"id":"PollCount","type":"counter"}`,
			method: "POST",
			want: want{
				code: http.StatusBadRequest,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := New(cache.Create())

			r := func(m MetricsHandler) chi.Router {
				r := chi.NewRouter()
				r.Use(middleware.WithValue("Key", ""))
				r.Post("/update/", m.PostUpdateMetric)
				return r
			}(a)
			ts := httptest.NewServer(r)
			defer ts.Close()

			statusCode, _ := testRequest(t, ts, tt.url, tt.method, strings.NewReader(tt.json))
			assert.Equal(t, tt.want.code, statusCode)
		})
	}
}

func TestMetricsHandler_GetAllMetrics(t *testing.T) {
	type want struct {
		code int
	}
	tests := []struct {
		name   string
		url    string
		json   string
		method string
		want   want
	}{
		{
			name:   "negative get all metric test #1",
			url:    "/",
			json:   `{}`,
			method: "GET",
			want: want{
				code: http.StatusBadRequest,
			},
		},
		{
			name:   "negative get all metric test #2",
			url:    "/",
			json:   `{}`,
			method: "POST",
			want: want{
				code: http.StatusMethodNotAllowed,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := New(cache.Create())
			var metric metrics.Metrics

			if err := json.NewDecoder(strings.NewReader(tt.json)).Decode(&metric); err == nil {
				switch metric.MType {
				case "gauge":
					a.storage.UpdateGauge(metric.ID, *metric.Value)
				case "counter":
					a.storage.UpdateCounter(metric.ID, *metric.Delta)
				}
			}

			r := func(m MetricsHandler) chi.Router {
				r := chi.NewRouter()
				r.Get("/", a.GetAllMetrics)
				return r
			}(a)
			ts := httptest.NewServer(r)
			defer ts.Close()

			statusCode, _ := testRequest(t, ts, tt.url, tt.method, nil)
			assert.Equal(t, tt.want.code, statusCode)
		})
	}
}

func TestMetricsHandler_PostGetMetric(t *testing.T) {
	type want struct {
		code int
	}
	tests := []struct {
		name   string
		url    string
		json   string
		added  bool
		method string
		want   want
	}{
		{
			name:   "positive get metric test #1",
			url:    "/value/",
			added:  true,
			json:   `{"id":"Testmetricid1","type":"gauge","value":434.32}`,
			method: "POST",
			want: want{
				code: http.StatusOK,
			},
		},
		{
			name:   "negative get all metric test #2",
			url:    "/value/",
			added:  true,
			json:   `{}`,
			method: "POST",
			want: want{
				code: http.StatusNotFound,
			},
		},
		{
			name:   "positive get metric test #3",
			added:  true,
			url:    "/value/",
			json:   `{"id":"Testmetric","type":"counter","delta":434}`,
			method: "POST",
			want: want{
				code: http.StatusOK,
			},
		},
		{
			name:   "negative metric test #4",
			url:    "/value/",
			added:  true,
			json:   `{"id":"Testmetricid2","type":"counter","delta":434}`,
			method: "GET",
			want: want{
				code: http.StatusMethodNotAllowed,
			},
		},
		{
			name:   "negative metric test #5",
			added:  false,
			url:    "/value/",
			json:   `{"id":"Testmetricid","type":"counter","delta":434}`,
			method: "POST",
			want: want{
				code: http.StatusNotFound,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := New(cache.Create())
			var metric metrics.Metrics

			if err := json.NewDecoder(strings.NewReader(tt.json)).Decode(&metric); err == nil && tt.added {
				switch metric.MType {
				case "gauge":
					a.storage.UpdateGauge(metric.ID, *metric.Value)
				case "counter":
					a.storage.UpdateCounter(metric.ID, *metric.Delta)
				}
			}

			r := func(m MetricsHandler) chi.Router {
				r := chi.NewRouter()
				r.Use(middleware.WithValue("Key", ""))
				r.Post("/value/", a.PostGetMetric)
				return r
			}(a)
			ts := httptest.NewServer(r)
			defer ts.Close()

			statusCode, _ := testRequest(t, ts, tt.url, tt.method, strings.NewReader(tt.json))
			assert.Equal(t, tt.want.code, statusCode)
		})
	}
}

func TestMetricsHandler_UpdateMetric(t *testing.T) {
	type want struct {
		code int
	}
	tests := []struct {
		name   string
		url    string
		method string
		want   want
	}{
		{
			name:   "positive gauge test #1",
			url:    "/update/gauge/HeapInuse/933888.43",
			method: "POST",
			want: want{
				code: http.StatusOK,
			},
		},
		{
			name:   "negative gauge test #2",
			url:    "/update/gauge/HeapInuse/933888.43",
			method: "GET",
			want: want{
				code: http.StatusMethodNotAllowed,
			},
		},
		{
			name:   "negative gauge test #3",
			url:    "/update/gauge/HeapInuse/933888fdfd",
			method: "POST",
			want: want{
				code: http.StatusBadRequest,
			},
		},
		{
			name:   "negative gauge test #4",
			url:    "/update/gauge/HeapInuse/",
			method: "POST",
			want: want{
				code: http.StatusNotFound,
			},
		},
		{
			name:   "positive counter test #1",
			url:    "/update/counter/PollCount/13",
			method: "POST",
			want: want{
				code: http.StatusOK,
			},
		},
		{
			name:   "negative counter test #2",
			url:    "/update/counter/PollCount/13",
			method: "GET",
			want: want{
				code: http.StatusMethodNotAllowed,
			},
		},
		{
			name:   "negative counter test #3",
			url:    "/update/counter/PollCount/13cd",
			method: "POST",
			want: want{
				code: http.StatusBadRequest,
			},
		},
		{
			name:   "negative counter test #4",
			url:    "/update/counter/PollCount/13.33",
			method: "POST",
			want: want{
				code: http.StatusBadRequest,
			},
		},
		{
			name:   "negative counter test #5",
			url:    "/update/counter/PollCount/",
			method: "POST",
			want: want{
				code: http.StatusNotFound,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := New(cache.Create())
			r := func(m MetricsHandler) chi.Router {
				r := chi.NewRouter()
				r.Post("/update/{typeMetric}/{nameMetric}/{value}", m.UpdateMetric)
				return r
			}(a)
			ts := httptest.NewServer(r)
			defer ts.Close()

			statusCode, _ := testRequest(t, ts, tt.url, tt.method, nil)
			assert.Equal(t, tt.want.code, statusCode)
		})
	}
}
