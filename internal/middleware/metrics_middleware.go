package internal_middleware

import (
	"context"
	"platform-service/internal/config"
	"time"

	"github.com/labstack/echo/v4"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
)

type MetricsMiddleware struct {
	requestCounter  metric.Int64Counter
	requestDuration metric.Float64Histogram
	responseCounter metric.Int64Counter
	meter           metric.Meter
	LocalMetrics    *LocalMetrics
	UseLocalMetrics bool
}

func NewMetricsMiddleware() (*MetricsMiddleware, error) {
	UseLocalMetrics := config.IsOpenTelemetryDisabled()

	m := &MetricsMiddleware{
		UseLocalMetrics: UseLocalMetrics,
	}

	if UseLocalMetrics {
		m.LocalMetrics = NewLocalMetrics()
	} else {
		serviceName, instrumentationVersion := config.GetAppTelemetryInfo()
		meter := otel.GetMeterProvider().Meter(
			serviceName,
			metric.WithInstrumentationVersion(instrumentationVersion),
		)

		var err error
		m.requestCounter, err = meter.Int64Counter(
			"http.server.request_count",
			metric.WithDescription("Total number of HTTP requests"),
		)
		if err != nil {
			return nil, err
		}

		m.requestDuration, err = meter.Float64Histogram(
			"http.server.duration",
			metric.WithDescription("Duration of HTTP requests"),
			metric.WithUnit("ms"),
		)
		if err != nil {
			return nil, err
		}

		m.responseCounter, err = meter.Int64Counter(
			"http.server.response_count",
			metric.WithDescription("Total number of HTTP responses by status code"),
		)
		if err != nil {
			return nil, err
		}

		m.meter = meter
	}

	return m, nil
}

func (m *MetricsMiddleware) Middleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			start := time.Now()
			path := c.Path()
			method := c.Request().Method

			if m.UseLocalMetrics {
				m.LocalMetrics.RecordRequest(method, path)
			} else {
				attrs := []attribute.KeyValue{
					attribute.String("http.method", method),
					attribute.String("http.path", path),
				}
				m.requestCounter.Add(context.Background(), 1, metric.WithAttributes(attrs...))
			}

			err := next(c)

			duration := float64(time.Since(start).Milliseconds())

			if m.UseLocalMetrics {
				m.LocalMetrics.RecordDuration(method, path, duration)
				m.LocalMetrics.RecordStatus(method, path, c.Response().Status)
			} else {
				attrs := []attribute.KeyValue{
					attribute.String("http.method", method),
					attribute.String("http.path", path),
				}
				m.requestDuration.Record(context.Background(), duration,
					metric.WithAttributes(attrs...))

				responseAttrs := append(attrs, attribute.Int("http.status_code", c.Response().Status))
				m.responseCounter.Add(context.Background(), 1,
					metric.WithAttributes(responseAttrs...))
			}

			return err
		}
	}
}
