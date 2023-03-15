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

func (m MetricsHandler) Update(w http.ResponseWriter, r *http.Request) {
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

	ValueStr := path[4]
	NameMetric := path[3]
	switch TypeMetric := path[2]; TypeMetric {
	case "gauge":
		if value, err := strconv.ParseFloat(ValueStr, 32); err == nil {
			m.storage.UpdateGauge(NameMetric, value)
			w.WriteHeader(http.StatusOK)
		} else {
			http.Error(w, "The value does not match the type!", http.StatusBadRequest)
			return
		}
	case "counter":
		if value, err := strconv.ParseInt(path[4], 10, 64); err == nil {
			m.storage.UpdateCounter(NameMetric, value)
			w.WriteHeader(http.StatusOK)
		} else {
			http.Error(w, "The value does not match the type!", http.StatusBadRequest)
			return
		}
	default:
		http.Error(w, "Invalid type metric", http.StatusNotImplemented)
		return
	}

}
