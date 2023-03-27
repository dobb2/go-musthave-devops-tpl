package backup

import (
	"encoding/json"
	"github.com/dobb2/go-musthave-devops-tpl/internal/config"
	"github.com/dobb2/go-musthave-devops-tpl/internal/storage"
	"github.com/dobb2/go-musthave-devops-tpl/internal/storage/metrics"
	"io"
	"log"
	"os"
)

type MetricsBackuper struct {
	storage storage.MetricsBackuper
}

func New(metrics storage.MetricsBackuper) MetricsBackuper {
	return MetricsBackuper{storage: metrics}
}

func (b MetricsBackuper) Restore(cfg config.Config) {
	jsonFile, err := os.OpenFile(cfg.StoreFile, os.O_CREATE|os.O_RDONLY, 0777)
	if err != nil {
		log.Println(err)
	}
	defer jsonFile.Close()

	byteValue, err := io.ReadAll(jsonFile)
	if err != nil {
		log.Println(err)
	}
	metrics := make([]metrics.Metrics, 0)

	json.Unmarshal(byteValue, &metrics)

	b.storage.UploadMetrics(metrics)
	log.Println("old metrics download to storage")
}

func (b MetricsBackuper) UpdateBackup(cfg config.Config) {
	metrics, err := b.storage.GetAllMetrics()
	if err != nil {
		log.Println(err)
	}

	dataMetrics, err := json.Marshal(metrics)
	if err != nil {
		log.Println(err)
	}
	dataMetrics = append(dataMetrics, '\n')
	file, err := os.OpenFile(cfg.StoreFile, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0777)
	if err != nil {
		log.Println(err)
	}
	defer file.Close()

	_, err = file.Write(dataMetrics)
	if err != nil {
		log.Println(err)
	}
	log.Println("backup done")
}
