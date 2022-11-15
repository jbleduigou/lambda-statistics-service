package services

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/lambda"
	"github.com/aws/aws-sdk-go/service/lambda/lambdaiface"
	"lambda-stats/api"
)

type LambdaService interface {
	GetLambdaFunctions(ctx aws.Context) ([]api.LambdaFunction, error)
	GetTagsForFunction(ctx aws.Context, arn string) (map[string]*string, error)
}

func NewLambdaService(region string) (LambdaService, error) {
	sess := session.Must(session.NewSession())
	l := lambda.New(sess, &aws.Config{Region: aws.String(region)})

	return &lambdaServiceImpl{l: l}, nil
}

type lambdaServiceImpl struct {
	l lambdaiface.LambdaAPI
}

func (s *lambdaServiceImpl) GetLambdaFunctions(ctx aws.Context) ([]api.LambdaFunction, error) {
	lambdas := []api.LambdaFunction{}
	input := &lambda.ListFunctionsInput{}
	output, err := s.l.ListFunctionsWithContext(ctx, input)
	if err != nil {
		return lambdas, err
	}
	for _, fc := range output.Functions {
		f := api.LambdaFunction{
			FunctionName: aws.StringValue(fc.FunctionName),
			FunctionArn:  aws.StringValue(fc.FunctionArn),
			Description:  aws.StringValue(fc.Description),
			Runtime:      aws.StringValue(fc.Runtime),
		}
		lambdas = append(lambdas, f)
	}
	return lambdas, nil
}

func (s *lambdaServiceImpl) GetTagsForFunction(ctx aws.Context, arn string) (map[string]*string, error) {
	input := &lambda.ListTagsInput{Resource: aws.String(arn)}
	output, err := s.l.ListTagsWithContext(ctx, input)
	if err != nil {
		return map[string]*string{}, err
	}
	return output.Tags, nil
}
