package main

import (
	"context"
	"encoding/json"
	"lambda-stats/config"
	"lambda-stats/errors"
	"lambda-stats/log"
	"lambda-stats/services"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"go.uber.org/zap"
)

func handler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	log.InitLogger(ctx)
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

func getRequestedRegions(request events.APIGatewayProxyRequest) ([]string, error) {
	region, found := request.QueryStringParameters["region"]
	if !found {
		return config.SupportedRegions, nil
	}
	for _, v := range config.SupportedRegions {
		if region == v {
			return []string{region}, nil
		}
	}
	return []string{}, errors.ErrInvalidRegion
}
