package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/dobb2/go-musthave-devops-tpl/internal/crypto"
	"github.com/rs/zerolog"
	"html/template"
	"net/http"
	"path/filepath"
	"strconv"

	"github.com/dobb2/go-musthave-devops-tpl/internal/storage"
	"github.com/dobb2/go-musthave-devops-tpl/internal/storage/metrics"
	"github.com/go-chi/chi/v5"
)

type MetricsHandler struct {
	storage storage.MetricGetterCreatorUpdater
	logger  zerolog.Logger
}

func New(metrics storage.MetricGetterCreatorUpdater, logger zerolog.Logger) MetricsHandler {
	return MetricsHandler{
		storage: metrics,
		logger:  logger,
	}
}

func (m MetricsHandler) GetPing(w http.ResponseWriter, r *http.Request) {
	err := m.storage.GetPing()
	if err != nil {
		m.logger.Error().Err(err).Msg("cannot connect to db")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
}

func (m MetricsHandler) GetAllMetrics(w http.ResponseWriter, r *http.Request) {
	metrics, err := m.storage.GetAllMetrics()
	w.Header().Set("Content-Type", "text/html")

	if err != nil {
		m.logger.Warn().Err(err).Msg("cannot get all metrics")
		http.Error(w, "No metrics", http.StatusBadRequest)
		return
	}

	//main := filepath.Join("..", "..", "internal", "static", "dynamicMetricsPage.html")
	//Join for autotests in git
	main := filepath.Join("internal", "static", "dynamicMetricsPage.html")
	tmpl, err := template.ParseFiles(main)
	if err != nil {
		m.logger.Warn().Stack().Err(err).Msg("problem parse files for template html")
		http.Error(w, "", http.StatusBadRequest)
		return
	}
	err = tmpl.ExecuteTemplate(w, "metrics", metrics)
	if err != nil {
		m.logger.Warn().Stack().Err(err)
		http.Error(w, "", http.StatusBadRequest)
		return
	}
}

func (m MetricsHandler) PostUpdateMetric(w http.ResponseWriter, r *http.Request) {
	var metric metrics.Metrics
	if err := json.NewDecoder(r.Body).Decode(&metric); err != nil {
		m.logger.Debug().Err(err).Msg("invalid json")
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}

	key := r.Context().Value("Key").(string)
	switch TypeMetric := metric.MType; TypeMetric {
	case "gauge":
		if value := metric.Value; value != nil {
			if !crypto.ValidMAC(fmt.Sprintf("%s:gauge:%f", metric.ID, *metric.Value), metric.Hash, key) {
				m.logger.Debug().Msg("obtained and computed hashes do not match")
				m.logger.Debug().Msg("gauge")
				m.logger.Debug().Msg("Key" + key)
				m.logger.Debug().Msg("Hash" + metric.Hash)

				m.logger.Debug().Msg(fmt.Sprintf("%f", *metric.Value))
				http.Error(w, "obtained and computed hashes do not match", http.StatusBadRequest)
				return
			}

			err := m.storage.UpdateGauge(metric.ID, *value)
			if err != nil {
				m.logger.Error().Stack().Err(err).Msg("problem with update value in storage")
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			w.Header().Set("Content-Type", "text/plain")
			w.WriteHeader(http.StatusOK)
		} else {
			http.Error(w, "the value does not match the type!", http.StatusBadRequest)
			return
		}
	case "counter":
		if delta := metric.Delta; delta != nil {
			if !crypto.ValidMAC(fmt.Sprintf("%s:counter:%d", metric.ID, *metric.Delta), metric.Hash, key) {
				m.logger.Debug().Msg("counter")
				m.logger.Debug().Msg("Key" + key)
				m.logger.Debug().Msg("Hash" + metric.Hash)
				m.logger.Debug().Msg(fmt.Sprintf("%d", *metric.Delta))
				m.logger.Debug().Msg("obtained and computed hashes do not match")
				http.Error(w, "obtained and computed hashes do not match", http.StatusBadRequest)
				return
			}

			err := m.storage.UpdateCounter(metric.ID, *delta)
			if err != nil {
				m.logger.Error().Stack().Err(err).Msg("problem with update value in storage")
				http.Error(w, "", http.StatusInternalServerError)
			}
			w.Header().Set("Content-Type", "text/plain")
			w.WriteHeader(http.StatusOK)
		} else {
			m.logger.Debug().Msg("the value does not match the expected type:")
			http.Error(w, "The value does not match the type!", http.StatusBadRequest)
			return
		}
	default:
		m.logger.Debug().Msg("The resulting value has an unknown type")
		http.Error(w, "Invalid type metric", http.StatusNotImplemented)
		return
	}
}

func (m MetricsHandler) PostGetMetric(w http.ResponseWriter, r *http.Request) {
	var metricGet metrics.Metrics
	if err := json.NewDecoder(r.Body).Decode(&metricGet); err != nil {
		m.logger.Debug().Err(err).Msg("invalid json")
		w.Header().Set("Content-Type", "text/plain")
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}

	metricSend, err := m.storage.GetValue(metricGet.MType, metricGet.ID)
	if err != nil {
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
		m.logger.Debug().Msg("invalid type metric for create hash")
	}
	out, err := json.Marshal(metricSend)
	if err != nil {
		m.logger.Debug().Err(err).Msg("invalid marshal json")
		http.Error(w, "problem marshal metric to json", http.StatusInternalServerError)
		return
	}
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
			err = m.storage.UpdateGauge(nameMetric, value)
			if err != nil {
				m.logger.Error().Stack().Err(err).Msg("problem with update value in storage")
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			w.Header().Set("Content-Type", "text/plain")
			w.WriteHeader(http.StatusOK)
		} else {
			m.logger.Debug().Err(err).Msg("value does not match the type")
			http.Error(w, "The value does not match the type!", http.StatusBadRequest)
			return
		}
	case "counter":
		if value, err := strconv.ParseInt(valueStr, 10, 64); err == nil {
			err = m.storage.UpdateCounter(nameMetric, value)
			if err != nil {
				m.logger.Error().Stack().Err(err).Msg("problem with update value in storage")
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			w.Header().Set("Content-Type", "text/plain")
			w.WriteHeader(http.StatusOK)
		} else {
			m.logger.Debug().Err(err).Msg("value does not match the type")
			http.Error(w, "The value does not match the type!", http.StatusBadRequest)
			return
		}
	default:
		m.logger.Debug().Msg("invalid type metric for create hash")
		http.Error(w, "Invalid type metric", http.StatusNotImplemented)
		return
	}
}

func (m MetricsHandler) GetMetric(w http.ResponseWriter, r *http.Request) {
	typeMetric := chi.URLParam(r, "typeMetric")
	nameMetric := chi.URLParam(r, "nameMetric")

	strValue, err := m.storage.GetValue(typeMetric, nameMetric)
	if err != nil {
		m.logger.Debug().Err(err).Msg("Not found metric")
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

func (m MetricsHandler) PostUpdateBatchMetrics(w http.ResponseWriter, r *http.Request) {
	var metrics []metrics.Metrics

	if err := json.NewDecoder(r.Body).Decode(&metrics); err != nil {
		m.logger.Debug().Err(err).Msg("invalid json")
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}
	key := r.Context().Value("Key").(string)

	for i := range metrics {
		switch TypeMetric := metrics[i].MType; TypeMetric {
		case "gauge":
			if value := metrics[i].Value; value != nil {
				if !crypto.ValidMAC(fmt.Sprintf("%s:gauge:%f", metrics[i].ID, *metrics[i].Value), metrics[i].Hash, key) {
					m.logger.Debug().Msg("obtained and computed hashes do not match")
					m.logger.Debug().Msg(metrics[i].ID)
					m.logger.Debug().Msg("gauge")
					m.logger.Debug().Msg("Key" + key)
					m.logger.Debug().Msg("Hash" + metrics[i].Hash)
					m.logger.Debug().Msg(fmt.Sprintf("%f", *metrics[i].Value))
					http.Error(w, "", http.StatusBadRequest)
					return
				}
			} else {
				m.logger.Debug().Msg("the value for does not match the type")
				http.Error(w, "", http.StatusBadRequest)
				return
			}
		case "counter":
			if delta := metrics[i].Delta; delta != nil {
				if !crypto.ValidMAC(fmt.Sprintf("%s:counter:%d", metrics[i].ID, *metrics[i].Delta), metrics[i].Hash, key) {
					m.logger.Debug().Msg("obtained and computed hashes do not match")
					m.logger.Debug().Msg("counter")
					m.logger.Debug().Msg("Key" + key)
					m.logger.Debug().Msg("Hash" + metrics[i].Hash)
					m.logger.Debug().Msg(fmt.Sprintf("%d", *metrics[i].Delta))
					http.Error(w, "", http.StatusBadRequest)
					return
				}
			} else {
				m.logger.Debug().Msg("the value for does not match the type")
				http.Error(w, "", http.StatusBadRequest)
				return
			}
		default:
			m.logger.Debug().Msg("unknown type for value")
			http.Error(w, "unknown type for"+metrics[i].ID, http.StatusBadRequest)
			return
		}
	}

	err := m.storage.UpdateMetrics(metrics)
	if err != nil {
		m.logger.Error().Stack().Err(err).Msg("metrics are not updated:")
		http.Error(w, "", http.StatusInternalServerError)
	}

	out, err := json.Marshal(metrics)
	if err != nil {
		m.logger.Info().Stack().Err(err).Msg("metrics are not marshaling")
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(out)
}
