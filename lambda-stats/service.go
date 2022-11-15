package main

import (
	"github.com/aws/aws-sdk-go/aws"
)

type LambdaService interface {
	GetLambdaFunctions(ctx aws.Context) ([]string, error)
}

func NewLambdaService(region string) (LambdaService, error) {
	return &lambdaServiceImpl{}, nil
}

type lambdaServiceImpl struct {
}

func (s *lambdaServiceImpl) GetLambdaFunctions(ctx aws.Context) ([]string, error) {
	return []string{"pouet"}, nil
}
