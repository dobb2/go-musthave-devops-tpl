package backup

import (
	"encoding/json"
	"io"
	"os"

	"github.com/dobb2/go-musthave-devops-tpl/internal/config"
	"github.com/dobb2/go-musthave-devops-tpl/internal/storage"
	"github.com/dobb2/go-musthave-devops-tpl/internal/storage/metrics"
)

type MetricsBackuper struct {
	storage storage.MetricCreatorUpdaterBackuper
}

func New(metrics storage.MetricCreatorUpdaterBackuper) MetricsBackuper {
	return MetricsBackuper{storage: metrics}
}

func (b MetricsBackuper) Restore(cfg config.ServerConfig) error {
	jsonFile, err := os.OpenFile(cfg.StoreFile, os.O_CREATE|os.O_RDONLY, 0777)
	if err != nil {
		return err
	}
	defer jsonFile.Close()

	byteValue, err := io.ReadAll(jsonFile)
	if err != nil {
		return err
	}
	metrics := make([]metrics.Metrics, 0)

	json.Unmarshal(byteValue, &metrics)

	b.storage.UploadMetrics(metrics)
	return nil
}

func (b MetricsBackuper) UpdateBackup(cfg config.ServerConfig) error {
	metrics, err := b.storage.GetAllMetrics()
	if err != nil {
		return err
	}

	dataMetrics, err := json.Marshal(metrics)
	if err != nil {
		return err
	}
	dataMetrics = append(dataMetrics, '\n')
	file, err := os.OpenFile(cfg.StoreFile, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0777)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.Write(dataMetrics)
	if err != nil {
		return err
	}
	return nil
}
