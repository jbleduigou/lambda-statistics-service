package main

import (
	"context"
	"encoding/json"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"go.uber.org/zap"
	"lambda-stats/config"
	"lambda-stats/errors"
	"lambda-stats/log"
	"lambda-stats/services"
	"net/http"
)

func handler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	log.InitLogger(ctx)
	zap.S().Debug("Received request")

	region, err := getRequestedRegion(request)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusUnprocessableEntity,
			Body:       err.Error(),
		}, nil
	}

	s, err := services.NewLambdaService(region)
	if err != nil {
		zap.S().Error("Error while creating service", zap.Error(err))
		return events.APIGatewayProxyResponse{}, err
	}
	stats, err := s.GetLambdaFunctions(ctx)
	if err != nil {
		zap.S().Error("Error while retrieving list of lambda functions", zap.Error(err))
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

func getRequestedRegion(request events.APIGatewayProxyRequest) (string, error) {
	region, found := request.QueryStringParameters["region"]
	if !found {
		return region, errors.ErrNoRegion
	}
	for _, v := range config.SupportedRegions {
		if region == v {
			return region, nil
		}
	}
	return region, errors.ErrInvalidRegion
}
