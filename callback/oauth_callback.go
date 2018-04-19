package main

import (
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/tkeech1/golambda_helper"
	"github.com/tkeech1/goshopify"
)

func Handler(request events.APIGatewayProxyRequest) (golambda_helper.Response, error) {
	h := LambdaHandler{
		Handler: &HandlerImpl{},
		Oauth:   &goshopify.Oauth{},
	}
	return h.HandleRequest(request)
}

func main() {
	lambda.Start(Handler)
}
