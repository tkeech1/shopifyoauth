package main

import (
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/tkeech1/golambda_helper"
)

func Handler(request events.APIGatewayProxyRequest) (golambda_helper.Response, error) {
	envHandler := &EnvHandler{}
	h := EnvHandler{
		Env: envHandler,
	}
	return h.Handler(request)
}

func main() {
	lambda.Start(Handler)
}
