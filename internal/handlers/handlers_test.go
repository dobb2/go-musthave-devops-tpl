package handlers

import (
	"github.com/dobb2/go-musthave-devops-tpl/internal/storage/metrics/cache"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestMetricsHandler_Gauge(t *testing.T) {
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
			name:   "positive test #1",
			url:    "/update/gauge/HeapInuse/933888.43",
			method: "POST",
			want: want{
				code: http.StatusOK,
			},
		},
		{
			name:   "negative test #2",
			url:    "/update/gauge/HeapInuse/933888.43",
			method: "GET",
			want: want{
				code: http.StatusMethodNotAllowed,
			},
		},
		{
			name:   "negative test #3",
			url:    "/update/gauge/HeapInuse/933888fdfd",
			method: "POST",
			want: want{
				code: http.StatusBadRequest,
			},
		},
		{
			name:   "negative test #4",
			url:    "/update/gauge/HeapInuse/",
			method: "POST",
			want: want{
				code: http.StatusNotFound,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(tt.method, tt.url, nil)
			w := httptest.NewRecorder()

			a := New(cache.Create())
			a.Gauge(w, req)

			result := w.Result()
			assert.Equal(t, tt.want.code, result.StatusCode)

			_, err := ioutil.ReadAll(result.Body)
			require.NoError(t, err)
			err = result.Body.Close()
			require.NoError(t, err)
		})
	}
}

func TestMetricsHandler_Counter(t *testing.T) {
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
			name:   "positive test #1",
			url:    "/update/counter/PollCount/13",
			method: "POST",
			want: want{
				code: http.StatusOK,
			},
		},
		{
			name:   "negative test #2",
			url:    "/update/counter/PollCount/13",
			method: "GET",
			want: want{
				code: http.StatusMethodNotAllowed,
			},
		},
		{
			name:   "negative test #3",
			url:    "/update/counter/PollCount/13cd",
			method: "POST",
			want: want{
				code: http.StatusBadRequest,
			},
		},
		{
			name:   "negative test #4",
			url:    "/update/counter/PollCount/13.33",
			method: "POST",
			want: want{
				code: http.StatusBadRequest,
			},
		},
		{
			name:   "negative test #5",
			url:    "/update/counter/PollCount/",
			method: "POST",
			want: want{
				code: http.StatusNotFound,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(tt.method, tt.url, nil)
			w := httptest.NewRecorder()

			a := New(cache.Create())
			a.Counter(w, req)

			result := w.Result()
			assert.Equal(t, tt.want.code, result.StatusCode)

			_, err := ioutil.ReadAll(result.Body)
			require.NoError(t, err)
			err = result.Body.Close()
			require.NoError(t, err)
		})
	}
}
