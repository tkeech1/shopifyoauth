package main_test

import (
	"errors"
	"testing"

	"github.com/aws/aws-lambda-go/events"
	"github.com/stretchr/testify/assert"
	"github.com/tkeech1/golambda_helper"
	main "github.com/tkeech1/shopifyoauth/callback"
	"github.com/tkeech1/shopifyoauth/callback/mocks"
)

func TestHandlerShopifyOauth_OauthCallback(t *testing.T) {

	tests := map[string]struct {
		Request       events.APIGatewayProxyRequest
		Response      golambda_helper.Response
		ResponseError error
		EnvApiKey     string
		EnvScope      string
		EnvCallback   string
		RedirectUrl   string
		OauthTable    string
		GetByIdError  error
		InstallState  string
	}{
		"success": {
			Request: events.APIGatewayProxyRequest{
				QueryStringParameters: map[string]string{
					"shop":  "testshop.myshopify.com",
					"state": "12345",
				},
			},
			Response: golambda_helper.Response{
				StatusCode: 302,
				Header: golambda_helper.Header{
					Location: "SomeTestURL",
				},
			},
			RedirectUrl:   "SomeTestURL",
			ResponseError: nil,
			GetByIdError:  nil,
			OauthTable:    "someTable",
			InstallState:  "12345",
		},
		"error_no_shop": {
			Request: events.APIGatewayProxyRequest{
				QueryStringParameters: map[string]string{
					"shopxx": "testshop.myshopify.com",
					"state":  "12345",
				},
			},
			RedirectUrl: "SomeTestURL",
			Response: golambda_helper.Response{
				Body:       `{"message":"An error occurred processing the request."}`,
				StatusCode: 400,
				Header: golambda_helper.Header{
					ContentType:              "application/json",
					AccessControlAllowOrigin: "*",
				},
			},
			ResponseError: nil,
			GetByIdError:  nil,
			OauthTable:    "someTable",
			InstallState:  "12345",
		},
		"error_no_state": {
			Request: events.APIGatewayProxyRequest{
				QueryStringParameters: map[string]string{
					"shop":   "testshop.myshopify.com",
					"statex": "12345",
				},
			},
			RedirectUrl: "SomeTestURL",
			Response: golambda_helper.Response{
				Body:       `{"message":"An error occurred processing the request."}`,
				StatusCode: 400,
				Header: golambda_helper.Header{
					ContentType:              "application/json",
					AccessControlAllowOrigin: "*",
				},
			},
			ResponseError: nil,
			GetByIdError:  nil,
			OauthTable:    "someTable",
			InstallState:  "12345",
		},
		"error_not_a_shopify_url": {
			Request: events.APIGatewayProxyRequest{
				QueryStringParameters: map[string]string{
					"shop":  "testshop.myshopifyx.com",
					"state": "12345",
				},
			},
			Response: golambda_helper.Response{
				Body:       `{"message":"An error occurred processing the request."}`,
				StatusCode: 400,
				Header: golambda_helper.Header{
					ContentType:              "application/json",
					AccessControlAllowOrigin: "*",
				},
			},
			RedirectUrl:   "SomeTestURL",
			ResponseError: nil,
			GetByIdError:  nil,
			OauthTable:    "someTable",
			InstallState:  "12345",
		},
		"error_getbyid_error": {
			Request: events.APIGatewayProxyRequest{
				QueryStringParameters: map[string]string{
					"shop":  "testshop.myshopify.com",
					"state": "12345",
				},
			},
			Response: golambda_helper.Response{
				Body:       `{"message":"An error occurred processing the request."}`,
				StatusCode: 400,
				Header: golambda_helper.Header{
					ContentType:              "application/json",
					AccessControlAllowOrigin: "*",
				},
			},
			RedirectUrl:   "SomeTestURL",
			ResponseError: nil,
			GetByIdError:  errors.New("An Error"),
			OauthTable:    "someTable",
			InstallState:  "12345",
		},
		"error_state_does_not_match": {
			Request: events.APIGatewayProxyRequest{
				QueryStringParameters: map[string]string{
					"shop":  "testshop.myshopify.com",
					"state": "12345",
				},
			},
			Response: golambda_helper.Response{
				Body:       `{"message":"An error occurred processing the request."}`,
				StatusCode: 400,
				Header: golambda_helper.Header{
					ContentType:              "application/json",
					AccessControlAllowOrigin: "*",
				},
			},
			RedirectUrl:   "SomeTestURL",
			ResponseError: nil,
			GetByIdError:  errors.New("An Error"),
			OauthTable:    "someTable",
			InstallState:  "123456",
		},
	}

	for name, test := range tests {
		t.Logf("Running test case: %s", name)

		handlerInt := &mocks.HandlerInterface{}
		oauth := golambda_helper.Oauth{
			InstallState: test.InstallState,
		}
		handlerInt.
			On("Getenv", "SUCCESS_URI").
			Return(test.RedirectUrl)
		handlerInt.
			On("Getenv", "TABLE_OAUTH").
			Return(test.OauthTable)
		handlerInt.
			On("GetById", "shop_name", test.Request.QueryStringParameters["shop"], test.OauthTable, &oauth).
			Return(test.GetByIdError)

		h := main.LambdaHandler{
			Handler: handlerInt,
			Oauth:   oauth,
		}

		response, err := h.HandleRequest(test.Request)
		assert.Equal(t, test.Response, response)
		assert.Equal(t, test.ResponseError, err)
	}
}
