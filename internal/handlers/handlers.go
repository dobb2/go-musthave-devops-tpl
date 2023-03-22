package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/dobb2/go-musthave-devops-tpl/internal/entities"
	"github.com/dobb2/go-musthave-devops-tpl/internal/storage"
	"html/template"
	"net/http"
	"path/filepath"
)

type MetricsHandler struct {
	storage storage.MetricCreatorUpdater
}

func New(metrics storage.MetricCreatorUpdater) MetricsHandler {
	return MetricsHandler{storage: metrics}
}

func (m MetricsHandler) GetAllMetrics(w http.ResponseWriter, r *http.Request) {
	metrics, err := m.storage.GetAllMetrics()
	if err != nil {
		http.Error(w, "No metrics", http.StatusBadRequest)
		return
	}

	main := filepath.Join("..", "..", "internal", "static", "dynamicMetricsPage.html")

	tmpl, err := template.ParseFiles(main)
	if err != nil {
		http.Error(w, "", http.StatusBadRequest)
		return
	}

	err = tmpl.ExecuteTemplate(w, "metrics", metrics)
	if err != nil {
		http.Error(w, "", http.StatusBadRequest)
		return
	}
}

func (m MetricsHandler) UpdateMetric(w http.ResponseWriter, r *http.Request) {
	var metric entities.Metrics // целевой объект

	if err := json.NewDecoder(r.Body).Decode(&metric); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}

	switch TypeMetric := metric.MType; TypeMetric {
	case "gauge":
		if value := metric.Value; value != nil {
			m.storage.UpdateGauge(metric.ID, *value)
			w.WriteHeader(http.StatusOK)
		} else {
			http.Error(w, "the value does not match the type!", http.StatusBadRequest)
			return
		}
	case "counter":
		if delta := metric.Delta; delta != nil {
			m.storage.UpdateCounter(metric.ID, *delta)
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

func (m MetricsHandler) GetMetric(w http.ResponseWriter, r *http.Request) {
	var metric entities.Metrics

	if err := json.NewDecoder(r.Body).Decode(&metric); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}

	strValue, err := m.storage.GetValue(metric.MType, metric.ID)
	if err != nil {
		http.Error(w, "Not found metric", http.StatusNotFound)
		return
	} else {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, strValue.Value)
	}
}
