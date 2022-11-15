package main

import (
	"context"
	"encoding/json"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"go.uber.org/zap"
	"lambda-stats/api"
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

	runtime, err := getRequestedRuntime(request)
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
	lf, err := s.GetLambdaFunctions(ctx)
	if err != nil {
		zap.S().Error("Error while retrieving list of lambda functions", zap.Error(err))
		return events.APIGatewayProxyResponse{}, err
	}

	stats := filterFunctionsByRuntime(lf, runtime)
	statsWithTags, err := filterFunctionsByTags(ctx, stats, s)
	if err != nil {
		zap.S().Error("Error while retrieving tags of lambda functions", zap.Error(err))
		return events.APIGatewayProxyResponse{}, err
	}
	payload, err := json.Marshal(statsWithTags)

	if err != nil {
		zap.S().Error("Error while serializing to json", zap.Error(err))
		return events.APIGatewayProxyResponse{}, err
	}

	return events.APIGatewayProxyResponse{
		Body:       string(payload),
		StatusCode: 200,
	}, nil
}

func filterFunctionsByRuntime(lf []api.LambdaFunction, runtime string) []api.LambdaFunction {
	stats := []api.LambdaFunction{}
	for _, v := range lf {
		if v.Runtime == runtime {
			stats = append(stats, v)
		}
	}
	return stats
}

func filterFunctionsByTags(ctx aws.Context, lambdas []api.LambdaFunction, s services.LambdaService) ([]api.LambdaFunction, error) {
	errChan := make(chan error)
	defer close(errChan)

	statChan := make(chan api.LambdaFunction)
	defer close(statChan)

	stats := []api.LambdaFunction{}
	var errOut error

	for _, l := range lambdas {
		go func(ctx aws.Context, lf api.LambdaFunction) {
			tags, err := s.GetTagsForFunction(ctx, lf.FunctionArn)
			errChan <- err
			statChan <- api.LambdaFunction{
				FunctionName: lf.FunctionName,
				FunctionArn:  lf.FunctionArn,
				Description:  lf.Description,
				Runtime:      lf.Runtime,
				Tags:         tags,
			}
		}(ctx, l)
	}

	for range lambdas {
		if err := <-errChan; err != nil {
			errOut = err
		}
	}
	for range lambdas {
		lf := <-statChan
		stats = append(stats, lf)
	}
	return stats, errOut
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

func getRequestedRuntime(request events.APIGatewayProxyRequest) (string, error) {
	runtime, found := request.QueryStringParameters["runtime"]
	if !found {
		return runtime, errors.ErrNoRuntime
	}
	for _, v := range config.SupportedRuntimes {
		if runtime == v {
			return runtime, nil
		}
	}
	return runtime, errors.ErrInvalidRuntime
}
