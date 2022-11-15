package main

import (
	"context"
	"encoding/json"
	"lambda-stats/config"
	"lambda-stats/log"
	"lambda-stats/services"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"go.uber.org/zap"
)

func handler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	log.InitLogger(ctx)
	zap.S().Debug("Received request")

	stats, err := retrieveStatistics(ctx, config.SupportedRegions)
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

			s, err := services.NewLambdaService(region)

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
