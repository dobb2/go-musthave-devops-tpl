package handlers

import (
	"github.com/dobb2/go-musthave-devops-tpl/internal/storage"
	"net/http"
	"strconv"
	"strings"
)

type MetricsHandler struct {
	storage storage.MetricCreatorUpdater
}

func New(metrics storage.MetricCreatorUpdater) MetricsHandler {
	return MetricsHandler{storage: metrics}
}

func (m MetricsHandler) Other(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "Not found!", http.StatusNotFound)
	return
}

func (m MetricsHandler) Gauge(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only GET requests are allowed!", http.StatusMethodNotAllowed)
		return
	}
	path := strings.Split(r.URL.Path, "/")
	if len(path) < 5 {
		http.Error(w, "Incorrect url", http.StatusNotFound)
		return
	} else if path[4] == "" {
		http.Error(w, "Incorrect url", http.StatusNotFound)
		return
	}

	NameMetric := path[3]
	if value, err := strconv.ParseFloat(path[4], 32); err == nil {
		m.storage.UpdateGauge(NameMetric, value)
		w.WriteHeader(http.StatusOK)
	} else {
		http.Error(w, "The value does not match the type!", http.StatusBadRequest)
		return
	}

}

func (m MetricsHandler) Counter(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only GET requests are allowed!", http.StatusMethodNotAllowed)
		return
	}
	path := strings.Split(r.URL.Path, "/")
	if len(path) < 5 {
		http.Error(w, "Incorrect url", http.StatusNotFound)
		return
	} else if path[4] == "" {
		http.Error(w, "Incorrect url", http.StatusNotFound)
		return
	}

	NameMetric := path[3]
	value, err := strconv.ParseInt(path[4], 10, 64)
	if err != nil {
		http.Error(w, "The value does not match the type!", http.StatusBadRequest)
		return
	} else {
		m.storage.UpdateCounter(NameMetric, value)
		w.WriteHeader(http.StatusOK)
	}

}
