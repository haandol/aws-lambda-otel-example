package instrument

import (
	"log/slog"

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

func HTTPResponseCode(span trace.Span, code int) {
	span.SetAttributes(
		attribute.KeyValue{
			Key:   semconv.HTTPStatusCodeKey,
			Value: attribute.IntValue(code),
		},
	)
}

func RecordError(logger *slog.Logger, span trace.Span, err error, msg string) {
	logger.Error(msg, "err", err)
	span.RecordError(errors.Wrap(err, msg))
	span.SetStatus(codes.Error, msg)
}
