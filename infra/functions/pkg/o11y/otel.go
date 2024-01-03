package o11y

import (
	"context"
	"log"

	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
)

var (
	tracer         trace.Tracer
	tracerProvider *sdktrace.TracerProvider
	tracerShutdown ShutdownFunc
	metricShutdown ShutdownFunc
)

func InitOtel() {
	tracerShutdown = initTracer()
	metricShutdown = initMetricProvider()

	log.Printf("OTEL initialized")
}

func Close(ctx context.Context) error {
	if tracerShutdown != nil {
		if err := tracerShutdown(ctx); err != nil {
			log.Printf("failed to shutdown tracer: %v", err)
		} else {
			log.Println("tracer shutdown")
		}
	}

	if metricShutdown != nil {
		if err := metricShutdown(ctx); err != nil {
			log.Printf("failed to shutdown metric: %v", err)
		} else {
			log.Println("metric shutdown")
		}
	}

	return nil
}
