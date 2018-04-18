package main_test

import (
	"errors"
	"testing"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/stretchr/testify/assert"
	"github.com/tkeech1/golambda_helper"
	main "github.com/tkeech1/shopifyoauth/callback"
	"github.com/tkeech1/shopifyoauth/callback/mocks"
)

func TestHandlerShopifyOauth_OauthCallback(t *testing.T) {

	tests := map[string]struct {
		Request            events.APIGatewayProxyRequest
		Response           golambda_helper.Response
		ResponseError      error
		EnvApiKey          string
		EnvSecret          string
		EnvScope           string
		EnvCallback        string
		RedirectUrl        string
		OauthTable         string
		GetByIdError       error
		InstallState       string
		TokenResponse      string
		TokenResponseError error
		Time               time.Time
		PutError           error
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
			RedirectUrl:        "SomeTestURL",
			ResponseError:      nil,
			GetByIdError:       nil,
			OauthTable:         "someTable",
			InstallState:       "12345",
			EnvApiKey:          "API_KEY",
			EnvSecret:          "SHARED_SECRET",
			TokenResponse:      "MYSHOPIFYTOKEN",
			TokenResponseError: nil,
			Time:               time.Date(2009, 11, 17, 20, 34, 58, 651387237, time.UTC),
			PutError:           nil,
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
				Body:       `{"message":"shop not found in callback"}`,
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
				Body:       `{"message":"state not found in callback"}`,
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
				Body:       `{"message":"not a shopify url"}`,
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
				Body:       `{"message":"error in GetById: An Error"}`,
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
				Body:       `{"message":"install state does not match"}`,
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
			InstallState:  "123456",
		},
		"error_token_response_error": {
			Request: events.APIGatewayProxyRequest{
				QueryStringParameters: map[string]string{
					"shop":  "testshop.myshopify.com",
					"state": "12345",
				},
			},
			Response: golambda_helper.Response{
				Body:       `{"message":"error in RequestToken: An Error"}`,
				StatusCode: 400,
				Header: golambda_helper.Header{
					ContentType:              "application/json",
					AccessControlAllowOrigin: "*",
				},
			},
			RedirectUrl:        "SomeTestURL",
			ResponseError:      nil,
			GetByIdError:       nil,
			OauthTable:         "someTable",
			InstallState:       "12345",
			EnvApiKey:          "API_KEY",
			EnvSecret:          "SHARED_SECRET",
			TokenResponse:      "",
			TokenResponseError: errors.New("An Error"),
		},
		"error_put_error": {
			Request: events.APIGatewayProxyRequest{
				QueryStringParameters: map[string]string{
					"shop":  "testshop.myshopify.com",
					"state": "12345",
				},
			},
			Response: golambda_helper.Response{
				Body:       `{"message":"error in Put: An Error"}`,
				StatusCode: 400,
				Header: golambda_helper.Header{
					ContentType:              "application/json",
					AccessControlAllowOrigin: "*",
				},
			},
			RedirectUrl:        "SomeTestURL",
			ResponseError:      nil,
			GetByIdError:       nil,
			OauthTable:         "someTable",
			InstallState:       "12345",
			EnvApiKey:          "API_KEY",
			EnvSecret:          "SHARED_SECRET",
			TokenResponse:      "MYSHOPIFYTOKEN",
			TokenResponseError: nil,
			Time:               time.Date(2009, 11, 17, 20, 34, 58, 651387237, time.UTC),
			PutError:           errors.New("An Error"),
		},
	}

	for name, test := range tests {
		t.Logf("Running test case: %s", name)

		handlerInt := &mocks.HandlerInterface{}
		oauth := &golambda_helper.Oauth{
			InstallState: test.InstallState,
		}
		handlerInt.
			On("Getenv", "SUCCESS_URI").
			Return(test.RedirectUrl)
		handlerInt.
			On("Getenv", "TABLE_OAUTH").
			Return(test.OauthTable)
		handlerInt.
			On("Getenv", "API_KEY").
			Return(test.EnvApiKey)
		handlerInt.
			On("Getenv", "SHARED_SECRET").
			Return(test.EnvSecret)
		handlerInt.
			On("RequestToken", test.Request.QueryStringParameters, test.EnvSecret, test.EnvApiKey).
			Return(test.TokenResponse, test.TokenResponseError)
		handlerInt.
			On("GetById", "shop_name", test.Request.QueryStringParameters["shop"], test.OauthTable, oauth).
			Return(test.GetByIdError)
		handlerInt.
			On("Now").
			Return(test.Time)
		handlerInt.
			On("Put", oauth, test.OauthTable).
			Return(test.PutError)

		h := main.LambdaHandler{
			Handler: handlerInt,
			Oauth:   oauth,
		}

		response, err := h.HandleRequest(test.Request)
		assert.Equal(t, test.Response, response)
		assert.Equal(t, test.ResponseError, err)
	}
}
