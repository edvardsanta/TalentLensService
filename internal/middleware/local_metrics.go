package internal_middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/labstack/echo/v4"
)

type LocalMetrics struct {
	mu            sync.RWMutex
	RequestCounts map[string]map[string]int64
	StatusCounts  map[string]map[string]map[int]int64
	Durations     map[string]map[string][]float64
	StartTime     time.Time
}

type MetricsSummary struct {
	Uptime          string                              `json:"uptime"`
	TotalRequests   int64                               `json:"total_requests"`
	RequestsPerPath map[string]map[string]int64         `json:"requests_per_path"`
	StatusCodes     map[string]map[string]map[int]int64 `json:"status_codes"`
	AverageDuration map[string]map[string]float64       `json:"average_duration"`
}

func NewLocalMetrics() *LocalMetrics {
	return &LocalMetrics{
		RequestCounts: make(map[string]map[string]int64),
		StatusCounts:  make(map[string]map[string]map[int]int64),
		Durations:     make(map[string]map[string][]float64),
		StartTime:     time.Now(),
	}
}

func (m *LocalMetrics) RecordRequest(method, path string) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.RequestCounts[method] == nil {
		m.RequestCounts[method] = make(map[string]int64)
	}
	m.RequestCounts[method][path]++
}

func (m *LocalMetrics) RecordStatus(method, path string, status int) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.StatusCounts[method] == nil {
		m.StatusCounts[method] = make(map[string]map[int]int64)
	}
	if m.StatusCounts[method][path] == nil {
		m.StatusCounts[method][path] = make(map[int]int64)
	}
	m.StatusCounts[method][path][status]++
}

func (m *LocalMetrics) RecordDuration(method, path string, duration float64) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.Durations[method] == nil {
		m.Durations[method] = make(map[string][]float64)
	}
	m.Durations[method][path] = append(m.Durations[method][path], duration)
}

func (m *LocalMetrics) GetSummary() MetricsSummary {
	m.mu.RLock()
	defer m.mu.RUnlock()

	summary := MetricsSummary{
		Uptime:          time.Since(m.StartTime).Round(time.Second).String(),
		RequestsPerPath: make(map[string]map[string]int64),
		StatusCodes:     make(map[string]map[string]map[int]int64),
		AverageDuration: make(map[string]map[string]float64),
	}

	for method, paths := range m.RequestCounts {
		summary.RequestsPerPath[method] = make(map[string]int64)
		for path, count := range paths {
			summary.RequestsPerPath[method][path] = count
			summary.TotalRequests += count
		}
	}

	for method, paths := range m.StatusCounts {
		summary.StatusCodes[method] = make(map[string]map[int]int64)
		for path, statuses := range paths {
			summary.StatusCodes[method][path] = make(map[int]int64)
			for status, count := range statuses {
				summary.StatusCodes[method][path][status] = count
			}
		}
	}

	summary.AverageDuration = make(map[string]map[string]float64)
	for method, paths := range m.Durations {
		summary.AverageDuration[method] = make(map[string]float64)
		for path, durations := range paths {
			if len(durations) > 0 {
				var sum float64
				for _, d := range durations {
					sum += d
				}
				summary.AverageDuration[method][path] = sum / float64(len(durations))
			}
		}
	}

	return summary
}

func (m *LocalMetrics) MetricsHandler(c echo.Context) error {
	summary := m.GetSummary()
	return c.JSON(http.StatusOK, summary)
}
