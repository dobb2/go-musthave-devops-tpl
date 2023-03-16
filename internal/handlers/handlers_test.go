package handlers

import (
	"github.com/dobb2/go-musthave-devops-tpl/internal/storage/metrics"
	"github.com/dobb2/go-musthave-devops-tpl/internal/storage/metrics/cache"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
)

func testRequest(t *testing.T, ts *httptest.Server, method, path string) (int, string) {
	req, err := http.NewRequest(method, ts.URL+path, nil)
	require.NoError(t, err)

	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)

	respBody, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	defer resp.Body.Close()

	return resp.StatusCode, string(respBody)
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

			statusCode, _ := testRequest(t, ts, tt.method, tt.url)
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
		metric metrics.Metric
		method string
		want   want
	}{
		{
			name:   "negative get all metric test #1",
			url:    "/",
			metric: metrics.Metric{},
			method: "GET",
			want: want{
				code: http.StatusBadRequest,
			},
		},
		{
			name:   "negative get all metric test #2",
			url:    "/",
			metric: metrics.Metric{},
			method: "POST",
			want: want{
				code: http.StatusMethodNotAllowed,
			},
		},
		{
			name: "positive get all metric test #3",
			url:  "/",
			metric: metrics.Metric{
				TypeMetric: "gauge",
				NameMetric: "Testmetricid1",
				Value:      "434.32",
			},
			method: "GET",
			want: want{
				code: http.StatusOK,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := New(cache.Create())

			switch tt.metric.TypeMetric {
			case "gauge":
				value, _ := strconv.ParseFloat(tt.metric.Value, 64)
				a.storage.UpdateGauge(tt.metric.NameMetric, value)
			case "counter":
				value, _ := strconv.ParseInt(tt.metric.Value, 10, 64)
				a.storage.UpdateCounter(tt.metric.NameMetric, value)
			}

			r := func(m MetricsHandler) chi.Router {
				r := chi.NewRouter()
				r.Get("/", a.GetAllMetrics)
				return r
			}(a)
			ts := httptest.NewServer(r)
			defer ts.Close()

			statusCode, _ := testRequest(t, ts, tt.method, tt.url)
			assert.Equal(t, tt.want.code, statusCode)
		})
	}
}

func TestMetricsHandler_GetMetric(t *testing.T) {
	type want struct {
		code int
	}
	tests := []struct {
		name   string
		url    string
		metric metrics.Metric
		method string
		want   want
	}{
		{
			name: "positive get metric test #1",
			url:  "/value/gauge/Testmetricid1",
			metric: metrics.Metric{
				TypeMetric: "gauge",
				NameMetric: "Testmetricid1",
				Value:      "434.32",
			},
			method: "GET",
			want: want{
				code: http.StatusOK,
			},
		},
		{
			name:   "negative get all metric test #2",
			url:    "/value/counter/Testmetric",
			metric: metrics.Metric{},
			method: "GET",
			want: want{
				code: http.StatusNotFound,
			},
		},
		{
			name:   "positive get all metric test #3",
			url:    "/value/counter/Testmetric",
			method: "GET",
			metric: metrics.Metric{
				TypeMetric: "counter",
				NameMetric: "Testmetric",
				Value:      "434",
			},
			want: want{
				code: http.StatusOK,
			},
		},
		{
			name:   "negative metric test #4",
			url:    "/value/counter/Testmetricid2",
			method: "GET",
			metric: metrics.Metric{
				TypeMetric: "counter",
				NameMetric: "Testmetric",
				Value:      "434",
			},
			want: want{
				code: http.StatusNotFound,
			},
		},
		{
			name:   "negative metric test #5",
			url:    "/value/counter/Testmetricid2",
			method: "POST",
			metric: metrics.Metric{
				TypeMetric: "counter",
				NameMetric: "Testmetric",
				Value:      "434",
			},
			want: want{
				code: http.StatusMethodNotAllowed,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := New(cache.Create())

			switch tt.metric.TypeMetric {
			case "gauge":
				value, _ := strconv.ParseFloat(tt.metric.Value, 64)
				a.storage.UpdateGauge(tt.metric.NameMetric, value)
			case "counter":
				value, _ := strconv.ParseInt(tt.metric.Value, 10, 64)
				a.storage.UpdateCounter(tt.metric.NameMetric, value)
			}

			r := func(m MetricsHandler) chi.Router {
				r := chi.NewRouter()
				r.Get("/value/{typeMetric}/{nameMetric}", a.GetMetric)
				return r
			}(a)
			ts := httptest.NewServer(r)
			defer ts.Close()

			statusCode, _ := testRequest(t, ts, tt.method, tt.url)
			assert.Equal(t, tt.want.code, statusCode)
		})
	}
}
