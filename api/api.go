package api

type LambdaFunction struct {
	FunctionName string `json:"function-name"`
	FunctionArn  string `json:"function-arn"`
	Description  string `json:"description"`
	Runtime      string `json:"runtime"`
}
