package o11y

import (
	"context"
	"log"
	"time"

	"go.opentelemetry.io/contrib/detectors/aws/lambda"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetrichttp"
	"go.opentelemetry.io/otel/metric"
	noopmetric "go.opentelemetry.io/otel/metric/noop"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
)

func initMetricProvider() ShutdownFunc {
	ctx := context.Background()

	exporter, err := otlpmetrichttp.New(
		ctx,
		otlpmetrichttp.WithInsecure(),
	)
	if err != nil {
		log.Printf("failed to create new exporter: %v", err)
		otel.SetMeterProvider(noopmetric.NewMeterProvider())
		return func(ctx context.Context) error { return nil }
	}

	resourceDetector := lambda.NewResourceDetector()
	resource, err := resourceDetector.Detect(ctx)
	if err != nil {
		// just use nil-resource if failed to detect resource
		log.Printf("Failed to create new resource: %v", err)
	}

	provider := sdkmetric.NewMeterProvider(
		sdkmetric.WithResource(resource),
		sdkmetric.WithReader(sdkmetric.NewPeriodicReader(exporter, sdkmetric.WithInterval(1*time.Second))),
	)
	otel.SetMeterProvider(provider)

	return exporter.Shutdown
}

func NewIntCounter(name string) (metric.Int64Counter, error) {
	counter, err := otel.Meter(name).Int64Counter(name, metric.WithUnit("Count"))
	if err != nil {
		return noopmetric.Int64Counter{}, err
	}
	return counter, nil
}

func NewIntHistogram(name string) (metric.Int64Histogram, error) {
	histogram, err := otel.Meter(name).Int64Histogram(name, metric.WithUnit("Count"))
	if err != nil {
		return noopmetric.Int64Histogram{}, err
	}
	return histogram, nil
}
