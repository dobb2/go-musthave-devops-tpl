package handlers

import (
	"fmt"
	"github.com/dobb2/go-musthave-devops-tpl/internal/storage"
	"github.com/go-chi/chi/v5"
	"html/template"
	"net/http"
	"path/filepath"
	"strconv"
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
	valueStr := chi.URLParam(r, "value")
	nameMetric := chi.URLParam(r, "nameMetric")

	switch TypeMetric := chi.URLParam(r, "typeMetric"); TypeMetric {
	case "gauge":
		if value, err := strconv.ParseFloat(valueStr, 64); err == nil {
			m.storage.UpdateGauge(nameMetric, value)
			w.WriteHeader(http.StatusOK)
		} else {
			http.Error(w, "The value does not match the type!", http.StatusBadRequest)
			return
		}
	case "counter":
		if value, err := strconv.ParseInt(valueStr, 10, 64); err == nil {
			m.storage.UpdateCounter(nameMetric, value)
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
	typeMetric := chi.URLParam(r, "typeMetric")
	nameMetric := chi.URLParam(r, "nameMetric")

	strValue, err := m.storage.GetValue(typeMetric, nameMetric)
	if err != nil {
		http.Error(w, "Not found metric", http.StatusNotFound)
		return
	} else {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, strValue.Value)
	}
}
