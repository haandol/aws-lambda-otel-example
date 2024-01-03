package handler

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/haandol/aws-lambda-otel-example/demo/pkg/o11y"
	"github.com/haandol/aws-lambda-otel-example/demo/pkg/util"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	semconv "go.opentelemetry.io/otel/semconv/v1.10.0"
)

func HelloWorldController(c *gin.Context) {
	logger := util.GetLogger().With(
		"feature", "hello",
		"usecase", "hello_world",
		"component", "controller",
	)
	logger.Info("hello world")

	ctx, cancel := context.WithTimeout(c.Request.Context(), time.Second*10)
	defer cancel()

	ctx, span := o11y.BeginSpan(ctx, "helloworld controller")
	defer span.End()

	res, err := HelloWorldService(ctx)
	if err != nil {
		logger.Error("failed to say hello world", "err", err)
		span.RecordError(err)
		span.SetStatus(o11y.GetStatus(err))
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": "failed to say hello world",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": fmt.Sprintf("hello world with code: %v", res),
	})
}

func HelloWorldService(ctx context.Context) (int, error) {
	logger := util.GetLogger().With(
		"feature", "hello",
		"usecase", "hello_world",
		"component", "service",
	)
	logger.Info("hello world")

	ctx, span := o11y.BeginSubSpan(ctx, "helloworld service")
	defer span.End()

	url := "https://www.google.com"
	resp, err := makeRequest(ctx, url)
	span.SetAttributes(
		o11y.AttrInt(string(semconv.HTTPStatusCodeKey), resp.StatusCode),
	)
	if err != nil {
		logger.Error("failed to make request", "err", err)
		span.RecordError(err)
		span.SetStatus(o11y.GetStatus(err))
		return 0, err
	}

	logger.Info("response", "status", resp.StatusCode)

	return resp.StatusCode, nil
}

func makeRequest(ctx context.Context, url string) (*http.Response, error) {
	ctx, span := o11y.BeginSubSpan(ctx, "http request")
	defer span.End()

	span.SetAttributes(
		o11y.AttrString(string(semconv.HTTPURLKey), url),
		o11y.AttrString(string(semconv.HTTPMethodKey), http.MethodGet),
	)

	resp, err := otelhttp.Get(ctx, url)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(o11y.GetStatus(err))
		return nil, err
	}

	span.SetAttributes(
		o11y.AttrInt(string(semconv.HTTPStatusCodeKey), resp.StatusCode),
	)
	return resp, err
}
