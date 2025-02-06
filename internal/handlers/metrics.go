package handlers

import (
	internal_middleware "platform-service/internal/middleware"

	"github.com/labstack/echo/v4"
)

func GetMetricsHandler(m *internal_middleware.MetricsMiddleware) echo.HandlerFunc {
	if !m.UseLocalMetrics {
		return func(c echo.Context) error {
			return c.JSON(404, map[string]string{
				"error": "Metrics endpoint only available when OTEL_SDK_DISABLED=true",
			})
		}
	}
	return m.LocalMetrics.MetricsHandler
}
