package measurements

import (
	"fmt"
	"time"

	"context"

	"github.com/go-logr/logr"
	"github.com/konstellation-io/kai-sdk/go-sdk/internal/common"
	"github.com/konstellation-io/kai-sdk/go-sdk/sdk/metadata"
	"github.com/spf13/viper"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/metric"
	sdkMetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.21.0"
)

const (
	_persistentStorageLoggerName = "[MEASUREMENTS]"
)

type Measurements struct {
	logger        logr.Logger
	metricsClient metric.Meter
	metadata      *metadata.Metadata
}

func New(logger logr.Logger) (*Measurements, error) {
	endpoint := viper.GetString(common.ConfigOpenTelemetryEndpointKey)
	insecure := viper.GetBool(common.ConfigOpenTelemetryInsecureKey)
	timeout := viper.GetInt(common.ConfigOpenTelemetryTimeoutKey)
	interval := viper.GetInt(common.ConfigOpenTelemetryMetricsIntervalKey)
	metadata := metadata.NewMetadata()

	metricsClient, err := initMetrics(logger, endpoint, insecure, timeout, interval, metadata)
	if err != nil {
		return nil, err
	}

	return &Measurements{
		logger:        logger,
		metricsClient: metricsClient,
		metadata:      metadata,
	}, nil
}

func initMetrics(logger logr.Logger, endpoint string, insecure bool, timeout, interval int, metadata *metadata.Metadata) (metric.Meter, error) {
	resource, err := initResource(metadata)
	if err != nil {
		return nil, fmt.Errorf("error initializing metrics: %w", err)
	}

	exporter, err := initExporter(endpoint, insecure, timeout)
	if err != nil {
		return nil, fmt.Errorf("error initializing metrics: %w", err)
	}

	provider := initProvider(exporter, resource, interval)

	otel.SetMeterProvider(provider)

	logger.WithName(_persistentStorageLoggerName).Info("Successfully initialized metrics")

	return provider.Meter("measurements"), nil
}

func initResource(metadata *metadata.Metadata) (*resource.Resource, error) {
	return resource.Merge(resource.Default(),
		resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceName(metadata.GetProduct()),
			semconv.ServiceVersion(metadata.GetVersion()),
			semconv.ServiceNamespace(metadata.GetWorkflow()),
			semconv.ServiceInstanceID(metadata.GetProcess()),
		),
	)
}

func initExporter(endpoint string, insecure bool, timeout int) (*otlpmetricgrpc.Exporter, error) {
	if insecure {
		return otlpmetricgrpc.New(
			context.Background(),
			otlpmetricgrpc.WithEndpoint(endpoint),
			otlpmetricgrpc.WithTimeout(time.Duration(timeout)*time.Second),
			otlpmetricgrpc.WithInsecure(),
		)
	}

	return otlpmetricgrpc.New(
		context.Background(),
		otlpmetricgrpc.WithEndpoint(endpoint),
		otlpmetricgrpc.WithTimeout(time.Duration(timeout)*time.Second),
	)
}

func initProvider(exporter *otlpmetricgrpc.Exporter, resource *resource.Resource, interval int) *sdkMetric.MeterProvider {
	return sdkMetric.NewMeterProvider(
		sdkMetric.WithResource(resource),
		sdkMetric.WithReader(
			sdkMetric.NewPeriodicReader(
				exporter,
				sdkMetric.WithInterval(time.Duration(interval)*time.Second),
			),
		),
	)
}
