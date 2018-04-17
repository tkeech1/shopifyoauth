package main

import (
	"errors"
	"os"
	"time"

	"github.com/tkeech1/goshopify"

	"github.com/aws/aws-lambda-go/events"
	"github.com/tkeech1/golambda_helper"
)

type LambdaHandler struct {
	Handler HandlerInterface
}

type HandlerInterface interface {
	Getenv(string) string
	GenerateState() (string, error)
	Put(interface{}, string) error
	Now() time.Time
}

type HandlerImpl struct{}

func (h *HandlerImpl) Getenv(variableName string) string {
	return os.Getenv(variableName)
}

func (h *HandlerImpl) Now() time.Time {
	return time.Now()
}

func (h *HandlerImpl) GenerateState() (string, error) {
	u := golambda_helper.UuidHandler{
		Uuid: &golambda_helper.UuidHandler{},
	}
	return u.GenerateState()
}

func (h *HandlerImpl) Put(i interface{}, s string) error {
	u := golambda_helper.DynamoHandler{
		Svc: &golambda_helper.DynamoHandler{},
	}
	return u.Put(i, s)
}

func (h *LambdaHandler) HandleRequest(request events.APIGatewayProxyRequest) (golambda_helper.Response, error) {
	if shopname, ok := request.PathParameters["shopname"]; ok {

		installState, err := h.Handler.GenerateState()
		if err != nil {
			return golambda_helper.GenerateError(errors.New("An error occurred processing the request."))
		}

		permissionUrl := goshopify.CreatePermissionUrl(h.Handler.Getenv("API_KEY"), h.Handler.Getenv("SCOPES"), h.Handler.Getenv("OAUTH_CALLBACK_URI"), installState, shopname)

		oauth := golambda_helper.Oauth{
			ShopName:        shopname,
			InstallState:    installState,
			InstallDateTime: h.Handler.Now().Format(time.RFC3339),
		}

		err = h.Handler.Put(oauth, h.Handler.Getenv("TABLE_OAUTH"))
		if err != nil {
			return golambda_helper.GenerateError(errors.New("An error occurred processing the request."))
		}

		return golambda_helper.GenerateRedirect(permissionUrl)

	}

	return golambda_helper.GenerateError(errors.New("An error occurred processing the request."))
}
