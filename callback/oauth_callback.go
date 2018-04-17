package main

import (
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/tkeech1/golambda_helper"
)

func Handler(request events.APIGatewayProxyRequest) (golambda_helper.Response, error) {
	h := LambdaHandler{
		Handler: &HandlerImpl{},
		Oauth:   golambda_helper.Oauth{},
	}
	return h.HandleRequest(request)
}

func main() {
	lambda.Start(Handler)
}
