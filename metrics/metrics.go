package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"go.uber.org/multierr"
)

const (
	NAMESPACE = "golib_"
	SUBSYSTEM = "backend_"
)

var (
	State                 prometheus.Gauge
	Uptime                prometheus.Gauge
	CountOnlineUser       prometheus.Gauge
	HttpRequestsTotal     *prometheus.CounterVec
	HttpRequestDuration   *prometheus.HistogramVec
	ResponseTime          prometheus.Histogram
	ErrorRate             prometheus.Counter
	MemoryTotal           prometheus.Gauge
	MemoryUsed            prometheus.Gauge
	MemoryUsedPercent     prometheus.Gauge
	CpuCount              prometheus.Gauge
	CpuUsedPercent        prometheus.Gauge
	ConcurrentConnections prometheus.Gauge
	DbConnectionsOpen     prometheus.Gauge
	CacheHit              *prometheus.CounterVec
	CacheMiss             *prometheus.CounterVec
	QueueSize             prometheus.Gauge
)

func Init() error {
	State = prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace: NAMESPACE,
		Subsystem: SUBSYSTEM,
		Name:      "state",
		Help:      "The state of the backend",
	})
	Uptime = prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace: NAMESPACE,
		Subsystem: SUBSYSTEM,
		Name:      "uptime",
		Help:      "The uptime of the backend",
	})
	// HttpRequestsTotal.WithLabelValues("GET").Inc()
	// HttpRequestsTotal.WithLabelValues("POST").Inc()
	HttpRequestsTotal = prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: NAMESPACE,
		Subsystem: SUBSYSTEM,
		Name:      "http_requests_total",
		Help:      "Total number of HTTP requests",
	},
		[]string{"method", "path", "status"},
	)
	HttpRequestDuration = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Namespace: NAMESPACE,
		Subsystem: SUBSYSTEM,
		Name:      "http_request_duration_seconds",
		Help:      "HTTP request latencies in seconds",
		Buckets:   prometheus.DefBuckets,
	},
		[]string{"method", "path", "status"},
	)

	CountOnlineUser = prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace: NAMESPACE,
		Subsystem: SUBSYSTEM,
		Name:      "count_online",
		Help:      "The count of online user",
	})

	ResponseTime = prometheus.NewHistogram(prometheus.HistogramOpts{
		Namespace: NAMESPACE,
		Subsystem: SUBSYSTEM,
		Name:      "response_time_seconds",
		Help:      "Response time in seconds",
		Buckets:   prometheus.DefBuckets,
	})
	ErrorRate = prometheus.NewCounter(prometheus.CounterOpts{
		Namespace: NAMESPACE,
		Subsystem: SUBSYSTEM,
		Name:      "error_total",
		Help:      "Total number of errors",
	})
	MemoryTotal = prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace: NAMESPACE,
		Subsystem: SUBSYSTEM,
		Name:      "memory_usage_total_bytes",
		Help:      "Current memory total in bytes",
	})
	MemoryUsed = prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace: NAMESPACE,
		Subsystem: SUBSYSTEM,
		Name:      "memory_usage_used_bytes",
		Help:      "Current memory used in bytes",
	})
	MemoryUsedPercent = prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace: NAMESPACE,
		Subsystem: SUBSYSTEM,
		Name:      "memory_usage_percent",
		Help:      "Current memory used in percent",
	})
	CpuCount = prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace: NAMESPACE,
		Subsystem: SUBSYSTEM,
		Name:      "cpu_count",
		Help:      "Current cpu count",
	})
	CpuUsedPercent = prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace: NAMESPACE,
		Subsystem: SUBSYSTEM,
		Name:      "cpu_used_percent",
		Help:      "Current cpu used in percent",
	})
	ConcurrentConnections = prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace: NAMESPACE,
		Subsystem: SUBSYSTEM,
		Name:      "concurrent_connections",
		Help:      "Number of concurrent connections",
	})
	DbConnectionsOpen = prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace: NAMESPACE,
		Subsystem: SUBSYSTEM,
		Name:      "db_connections_open",
		Help:      "Number of open database connections",
	})
	CacheHit = prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: NAMESPACE,
		Subsystem: SUBSYSTEM,
		Name:      "cache_hits_total",
		Help:      "Total number of cache hits",
	}, []string{"phase", "table"})
	CacheMiss = prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: NAMESPACE,
		Subsystem: SUBSYSTEM,
		Name:      "cache_misses_total",
		Help:      "Total number of cache misses",
	}, []string{"phase", "table"})
	QueueSize = prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace: NAMESPACE,
		Subsystem: SUBSYSTEM,
		Name:      "queue_size",
		Help:      "Current size of the task queue",
	})

	errs := make([]error, 0)
	errs = append(errs, prometheus.Register(State))
	errs = append(errs, prometheus.Register(Uptime))
	errs = append(errs, prometheus.Register(HttpRequestsTotal))
	errs = append(errs, prometheus.Register(CountOnlineUser))
	errs = append(errs, prometheus.Register(ResponseTime))
	errs = append(errs, prometheus.Register(ErrorRate))
	errs = append(errs, prometheus.Register(MemoryTotal))
	errs = append(errs, prometheus.Register(MemoryUsed))
	errs = append(errs, prometheus.Register(MemoryUsedPercent))
	errs = append(errs, prometheus.Register(CpuCount))
	errs = append(errs, prometheus.Register(CpuUsedPercent))
	errs = append(errs, prometheus.Register(ConcurrentConnections))
	errs = append(errs, prometheus.Register(DbConnectionsOpen))
	errs = append(errs, prometheus.Register(CacheHit))
	errs = append(errs, prometheus.Register(CacheMiss))
	errs = append(errs, prometheus.Register(QueueSize))

	errs = append(errs, prometheus.Register(collectors.NewBuildInfoCollector()))
	errs = append(errs, prometheus.Register(collectors.NewProcessCollector(collectors.ProcessCollectorOpts{Namespace: NAMESPACE})))
	// errs = append(errs, prometheus.Register(collectors.NewGoCollector()))
	// errs = append(errs, prometheus.Register(collectors.NewGoCollector(
	// 	collectors.WithGoCollections(collectors.GoRuntimeMetricsCollection),
	// 	collectors.WithGoCollections(collectors.GoRuntimeMemStatsCollection),
	// )))
	return multierr.Combine(errs...)
}
