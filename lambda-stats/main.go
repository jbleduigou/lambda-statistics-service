package main

import (
	"context"
	"encoding/json"
	"net/http"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-lambda-go/lambdacontext"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func handler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	initLogger(ctx)
	zap.S().Debug("Received request")

	regions, err := getRequestedRegions(request)

	stats := []string{}

	for _, region := range regions {
		if err != nil {
			return events.APIGatewayProxyResponse{
				StatusCode: http.StatusUnprocessableEntity,
				Body:       err.Error(),
			}, nil
		}

		s, err := NewLambdaService(region)

		if err != nil {
			zap.S().Error("Error creating lambda service", zap.Error(err))
			return events.APIGatewayProxyResponse{}, err
		}

		lf, err := s.GetLambdaFunctions(ctx)

		if err != nil {
			zap.S().Errorw("Error while getting lambda statistics", "error", zap.Error(err), "region", region)
			return events.APIGatewayProxyResponse{}, err
		}

		stats = append(stats, lf...)
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
