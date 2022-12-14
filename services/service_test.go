package services

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/stretchr/testify/mock"
	"lambda-stats/api"
	"testing"

	"github.com/aws/aws-sdk-go/service/lambda"
	"github.com/stretchr/testify/assert"
)

func TestNewLambdaServiceShouldConfigureRegion(t *testing.T) {
	svc, err := NewLambdaService("us-east-1")

	assert.Nil(t, err)

	assert.Equal(t, "us-east-1", *(svc.(*lambdaServiceImpl).l.(*lambda.Lambda).Client.Config.Region))
}

func TestGetLambdaFunctions(t *testing.T) {
	ctx := context.Background()

	for _, tc := range []struct {
		name          string
		setupMock     func(m *mockedLambda)
		wantFunctions []api.LambdaFunction
		wantErr       string
	}{
		{
			name: "should return error",
			setupMock: func(m *mockedLambda) {
				m.On("ListFunctionsWithContext", ctx, mock.Anything, mock.Anything).Return(nil, fmt.Errorf("fake error"))
			},
			wantFunctions: []api.LambdaFunction{},
			wantErr:       "fake error",
		},
		{
			name: "should return two functions",
			setupMock: func(m *mockedLambda) {
				m.On("ListFunctionsWithContext", ctx, mock.Anything, mock.Anything).
					Return(&lambda.ListFunctionsOutput{Functions: []*lambda.FunctionConfiguration{
						{FunctionArn: aws.String("arn:aws:lambda:us-east-1:123456789012:function:sam-app-HelloWorldFunction")},
						{FunctionArn: aws.String("arn:aws:lambda:us-east-1:123456789012:function:sam-app-AlexaFunction")},
					}},
						nil)
			},
			wantFunctions: []api.LambdaFunction{{FunctionName: "", FunctionArn: "arn:aws:lambda:us-east-1:123456789012:function:sam-app-HelloWorldFunction", Description: "", Runtime: "", Tags: map[string]*string(nil)}, {FunctionName: "", FunctionArn: "arn:aws:lambda:us-east-1:123456789012:function:sam-app-AlexaFunction", Description: "", Runtime: "", Tags: map[string]*string(nil)}},
		},
	} {
		m := &mockedLambda{}
		svc := &lambdaServiceImpl{l: m}
		tc.setupMock(m)

		f, err := svc.GetLambdaFunctions(ctx)

		if err != nil {
			assert.Equal(t, tc.wantErr, err.Error(), tc.name)
		}
		assert.Equal(t, tc.wantFunctions, f, tc.name)
	}
}
