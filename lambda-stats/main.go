package main

import (
	"context"
	"encoding/json"
	"net/http"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-lambda-go/lambdacontext"
	"github.com/aws/aws-sdk-go/aws"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func handler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	initLogger(ctx)
	zap.S().Debug("Received request")

	regions, err := getRequestedRegions(request)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusUnprocessableEntity,
			Body:       err.Error(),
		}, nil
	}

	stats, err := retrieveStatistics(ctx, regions)
	if err != nil {
		zap.S().Error("Error while retrieving statistics", zap.Error(err))
		return events.APIGatewayProxyResponse{}, err
	}

	payload, err := json.Marshal(stats)

	if err != nil {
		zap.S().Error("Error while serializing to json", zap.Error(err))
		return events.APIGatewayProxyResponse{}, err
	}

	return events.APIGatewayProxyResponse{
		Body:       string(payload),
		StatusCode: 200,
	}, nil
}

func main() {
	lambda.Start(handler)
}

func retrieveStatistics(ctx context.Context, regions []string) ([]string, error) {
	errChan := make(chan error)
	defer close(errChan)

	statChan := make(chan []string)
	defer close(statChan)

	stats := []string{}
	var errOut error

	for _, region := range regions {
		go func(ctx aws.Context, region string) {

			s, err := NewLambdaService(region)

			if err != nil {
				errChan <- err
				statChan <- []string{}
				return
			}

			lf, err := s.GetLambdaFunctions(ctx)
			errChan <- err
			statChan <- lf

		}(ctx, region)
	}

	for range regions {
		if err := <-errChan; err != nil {
			errOut = err
		}
	}
	for range regions {
		lf := <-statChan
		stats = append(stats, lf...)
	}
	return stats, errOut
}

func initLogger(ctx context.Context) {
	// Retrieve AWS Request ID
	lc, _ := lambdacontext.FromContext(ctx)
	requestID := lc.AwsRequestID
	cfg := zap.Config{
		Encoding:         "json",
		Level:            zap.NewAtomicLevelAt(getLogLevel()),
		OutputPaths:      []string{"stdout"},
		ErrorOutputPaths: []string{"stderr"},
		InitialFields:    map[string]interface{}{"request-id": requestID},
		EncoderConfig: zapcore.EncoderConfig{
			MessageKey: "message",

			LevelKey:    "level",
			EncodeLevel: zapcore.CapitalLevelEncoder,

			TimeKey:    "time",
			EncodeTime: zapcore.ISO8601TimeEncoder,

			CallerKey:    "caller",
			EncodeCaller: zapcore.ShortCallerEncoder,
		},
	}
	logger, _ := cfg.Build()
	zap.ReplaceGlobals(logger)
}

func getLogLevel() zapcore.Level {
	lvl, found := os.LookupEnv("LOG_LEVEL")
	if found {
		var l zapcore.Level
		if err := l.Set(lvl); err == nil {
			return l
		}
	}
	env := os.Getenv("ENVIRONMENT")
	if env == "dev" {
		return zap.DebugLevel
	}
	if env == "staging" {
		return zap.InfoLevel
	}
	return zap.WarnLevel
}

func getRequestedRegions(request events.APIGatewayProxyRequest) ([]string, error) {
	region, found := request.QueryStringParameters["region"]
	if !found {
		return regions, nil
	}
	for _, v := range regions {
		if region == v {
			return []string{region}, nil
		}
	}
	return []string{}, ErrInvalidRegion
}
