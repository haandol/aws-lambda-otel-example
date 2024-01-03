package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	ginadapter "github.com/awslabs/aws-lambda-go-api-proxy/gin"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/haandol/aws-lambda-otel-example/demo/internal/constant"
	"github.com/haandol/aws-lambda-otel-example/demo/internal/feature/hello/handler"
	"github.com/haandol/aws-lambda-otel-example/demo/pkg/o11y"
	"github.com/haandol/aws-lambda-otel-example/demo/pkg/util"
)

var (
	r         *gin.Engine
	ginLambda *ginadapter.GinLambdaV2
	isProd    bool = true
)

func init() {
	// setup logger
	logger := util.InitLogger(isProd)
	logger.Info("initializing...")

	// setup o11y
	o11y.InitOtel()

	// setup router
	r = gin.Default()
	r.Use(util.GinSlogWithConfig(logger, &util.Config{
		UTC: false,
	}))
	r.Use(cors.New(cors.Config{
		AllowOrigins:  []string{"*"},
		AllowMethods:  []string{"*"},
		AllowHeaders:  []string{"*"},
		ExposeHeaders: []string{"Content-Length"},
		MaxAge:        12 * time.Hour,
	}))
	r.Use(util.RecoveryWithSlog(logger, true))

	r.GET("/", handler.HelloWorldController)

	// setup ginLambda
	ginLambda = ginadapter.NewV2(r)
}

func LambdaHandler(ctx context.Context, req events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	return ginLambda.ProxyWithContext(ctx, req)
}

func main() {
	logger := util.GetLogger()

	lambda.Start(o11y.InstrumentHandler(LambdaHandler))

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit

	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		defer cancel()
		logger.Info("Closing o11y connection...")
		if err := o11y.Close(ctx); err != nil {
			logger.Error("error on closing o11y", err)
		} else {
			logger.Info("o11y connection closed.")
		}
	}()
	select {
	case <-ctx.Done():
		logger.Info("Graceful close complete")
	case <-time.After(constant.GracefulShutdownTimeout):
		logger.Error("closed by timeout")
	}
}
