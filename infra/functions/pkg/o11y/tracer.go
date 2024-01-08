package o11y

import (
	"context"
	"fmt"
	"log"

	lambdadetector "go.opentelemetry.io/contrib/detectors/aws/lambda"
	"go.opentelemetry.io/contrib/instrumentation/github.com/aws/aws-lambda-go/otellambda"
	"go.opentelemetry.io/contrib/instrumentation/github.com/aws/aws-lambda-go/otellambda/xrayconfig"
	"go.opentelemetry.io/contrib/propagators/aws/xray"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.10.0"
	"go.opentelemetry.io/otel/trace"
	nooptrace "go.opentelemetry.io/otel/trace/noop"
)

func initTracer() ShutdownFunc {
	ctx := context.Background()

	tp, err := newTracerProvider(ctx)
	// use noop tracer if xray is not available
	if err != nil {
		log.Printf("error creating tracer provider: %v", err)
		otel.SetTracerProvider(nooptrace.NewTracerProvider())
		otel.SetTextMapPropagator(xray.Propagator{})
		tracer = otel.Tracer("composebold")
		return func(ctx context.Context) error { return nil }
	}

	tracerProvider = tp
	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(xray.Propagator{})
	tracer = otel.Tracer("composebold")
	return tp.Shutdown
}

func newTracerProvider(ctx context.Context) (*sdktrace.TracerProvider, error) {
	exp, err := otlptracegrpc.New(ctx, otlptracegrpc.WithInsecure())
	if err != nil {
		return nil, err
	}

	detector := lambdadetector.NewResourceDetector()
	resource, err := detector.Detect(ctx)
	if err != nil {
		return nil, err
	}

	return sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
		sdktrace.WithBatcher(exp),
		sdktrace.WithIDGenerator(xray.NewIDGenerator()),
		sdktrace.WithResource(resource),
	), nil
}

func SpanFromContext(ctx context.Context) trace.Span {
	return trace.SpanFromContext(ctx)
}

func BeginSpan(ctx context.Context, name string) (context.Context, trace.Span) {
	ctx, span := tracer.Start(ctx, name, trace.WithSpanKind(trace.SpanKindServer))
	span.SetAttributes(
		attribute.String("service.name", name),
	)
	return ctx, span
}

func BeginSpanWithTraceID(ctx context.Context, corrID, parentID, name string) (context.Context, trace.Span) {
	traceID, err := trace.TraceIDFromHex(corrID)
	if err != nil {
		log.Printf("Failed to parse traceID: %v", err)
	}

	spanID, err := trace.SpanIDFromHex(parentID)
	if err != nil {
		log.Printf("Failed to parse spanID: %v", err)
	}

	spanContext := trace.NewSpanContext(trace.SpanContextConfig{
		TraceID:    traceID,
		SpanID:     spanID,
		TraceFlags: trace.FlagsSampled.WithSampled(true),
		Remote:     true,
	})

	ctx, span := tracer.Start(
		trace.ContextWithSpanContext(ctx, spanContext),
		name,
		trace.WithSpanKind(trace.SpanKindServer),
	)
	span.SetAttributes(
		attribute.String("TraceId", GetXrayTraceID(traceID.String())),
		attribute.String("ParentSpanId", parentID),
		attribute.KeyValue{
			Key:   semconv.ServiceNameKey,
			Value: attribute.StringValue(name),
		},
	)

	return ctx, span
}

func BeginSubSpan(ctx context.Context, name string) (context.Context, trace.Span) {
	return tracer.Start(ctx, name)
}

func BeginSubSpanWithNode(ctx context.Context, name string) (context.Context, trace.Span) {
	return tracer.Start(ctx, name, trace.WithSpanKind(trace.SpanKindServer))
}

func GetTraceSpanID(ctx context.Context) (traceID, spanID string) {
	spanContext := trace.SpanContextFromContext(ctx)
	if !spanContext.IsValid() {
		return traceID, spanID
	}

	traceID = spanContext.TraceID().String()
	spanID = spanContext.SpanID().String()
	return traceID, spanID
}

func GetXrayTraceID(traceID string) string {
	if traceID == "" {
		return ""
	}
	return fmt.Sprintf("1-%s-%s", traceID[0:8], traceID[8:])
}

func GetStatus(err error) (code codes.Code, msg string) {
	code = codes.Ok
	if err != nil {
		code = codes.Error
		msg = fmt.Sprintf("%v", err)
	}

	return
}

func InstrumentHandler(handlerFunc interface{}) interface{} {
	if tracerProvider == nil {
		return otellambda.InstrumentHandler(handlerFunc)
	} else {
		return otellambda.InstrumentHandler(handlerFunc, xrayconfig.WithRecommendedOptions(tracerProvider)...)
	}
}
