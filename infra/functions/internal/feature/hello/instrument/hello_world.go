package instrument

import (
	"context"
	"log/slog"

	"github.com/haandol/aws-lambda-otel-example/demo/pkg/o11y"
	"github.com/pkg/errors"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	semconv "go.opentelemetry.io/otel/semconv/v1.10.0"
	"go.opentelemetry.io/otel/trace"
)

func HTTPGetReq(logger *slog.Logger, span trace.Span, method string, url string) {
	span.SetAttributes(
		attribute.KeyValue{
			Key:   semconv.HTTPURLKey,
			Value: attribute.StringValue(url),
		},
		attribute.KeyValue{
			Key:   semconv.HTTPMethodKey,
			Value: attribute.StringValue(method),
		},
	)
}

func HTTPResponseCode(ctx context.Context, span trace.Span, code int) {
	span.SetAttributes(
		attribute.KeyValue{
			Key:   semconv.HTTPStatusCodeKey,
			Value: attribute.IntValue(code),
		},
	)
	if code >= 200 && code < 400 {
		increaseHttpSuccessCount(ctx)
	}
}

func RecordErrorWithMessage(ctx context.Context, logger *slog.Logger, span trace.Span, err error, msg string) {
	logger.Error(msg, "err", err)
	span.RecordError(errors.Wrap(err, msg))
	span.SetStatus(codes.Error, msg)
	increaseHttpErrorCount(ctx)
}

func increaseHttpErrorCount(ctx context.Context) {
	counter, _ := o11y.NewIntCounter("http_error_count")
	counter.Add(ctx, 1)
}

func increaseHttpSuccessCount(ctx context.Context) {
	counter, _ := o11y.NewIntCounter("http_success_count")
	counter.Add(ctx, 1)
}
