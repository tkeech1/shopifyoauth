package main_test

import (
	"errors"
	"testing"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/stretchr/testify/assert"
	"github.com/tkeech1/golambda_helper"
	main "github.com/tkeech1/shopifyoauth/install"
	"github.com/tkeech1/shopifyoauth/install/mocks"
)

func TestHandlerShopifyOauth_OauthInstall(t *testing.T) {

	tests := map[string]struct {
		Request               events.APIGatewayProxyRequest
		Response              golambda_helper.Response
		GenerateStateResponse string
		GenerateStateError    error
		DynamoError           error
		time                  time.Time
		Oauth                 golambda_helper.Oauth
		EnvApiKey             string
		EnvScope              string
		EnvCallback           string
		EnvTableOauth         string
	}{
		"success_redirect": {
			Request: events.APIGatewayProxyRequest{
				PathParameters: map[string]string{"shopname": "testshop.myshopify.com"},
			},
			Response: golambda_helper.Response{
				StatusCode: 302,
				Header: golambda_helper.Header{
					Location: "https://testshop.myshopify.com/admin/oauth/authorize?client_id=SOMEKEY&redirect_uri=http%3A%2F%2Fmycallback.myshopify.com&scope=scope234",
				},
			},
			GenerateStateResponse: "FakeState",
			GenerateStateError:    nil,
			DynamoError:           nil,
			EnvApiKey:             "SOMEKEY",
			EnvScope:              "scope234",
			EnvCallback:           "http://mycallback.myshopify.com",
			EnvTableOauth:         "TABLE_OAUTH",
			time:                  time.Date(2009, 11, 17, 20, 34, 58, 651387237, time.UTC),
			Oauth: golambda_helper.Oauth{
				ShopName:        "testshop.myshopify.com",
				State:           "FakeState",
				InstallDateTime: time.Date(2009, 11, 17, 20, 34, 58, 651387237, time.UTC).Format(time.RFC3339),
			},
		},
		"error_GenerateState": {
			Request: events.APIGatewayProxyRequest{
				PathParameters: map[string]string{"shopname": "testshop.myshopify.com"},
			},
			Response: golambda_helper.Response{
				Body:       `{"message":"An error occurred processing the request."}`,
				StatusCode: 400,
				Header: golambda_helper.Header{
					ContentType:              "application/json",
					AccessControlAllowOrigin: "*",
				},
			},
			GenerateStateResponse: "",
			GenerateStateError:    errors.New("some error"),
			DynamoError:           nil,
			EnvApiKey:             "SOMEKEY",
			EnvScope:              "scope234",
			EnvCallback:           "http://mycallback.myshopify.com",
			EnvTableOauth:         "TABLE_OAUTH",
			time:                  time.Date(2009, 11, 17, 20, 34, 58, 651387237, time.UTC),
			Oauth: golambda_helper.Oauth{
				ShopName:        "testshop.myshopify.com",
				State:           "FAKESTATE",
				InstallDateTime: time.Date(2009, 11, 17, 20, 34, 58, 651387237, time.UTC).Format(time.RFC3339),
			},
		},
		"error_No_Shopname": {
			Request: events.APIGatewayProxyRequest{
				PathParameters: map[string]string{"dxx": "testshop.myshopify.com"},
			},
			Response: golambda_helper.Response{
				Body:       `{"message":"An error occurred processing the request."}`,
				StatusCode: 400,
				Header: golambda_helper.Header{
					ContentType:              "application/json",
					AccessControlAllowOrigin: "*",
				},
			},
			GenerateStateResponse: "FakeState",
			GenerateStateError:    nil,
			DynamoError:           nil,
			EnvApiKey:             "SOMEKEY",
			EnvScope:              "scope234",
			EnvCallback:           "http://mycallback.myshopify.com",
			EnvTableOauth:         "TABLE_OAUTH",
			time:                  time.Date(2009, 11, 17, 20, 34, 58, 651387237, time.UTC),
			Oauth: golambda_helper.Oauth{
				ShopName:        "testshop.myshopify.com",
				State:           "FakeState",
				InstallDateTime: time.Date(2009, 11, 17, 20, 34, 58, 651387237, time.UTC).Format(time.RFC3339),
			},
		},
		"error_Put": {
			Request: events.APIGatewayProxyRequest{
				PathParameters: map[string]string{"shopname": "testshop.myshopify.com"},
			},
			Response: golambda_helper.Response{
				Body:       `{"message":"An error occurred processing the request."}`,
				StatusCode: 400,
				Header: golambda_helper.Header{
					ContentType:              "application/json",
					AccessControlAllowOrigin: "*",
				},
			},
			GenerateStateResponse: "FakeState",
			GenerateStateError:    nil,
			DynamoError:           errors.New("some error"),
			EnvApiKey:             "SOMEKEY",
			EnvScope:              "scope234",
			EnvCallback:           "http://mycallback.myshopify.com",
			EnvTableOauth:         "TABLE_OAUTH",
			time:                  time.Date(2009, 11, 17, 20, 34, 58, 651387237, time.UTC),
			Oauth: golambda_helper.Oauth{
				ShopName:        "testshop.myshopify.com",
				State:           "FakeState",
				InstallDateTime: time.Date(2009, 11, 17, 20, 34, 58, 651387237, time.UTC).Format(time.RFC3339),
			},
		},
	}

	for name, test := range tests {
		t.Logf("Running test case: %s", name)

		handlerInt := &mocks.HandlerInterface{}
		handlerInt.
			On("Getenv", "API_KEY").
			Return(test.EnvApiKey)
		handlerInt.
			On("Getenv", "OAUTH_CALLBACK_URI").
			Return(test.EnvCallback)
		handlerInt.
			On("Getenv", "SCOPES").
			Return(test.EnvScope)
		handlerInt.
			On("Getenv", "TABLE_OAUTH").
			Return(test.EnvTableOauth)
		handlerInt.
			On("GenerateState").
			Return(test.GenerateStateResponse, test.GenerateStateError)
		handlerInt.
			On("Put", test.Oauth, "TABLE_OAUTH").
			Return(test.DynamoError)
		handlerInt.
			On("Now").
			Return(test.time)

		h := main.LambdaHandler{
			Handler: handlerInt,
		}
		response, _ := h.HandleRequest(test.Request)
		assert.Equal(t, test.Response.StatusCode, response.StatusCode)
		// exclude checking the state since it changes with every incovation
		if response.StatusCode == 302 {
			assert.Equal(t, test.Response.Header.Location, response.Header.Location[:136])
		}
	}
}
