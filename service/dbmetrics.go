package service

import (
	"database/sql"

	"github.com/prometheus/client_golang/prometheus"
)

type dbStatsCollector struct {
	db *sql.DB

	maxOpenDesc      *prometheus.Desc
	openDesc         *prometheus.Desc
	inUseDesc        *prometheus.Desc
	idleDesc         *prometheus.Desc
	waitCountDesc    *prometheus.Desc
	waitDurDesc      *prometheus.Desc
	maxIdleClosed    *prometheus.Desc
	maxLifetimeClosed *prometheus.Desc
}

func newDBStatsCollector(db *sql.DB) *dbStatsCollector {
	return &dbStatsCollector{
		db: db,
		maxOpenDesc: prometheus.NewDesc(
			"db_connections_max_open",
			"Maximum number of open connections to the database",
			nil, nil,
		),
		openDesc: prometheus.NewDesc(
			"db_connections_open",
			"Current number of open connections to the database",
			nil, nil,
		),
		inUseDesc: prometheus.NewDesc(
			"db_connections_in_use",
			"Current number of open connections in use",
			nil, nil,
		),
		idleDesc: prometheus.NewDesc(
			"db_connections_idle",
			"Current number of idle connections",
			nil, nil,
		),
		waitCountDesc: prometheus.NewDesc(
			"db_connections_wait_count_total",
			"Total number of connections waited for",
			nil, nil,
		),
		waitDurDesc: prometheus.NewDesc(
			"db_connections_wait_duration_seconds_total",
			"Total time blocked waiting for a new connection",
			nil, nil,
		),
		maxIdleClosed: prometheus.NewDesc(
			"db_connections_max_idle_closed_total",
			"Total number of connections closed due to SetMaxIdleConns",
			nil, nil,
		),
		maxLifetimeClosed: prometheus.NewDesc(
			"db_connections_max_lifetime_closed_total",
			"Total number of connections closed due to SetConnMaxLifetime",
			nil, nil,
		),
	}
}

func (c *dbStatsCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.maxOpenDesc
	ch <- c.openDesc
	ch <- c.inUseDesc
	ch <- c.idleDesc
	ch <- c.waitCountDesc
	ch <- c.waitDurDesc
	ch <- c.maxIdleClosed
	ch <- c.maxLifetimeClosed
}

func (c *dbStatsCollector) Collect(ch chan<- prometheus.Metric) {
	stats := c.db.Stats()

	ch <- prometheus.MustNewConstMetric(c.maxOpenDesc, prometheus.GaugeValue, float64(stats.MaxOpenConnections))
	ch <- prometheus.MustNewConstMetric(c.openDesc, prometheus.GaugeValue, float64(stats.OpenConnections))
	ch <- prometheus.MustNewConstMetric(c.inUseDesc, prometheus.GaugeValue, float64(stats.InUse))
	ch <- prometheus.MustNewConstMetric(c.idleDesc, prometheus.GaugeValue, float64(stats.Idle))
	ch <- prometheus.MustNewConstMetric(c.waitCountDesc, prometheus.CounterValue, float64(stats.WaitCount))
	ch <- prometheus.MustNewConstMetric(c.waitDurDesc, prometheus.CounterValue, stats.WaitDuration.Seconds())
	ch <- prometheus.MustNewConstMetric(c.maxIdleClosed, prometheus.CounterValue, float64(stats.MaxIdleClosed))
	ch <- prometheus.MustNewConstMetric(c.maxLifetimeClosed, prometheus.CounterValue, float64(stats.MaxLifetimeClosed))
}

func RegisterDBMetrics(db *sql.DB) {
	prometheus.MustRegister(newDBStatsCollector(db))
}
