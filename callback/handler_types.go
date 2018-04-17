package main

import (
	"errors"
	"net/http"
	"os"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/tkeech1/golambda_helper"
)

type LambdaHandler struct {
	Handler HandlerInterface
	Oauth   golambda_helper.Oauth
}

type HandlerInterface interface {
	Getenv(string) string
	GetById(string, string, string, interface{}) error
	Get(string) (*http.Response, error)
}

type HandlerImpl struct{}

func (h *HandlerImpl) Getenv(variableName string) string {
	return os.Getenv(variableName)
}

func (h *HandlerImpl) GetById(idName, idValue, tableName string, v interface{}) error {
	u := golambda_helper.DynamoHandler{
		Svc: &golambda_helper.DynamoHandler{},
	}
	return u.GetById(idName, idValue, tableName, v)
}

func (h *HandlerImpl) Get(s string) (*http.Response, error) {
	return http.Get(s)
}

func (h *LambdaHandler) HandleRequest(request events.APIGatewayProxyRequest) (golambda_helper.Response, error) {

	errorMessage := "An error occurred processing the request."

	var shopname, installState string
	var ok bool

	if shopname, ok = request.QueryStringParameters["shop"]; !ok {
		return golambda_helper.GenerateError(errors.New(errorMessage))
	}

	if installState, ok = request.QueryStringParameters["state"]; !ok {
		return golambda_helper.GenerateError(errors.New(errorMessage))
	}

	if !strings.HasSuffix(shopname, "myshopify.com") {
		return golambda_helper.GenerateError(errors.New(errorMessage))
	}

	//fmt.Println("shop:", shopname)
	//fmt.Println("state:", installState)

	err := h.Handler.GetById("shop_name", shopname, h.Handler.Getenv("TABLE_OAUTH"), &h.Oauth)
	if err != nil {
		return golambda_helper.GenerateError(errors.New(errorMessage))
	}

	if h.Oauth.InstallState != installState {
		return golambda_helper.GenerateError(errors.New(errorMessage))
	}

	// RequestToken
	// Put Oauth to DB

	//fmt.Println(oauth.ShopName)
	return golambda_helper.GenerateRedirect(h.Handler.Getenv("SUCCESS_URI"))
}
