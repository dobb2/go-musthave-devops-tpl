package postgres

import (
	"database/sql"
	"github.com/dobb2/go-musthave-devops-tpl/internal/storage/metrics"
)

type MetricsStorer struct {
	db *sql.DB
}

func Create(db *sql.DB) (*MetricsStorer, error) {
	query := `
		CREATE TABLE IF NOT EXISTS Metric (
    	id varchar(100) PRIMARY KEY,
    	mtype varchar(100) NOT NULL,
    	delta bigint,
    	val double precision
	);
	`
	_, err := db.Exec(query)
	if err != nil {
		return &MetricsStorer{}, err
	}
	return &MetricsStorer{db: db}, nil
}

func (m MetricsStorer) UpdateGauge(nameMetric string, value float64) error {
	query := `INSERT INTO Metric (id, mtype, val) VALUES($1, 'gauge', $2)
	ON CONFLICT (id)
    	DO
        UPDATE SET val = $2`

	_, err := m.db.Exec(query, nameMetric, value)
	if err != nil {
		return err
	}

	return nil
}

func (m MetricsStorer) UpdateCounter(nameMetric string, value int64) error {
	query := `INSERT INTO Metric (id, mtype, delta) VALUES($1, 'counter', $2)
	ON CONFLICT (id)
    	DO
        UPDATE SET delta = $2`

	_, err := m.db.Exec(query, nameMetric, value)
	if err != nil {
		return err
	}

	return nil
}

func (m MetricsStorer) GetValue(typeMetric string, NameMetric string) (metrics.Metrics, error) {
	query := `
		SELECT 
    		id,
    		mtype,
    		delta,
    		val
		FROM Metric
		WHERE 
		    id = $1 AND mtype = $2
	`
	var metric metrics.Metrics

	err := m.db.QueryRow(query, NameMetric, typeMetric).
		Scan(&metric.ID, &metric.MType, &metric.Delta, &metric.Value)

	if err != nil {
		return metrics.Metrics{}, err
	}

	return metric, nil

}

func (m MetricsStorer) GetAllMetrics() ([]metrics.Metrics, error) {
	query := `
		SELECT 
    		id,
    		mtype,
    		delta,
    		val
		FROM Metric
	`
	rows, err := m.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	batchMetric := make([]metrics.Metrics, 3)
	for rows.Next() {
		var rowMetric metrics.Metrics
		err = rows.Scan(&rowMetric.ID, &rowMetric.MType, &rowMetric.Delta, &rowMetric.Value)
		if err != nil {
			return nil, err
		}

		batchMetric = append(batchMetric, rowMetric)
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return batchMetric, nil
}

func (m MetricsStorer) GetPing() error {
	return m.db.Ping()
}

func (m MetricsStorer) UploadMetrics([]metrics.Metrics) error {
	return nil
}

func (m MetricsStorer) AddChannel(*chan struct{}) error {
	return nil
}
