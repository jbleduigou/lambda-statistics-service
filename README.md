# Lambda Statistics Service
![Go](https://github.com/jbleduigou/lambda-statistics-service/workflows/Go/badge.svg)
![Lint](https://github.com/jbleduigou/lambda-statistics-service/workflows/Linting/badge.svg)
![SAM Deploy](https://github.com/jbleduigou/lambda-statistics-service/workflows/SAM/badge.svg)
[![License](https://img.shields.io/badge/License-Apache_2.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)

## Overview

This repo is a service exposing endpoints for retrieving statistics about lambda functions deployed in an AWS account.  
Two different endpoints are exposed:
* One for listing all AWS Lambda functions created in the AWS account across regions.
* One for searching for Lambda functions using some query params.

## Getting Started

This project needs the following tools installed:
* Go 1.19 : `brew install go`
* SAM cli : `brew tap aws/tap && brew install aws-sam-cli`

## Deploying the Service

You will need an AWS account and your AWS credentials should be already configured.  
Two steps are then required to deploy the service.

```sam build```  
This will build the CloudFormation template and compile the code.


```sam deploy --guided```  
This will guide you towards deploying the service in the AWS cloud.

## API Key

The service is implemented to be protected by an API key.  
This key can be retrieved for the AWS console.  
The other option is to use the AWS CLI for that and export it as an environment variable.  
This can be achieved through the following commands:
```bash
export API_ID=$(aws apigateway get-rest-apis --profile cicd --region us-east-1 --query "items[?name == 'dev-lambda-stats-api']" | jq '.[].id' | tr -d "\"")
export API_KEY=$( aws apigateway get-api-keys --profile cicd --region us-east-1 --query "items[?stageKeys && contains(stageKeys, '$API_ID/dev')]" --include-value | jq '.[].value' | tr -d "\"")
```

## List Functions Endpoint

The list functions endpoint has no query parameters.  
The cURL tool can be used to access it:
```bash
curl -H "x-api-key: $API_KEY" "https://$API_ID.execute-api.us-east-1.amazonaws.com/dev/list"
```
The output will be something similar to:
```json
[
  {
    "function-name": "hello-world-python",
    "function-arn": "arn:aws:lambda:us-east-1:123456789012:function:hello-world-python",
    "description": "A starter AWS Lambda function.",
    "runtime": "python3.7"
  },
  {
    "function-name": "hello-world-node",
    "function-arn": "arn:aws:lambda:us-east-1:123456789012:function:hello-world-node",
    "description": "A starter AWS Lambda function.",
    "runtime": "nodejs14.x"
  },
  {
    "function-name": "dev-list-lambda-function",
    "function-arn": "arn:aws:lambda:us-east-1:123456789012:function:dev-list-lambda-function",
    "description": "",
    "runtime": "go1.x"
  },
  {
    "function-name": "dev-search-lambda-function",
    "function-arn": "arn:aws:lambda:us-east-1:123456789012:function:dev-search-lambda-function",
    "description": "",
    "runtime": "go1.x"
  },
  {
    "function-name": "sam-app-HelloWorldFunction-MiteidKyOkBk",
    "function-arn": "arn:aws:lambda:eu-west-1:123456789012:function:sam-app-HelloWorldFunction-MiteidKyOkBk",
    "description": "",
    "runtime": "go1.x"
  }
]
```

## Search Functions Endpoint

The search functions endpoint has 3 mandatory query parameters:  
    1. `runtime`: filter functions with a given runtime environment - [Read more about Lambda runtimes](https://docs.aws.amazon.com/lambda/latest/dg/lambda-runtimes.html)
    2. `tags`: filter functions with a given tag value or tag key - [Read more about tags on Lambda](https://docs.aws.amazon.com/lambda/latest/dg/configuration-tags.html)  
    3. `region`: filter functions running in a given AWS region - [Read more about list of AWS regions where AWS Lambda is supported](https://docs.aws.amazon.com/general/latest/gr/lambda-service.html)

The cURL tool can be used to access it:
```bash
curl -H "x-api-key: $API_KEY" "https://$API_ID.execute-api.us-east-1.amazonaws.com/dev/search?region=us-east-1&runtime=go1.x&tags=stats"
```

The output will be something similar to:
```json

[
  {
    "function-name": "dev-list-lambda-function",
    "function-arn": "arn:aws:lambda:us-east-1:123456789012:function:dev-list-lambda-function",
    "description": "",
    "runtime": "go1.x",
    "tags": {
      "Environment": "dev",
      "Service": "list-functions",
      "aws:cloudformation:logical-id": "ListFunctionsFunction",
      "aws:cloudformation:stack-id": "arn:aws:cloudformation:us-east-1:123456789012:stack/lambda-stats/11973c00-1234-11ed-aa46-0e70963405b0",
      "aws:cloudformation:stack-name": "lambda-stats",
      "lambda:createdBy": "SAM"
    }
  },
  {
    "function-name": "dev-search-lambda-function",
    "function-arn": "arn:aws:lambda:us-east-1:123456789012:function:dev-search-lambda-function",
    "description": "",
    "runtime": "go1.x",
    "tags": {
      "Environment": "dev",
      "Service": "search-functions",
      "aws:cloudformation:logical-id": "SearchFunctionsFunction",
      "aws:cloudformation:stack-id": "arn:aws:cloudformation:us-east-1:123456789012:stack/lambda-stats/11973c00-1234-11ed-aa46-0e70963405b0",
      "aws:cloudformation:stack-name": "lambda-stats",
      "lambda:createdBy": "SAM"
    }
  }
]
```

## Improvements

While this service is fully functionning, there are a couple of things which could be improved:
* 🤝 Having an OpenAPI definition would be nice
* ⚡️ Review performances of the filtering (successive loops)
* ♻️ Use enums for the constants in the config package
* ✅ Add more unit tests
* 📦️ In the search endpoint response, group lambda functions by region
* ⏰ Have an automated job for updating constants when new runtime or regions are adde by AWS
* 🍸 Use a web-framework such as Go Gin