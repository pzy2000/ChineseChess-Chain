package monitor_prometheus

import (
	"fmt"
	"sync"

	"github.com/prometheus/client_golang/prometheus"
)

var (
	counterVecs        map[string]*prometheus.CounterVec
	histogramVecs      map[string]*prometheus.HistogramVec
	gaugeVecs          map[string]*prometheus.GaugeVec
	counterVecsMutex   sync.Mutex
	histogramVecsMutex sync.Mutex
	gaugeVecsMutex     sync.Mutex
	namespace          string
)

// init 初始化
func init() {
	counterVecs = make(map[string]*prometheus.CounterVec)
	histogramVecs = make(map[string]*prometheus.HistogramVec)
	gaugeVecs = make(map[string]*prometheus.GaugeVec)
}

// NewCounterVec 新建Counter计数器
// @param subsystem
// @param name
// @param help
// @param labels
// @return *prometheus.CounterVec
func NewCounterVec(subsystem, name, help string, labels ...string) *prometheus.CounterVec {
	counterVecsMutex.Lock()
	defer counterVecsMutex.Unlock()
	s := fmt.Sprintf("%s_%s", subsystem, name)
	if metric, ok := counterVecs[s]; ok {
		return metric
	}
	metric := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: namespace,
			Subsystem: subsystem,
			Name:      name,
			Help:      help,
		}, labels)
	prometheus.MustRegister(metric)
	counterVecs[s] = metric
	return metric
}

// NewHistogramVec 新建统计图HistogramVec
// @param subsystem
// @param name
// @param help
// @param buckets
// @param labels
// @return *prometheus.HistogramVec
func NewHistogramVec(subsystem, name, help string, buckets []float64, labels ...string) *prometheus.HistogramVec {
	histogramVecsMutex.Lock()
	defer histogramVecsMutex.Unlock()
	s := fmt.Sprintf("%s_%s", subsystem, name)
	if metric, ok := histogramVecs[s]; ok {
		return metric
	}
	metric := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: namespace,
			Subsystem: subsystem,
			Name:      name,
			Help:      help,
			Buckets:   buckets,
		}, labels)
	prometheus.MustRegister(metric)
	histogramVecs[s] = metric
	return metric
}

// NewGaugeVec 新建计量器GaugeVec
// @param subsystem
// @param name
// @param help
// @param labels
// @return *prometheus.GaugeVec
func NewGaugeVec(subsystem, name, help string, labels ...string) *prometheus.GaugeVec {
	gaugeVecsMutex.Lock()
	defer gaugeVecsMutex.Unlock()
	s := fmt.Sprintf("%s_%s", subsystem, name)
	if metric, ok := gaugeVecs[s]; ok {
		return metric
	}
	metric := prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: namespace,
			Subsystem: subsystem,
			Name:      name,
			Help:      help,
		}, labels)
	prometheus.MustRegister(metric)
	gaugeVecs[s] = metric
	return metric
}

// NewHistogram 新建统计图Histogram
// @param subsystem
// @param name
// @param help
// @param buckets
// @return *prometheus.Histogram
func NewHistogram(subsystem, name, help string, buckets []float64) *prometheus.Histogram {
	metric := prometheus.NewHistogram(
		prometheus.HistogramOpts{
			Namespace: namespace,
			Subsystem: subsystem,
			Name:      name,
			Help:      help,
			Buckets:   buckets,
		})

	prometheus.MustRegister(metric)
	return &metric
}
