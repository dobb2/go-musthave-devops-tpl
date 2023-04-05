package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/dobb2/go-musthave-devops-tpl/internal/crypto"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	"strconv"

	"github.com/dobb2/go-musthave-devops-tpl/internal/storage"
	"github.com/dobb2/go-musthave-devops-tpl/internal/storage/metrics"
	"github.com/go-chi/chi/v5"
)

type MetricsHandler struct {
	storage storage.MetricCreatorUpdater
}

func New(metrics storage.MetricCreatorUpdater) MetricsHandler {
	return MetricsHandler{storage: metrics}
}

func (m MetricsHandler) GetAllMetrics(w http.ResponseWriter, r *http.Request) {
	metrics, err := m.storage.GetAllMetrics()
	w.Header().Set("Content-Type", "text/html")

	if err != nil {
		http.Error(w, "No metrics", http.StatusBadRequest)
		return
	}

	main := filepath.Join("internal", "static", "dynamicMetricsPage.html")
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

func (m MetricsHandler) PostUpdateMetric(w http.ResponseWriter, r *http.Request) {
	var metric metrics.Metrics
	if err := json.NewDecoder(r.Body).Decode(&metric); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}
	log.Println("value for update ", metric.ID, metric.MType, metric.Hash, *metric.Delta, *metric.Value)
	key := r.Context().Value("Key").(string)

	switch TypeMetric := metric.MType; TypeMetric {
	case "gauge":
		if value := metric.Value; value != nil {
			if crypto.ValidMAC(fmt.Sprintf("%s:gauge:%f", metric.ID, *metric.Value), metric.Hash, key) {
				http.Error(w, "obtained and computed hashes do not match", http.StatusBadRequest)
				return
			}
			m.storage.UpdateGauge(metric.ID, *value)
			w.Header().Set("Content-Type", "text/plain")
			w.WriteHeader(http.StatusOK)
		} else {
			http.Error(w, "the value does not match the type!", http.StatusBadRequest)
			return
		}
	case "counter":
		if delta := metric.Delta; delta != nil {
			if crypto.ValidMAC(fmt.Sprintf("%s:counter:%d", metric.ID, *metric.Delta), metric.Hash, key) {
				http.Error(w, "obtained and computed hashes do not match", http.StatusBadRequest)
				return
			}
			m.storage.UpdateCounter(metric.ID, *delta)
			w.Header().Set("Content-Type", "text/plain")
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

func (m MetricsHandler) PostGetMetric(w http.ResponseWriter, r *http.Request) {
	var metricGet metrics.Metrics
	if err := json.NewDecoder(r.Body).Decode(&metricGet); err != nil {
		w.Header().Set("Content-Type", "text/plain")
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}

	metricSend, err := m.storage.GetValue(metricGet.MType, metricGet.ID)
	log.Println("get value ", metricGet.ID, metricGet.MType, metricGet.Hash, *metricGet.Delta, *metricGet.Value)
	if err != nil {
		w.Header().Set("Content-Type", "text/plain")
		http.Error(w, "not found metric", http.StatusNotFound)
		return
	}
	key := r.Context().Value("Key").(string)
	switch metricSend.MType {
	case "counter":
		metricSend.Hash = crypto.Hash(fmt.Sprintf("%s:counter:%d", metricSend.ID, *metricSend.Delta), key)
	case "gauge":
		metricSend.Hash = crypto.Hash(fmt.Sprintf("%s:gauge:%f", metricSend.ID, *metricSend.Value), key)
	default:
		log.Println("invalid type metric for create hash")
	}

	out, err := json.Marshal(metricSend)
	if err != nil {
		w.Header().Set("Content-Type", "text/plain")
		http.Error(w, "problem marshal metric to json", http.StatusInternalServerError)
		return
	}
	log.Println(string(out))
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(out)
}

func (m MetricsHandler) UpdateMetric(w http.ResponseWriter, r *http.Request) {
	valueStr := chi.URLParam(r, "value")
	nameMetric := chi.URLParam(r, "nameMetric")

	switch TypeMetric := chi.URLParam(r, "typeMetric"); TypeMetric {
	case "gauge":
		if value, err := strconv.ParseFloat(valueStr, 64); err == nil {
			m.storage.UpdateGauge(nameMetric, value)
			w.Header().Set("Content-Type", "text/plain")
			w.WriteHeader(http.StatusOK)
		} else {
			http.Error(w, "The value does not match the type!", http.StatusBadRequest)
			return
		}
	case "counter":
		if value, err := strconv.ParseInt(valueStr, 10, 64); err == nil {
			m.storage.UpdateCounter(nameMetric, value)
			w.Header().Set("Content-Type", "text/plain")
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
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusOK)
		switch typeMetric {
		case "counter":
			fmt.Fprintln(w, *strValue.Delta)
		case "gauge":
			fmt.Fprintln(w, *strValue.Value)
		}
	}
}
