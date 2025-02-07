package main

import (
	"context"
	"log"
	"net/http"
	"platform-service/internal/config"
	"platform-service/internal/database"
	"platform-service/internal/handlers"
	internal_middleware "platform-service/internal/middleware"
	"platform-service/internal/utils"
	"time"

	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
)

func initMetrics() (func(), error) {
	ctx := context.Background()

	opts := []otlpmetricgrpc.Option{
		otlpmetricgrpc.WithTimeout(5 * time.Second),
		otlpmetricgrpc.WithRetry(otlpmetricgrpc.RetryConfig{
			Enabled:        true,
			MaxElapsedTime: 30 * time.Second,
		}),
	}
	exporter, err := otlpmetricgrpc.New(ctx, opts...)
	if err != nil {
		return nil, err
	}
	serviceName, serviceVersion := config.GetAppTelemetryInfo()

	resource, err := resource.New(ctx,
		resource.WithAttributes(
			semconv.ServiceName(serviceName),
			semconv.ServiceVersion(serviceVersion),
		),
	)
	if err != nil {
		return nil, err
	}

	meterProvider := metric.NewMeterProvider(
		metric.WithResource(resource),
		metric.WithReader(
			metric.NewPeriodicReader(
				exporter,
				metric.WithInterval(10*time.Second),
			),
		),
	)

	otel.SetMeterProvider(meterProvider)

	return func() {
		if err := meterProvider.Shutdown(ctx); err != nil {
			log.Printf("Error shutting down meter provider: %v", err)
		}
	}, nil
}

func main() {
	cleanup, err := initMetrics()
	if err != nil {
		log.Fatalf("Failed to initialize OpenTelemetry: %v", err)
	}
	defer cleanup()

	err = database.InitDB()
	if err != nil {
		panic(err)
	}

	metricsMiddleware, err := internal_middleware.NewMetricsMiddleware()
	if err != nil {
		log.Fatalf("Failed to create metrics middleware: %v", err)
	}

	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(metricsMiddleware.Middleware())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: config.GetAllowedOrigins(),
		AllowMethods: []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete, http.MethodOptions},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
	}))

	e.POST("/register", handlers.Register)
	e.POST("/login", handlers.Login)

	r := e.Group("/api")
	r.Use(echojwt.WithConfig(utils.JWTConfig()))

	r.GET("/protected", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{"message": "You are authenticated"})
	}, internal_middleware.AuthMiddleware)

	r.GET("/metrics", handlers.GetMetricsHandler(metricsMiddleware), internal_middleware.AdminAuthMiddleware)

	e.Logger.Fatal(e.Start(":8080"))
}
