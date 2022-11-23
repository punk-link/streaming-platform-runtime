package startup

import (
	"github.com/punk-link/logger"
	runtime "github.com/punk-link/streaming-platform-runtime"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/exporters/prometheus"
	"go.opentelemetry.io/otel/metric/global"
	metricSdk "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	traceSdk "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.12.0"
)

func configureOpenTelemetry(options *runtime.ServiceOptions) {
	configureTracing(options)
	configureMetrics(options)
}

func configureMetrics(options *runtime.ServiceOptions) {
	exporter, err := prometheus.New()
	logOpenTelemetryExceptionIfAny(options.Logger, err)

	metricProvider := metricSdk.NewMeterProvider(metricSdk.WithReader(exporter))
	global.SetMeterProvider(metricProvider)
}

func configureTracing(options *runtime.ServiceOptions) {
	jaegerSettingsValues, err := options.Consul.Get("JaegerSettings")
	if err != nil {
		options.Logger.LogInfo("Jaeger settings is empty")
		return
	}

	jaegerSettings := jaegerSettingsValues.(map[string]any)
	endpoint := jaegerSettings["Endpoint"].(string)
	if endpoint == "" {
		options.Logger.LogInfo("Jaeger endpoint is empty")
		return
	}

	exporter, err := jaeger.New(jaeger.WithCollectorEndpoint(jaeger.WithEndpoint(endpoint)))
	logOpenTelemetryExceptionIfAny(options.Logger, err)

	traceProvider := traceSdk.NewTracerProvider(traceSdk.WithBatcher(exporter), traceSdk.WithResource(resource.NewWithAttributes(
		semconv.SchemaURL,
		semconv.ServiceNameKey.String(options.ServiceName),
		attribute.String("environment", options.EnvironmentName),
	)))

	otel.SetTracerProvider(traceProvider)
}

func logOpenTelemetryExceptionIfAny(logger logger.Logger, err error) {
	if err == nil {
		return
	}

	logger.LogFatal(err, "OpenTelemetry error: %s", err)
}
