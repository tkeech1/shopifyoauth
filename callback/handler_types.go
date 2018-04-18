package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/tkeech1/goshopify"

	"github.com/aws/aws-lambda-go/events"
	"github.com/tkeech1/golambda_helper"
)

type LambdaHandler struct {
	Handler HandlerInterface
	Oauth   *golambda_helper.Oauth
}

type HandlerInterface interface {
	Getenv(string) string
	GetById(string, string, string, interface{}) error
	RequestToken(map[string]string, string, string) (string, error)
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

func (h *HandlerImpl) GetById(idName, idValue, tableName string, v interface{}) error {
	u := golambda_helper.DynamoHandler{
		Svc: &golambda_helper.DynamoHandler{},
	}
	return u.GetById(idName, idValue, tableName, v)
}

func (h *HandlerImpl) Put(i interface{}, s string) error {
	u := golambda_helper.DynamoHandler{
		Svc: &golambda_helper.DynamoHandler{},
	}
	return u.Put(i, s)
}

func (h *HandlerImpl) RequestToken(params map[string]string, secret string, apiKey string) (string, error) {
	u := goshopify.HttpRequestHandler{
		Req: &goshopify.HttpRequestHandler{},
	}
	return u.RequestToken(params, secret, apiKey)
}

func (h *LambdaHandler) HandleRequest(request events.APIGatewayProxyRequest) (golambda_helper.Response, error) {

	log.Print("Test")

	var shopname, installState string
	var ok bool

	if shopname, ok = request.QueryStringParameters["shop"]; !ok {
		return golambda_helper.GenerateError(errors.New("shop not found in callback"))
	}

	if installState, ok = request.QueryStringParameters["state"]; !ok {
		return golambda_helper.GenerateError(errors.New("state not found in callback"))
	}

	if !strings.HasSuffix(shopname, "myshopify.com") {
		return golambda_helper.GenerateError(errors.New("not a shopify url"))
	}

	err := h.Handler.GetById("shop_name", shopname, h.Handler.Getenv("TABLE_OAUTH"), h.Oauth)
	if err != nil {
		return golambda_helper.GenerateError(fmt.Errorf("error in GetById: %v", err))
	}

	if h.Oauth.InstallState != installState {
		return golambda_helper.GenerateError(errors.New("install state does not match"))
	}

	token, err := h.Handler.RequestToken(request.QueryStringParameters, h.Handler.Getenv("SHARED_SECRET"), h.Handler.Getenv("API_KEY"))
	if err != nil {
		return golambda_helper.GenerateError(fmt.Errorf("error in RequestToken: %v", err))
	}

	//fmt.Println(token)

	h.Oauth.OauthToken = token
	h.Oauth.CallbackDateTime = h.Handler.Now().Format(time.RFC3339)

	err = h.Handler.Put(h.Oauth, h.Handler.Getenv("TABLE_OAUTH"))
	if err != nil {
		return golambda_helper.GenerateError(fmt.Errorf("error in Put: %v", err))
	}

	return golambda_helper.GenerateRedirect(h.Handler.Getenv("SUCCESS_URI"))
}
