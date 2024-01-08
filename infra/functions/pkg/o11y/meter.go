package o11y

import (
	"context"
	"log"
	"time"

	"go.opentelemetry.io/contrib/detectors/aws/lambda"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/metric"
	noopmetric "go.opentelemetry.io/otel/metric/noop"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
)

func initMetricProvider() ShutdownFunc {
	ctx := context.Background()

	mp, err := newMetricProvider(ctx)
	// use noop metric provider if xray is not available
	if err != nil {
		log.Printf("failed to create metric provider: %v", err)
		otel.SetMeterProvider(noopmetric.NewMeterProvider())
		return func(ctx context.Context) error { return nil }
	}
	otel.SetMeterProvider(mp)

	return mp.Shutdown
}

func newMetricProvider(ctx context.Context) (*sdkmetric.MeterProvider, error) {
	exporter, err := otlpmetricgrpc.New(
		ctx,
		otlpmetricgrpc.WithInsecure(),
	)
	if err != nil {
		return nil, err
	}

	resourceDetector := lambda.NewResourceDetector()
	resource, err := resourceDetector.Detect(ctx)
	if err != nil {
		return nil, err
	}

	provider := sdkmetric.NewMeterProvider(
		sdkmetric.WithResource(resource),
		sdkmetric.WithReader(sdkmetric.NewPeriodicReader(exporter, sdkmetric.WithInterval(time.Second*1))),
	)
	return provider, nil
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
