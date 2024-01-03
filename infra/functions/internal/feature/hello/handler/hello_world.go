package handler

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/haandol/aws-lambda-otel-example/demo/internal/feature/hello/instrument"
	"github.com/haandol/aws-lambda-otel-example/demo/pkg/o11y"
	"github.com/haandol/aws-lambda-otel-example/demo/pkg/util"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

func HelloWorldController(c *gin.Context) {
	logger := util.GetLogger().With(
		"feature", "hello",
		"usecase", "hello_world",
		"component", "controller",
	)

	ctx, cancel := context.WithTimeout(c.Request.Context(), time.Second*10)
	defer cancel()

	ctx, span := o11y.BeginSpan(ctx, "helloworld controller")
	defer span.End()

	res, err := HelloWorldService(ctx)
	if err != nil {
		msg := "failed to say hello world"
		instrument.RecordErrorWithMessage(ctx, logger, span, err, msg)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": msg,
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
	if err != nil {
		instrument.RecordErrorWithMessage(ctx, logger, span, err, "failed to make request")
		return 0, err
	}

	logger.Info("response", "status", resp.StatusCode)

	return resp.StatusCode, nil
}

func makeRequest(ctx context.Context, url string) (*http.Response, error) {
	logger := util.GetLogger().With(
		"feature", "hello",
		"usecase", "hello_world",
		"component", "adapter",
	)

	ctx, span := o11y.BeginSubSpan(ctx, "http request")
	defer span.End()

	instrument.HTTPGetReq(logger, span, http.MethodGet, url)

	resp, err := otelhttp.Get(ctx, url)
	if err != nil {
		instrument.RecordErrorWithMessage(ctx, logger, span, err, "failed to make request")
		return nil, err
	}

	instrument.HTTPResponseCode(ctx, span, resp.StatusCode)

	return resp, err
}
