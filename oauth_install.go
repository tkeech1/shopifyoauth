package main

import (
	"errors"
	"os"

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

func (h *EnvHandler) Getenv(variableName string) string {
	return os.Getenv(variableName)
}

func (h *EnvHandler) Handler(request events.APIGatewayProxyRequest) (golambda_helper.Response, error) {
	if shopname, ok := request.PathParameters["shopname"]; ok {

		uuidHandler := &golambda_helper.UuidHandler{}
		u := golambda_helper.UuidHandler{
			Uuid: uuidHandler,
		}

		permissionUrl, err := u.Install(h.Env.Getenv("API_KEY"), h.Env.Getenv("SCOPES"), h.Env.Getenv("OAUTH_CALLBACK_URI"), shopname)
		if err == nil {
			return golambda_helper.GenerateRedirect(permissionUrl)
		}
	}

	return golambda_helper.GenerateError(errors.New("An error occurred processing the request."))
}
