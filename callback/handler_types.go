package main

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/tkeech1/golambda_helper"
)

type EnvInterface interface {
	Getenv(string) string
}

type EnvHandler struct {
	Env EnvInterface
}

func (h *EnvHandler) Getenv(variableName string) string {
	return os.Getenv(variableName)
}

func (h *EnvHandler) Handler(request events.APIGatewayProxyRequest) (golambda_helper.Response, error) {

	errorMessage := "An error occurred processing the request."

	if shopname, ok := request.QueryStringParameters["shop"]; ok {

		if state, ok := request.PathParameters["state"]; ok {

			if !strings.HasSuffix(shopname, "myshopify.com") {
				return golambda_helper.GenerateError(errors.New(errorMessage))
			}

			var ShopName golambda_helper.ShopName
			dynamoHandler := &golambda_helper.DynamoHandler{}
			m := golambda_helper.DynamoHandler{
				Svc: dynamoHandler,
			}
			err := m.GetById("id", "30174253-3caf-48e1-96d5-6595b783aa62", "goshopname-dev", &ShopName)
			if err != nil {
				return golambda_helper.GenerateError(err)
			}

			fmt.Println("shop:", shopname)
			fmt.Println("state:", state)

			/*if shop.state != state {
				return golambda_helper.GenerateError(errors.New(errorMessage))
			}

			if state == getOauthByShopName(shopname); ok {

			}*/

			/*uuidHandler := &golambda_helper.UuidHandler{}
			u := golambda_helper.UuidHandler{
				Uuid: uuidHandler,
			}

			permissionUrl, err := u.Install(h.Env.Getenv("API_KEY"), h.Env.Getenv("SCOPES"), h.Env.Getenv("OAUTH_CALLBACK_URI"), shopname)
			if err == nil {
				return golambda_helper.GenerateRedirect(permissionUrl)
			}*/
		}
		fmt.Println(shopname)
	}

	return golambda_helper.GenerateError(errors.New(errorMessage))
}
