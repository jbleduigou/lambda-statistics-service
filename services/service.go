package services

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/lambda"
	"github.com/aws/aws-sdk-go/service/lambda/lambdaiface"
)

type LambdaService interface {
	GetLambdaFunctions(ctx aws.Context) ([]string, error)
}

func NewLambdaService(region string) (LambdaService, error) {
	sess := session.Must(session.NewSession())
	l := lambda.New(sess, &aws.Config{Region: aws.String(region)})

	return &lambdaServiceImpl{l: l}, nil
}

type lambdaServiceImpl struct {
	l lambdaiface.LambdaAPI
}

func (s *lambdaServiceImpl) GetLambdaFunctions(ctx aws.Context) ([]string, error) {
	lambdas := []string{}
	input := &lambda.ListFunctionsInput{}
	output, err := s.l.ListFunctionsWithContext(ctx, input)
	if err != nil {
		return lambdas, err
	}
	for _, f := range output.Functions {
		lambdas = append(lambdas, *f.FunctionArn)
	}
	return lambdas, nil
}
