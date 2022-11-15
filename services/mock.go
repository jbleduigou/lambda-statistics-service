package services

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/service/lambda"
	"github.com/aws/aws-sdk-go/service/lambda/lambdaiface"
	"github.com/stretchr/testify/mock"
)

type mockedLambda struct {
	lambdaiface.LambdaAPI
	mock.Mock
}

func (m *mockedLambda) ListFunctionsWithContext(ctx aws.Context, input *lambda.ListFunctionsInput, opts ...request.Option) (*lambda.ListFunctionsOutput, error) {
	ret := m.Called(ctx, input, opts)
	if ret.Get(1) != nil {
		return nil, ret.Get(1).(error)
	}
	return ret.Get(0).(*lambda.ListFunctionsOutput), nil
}
