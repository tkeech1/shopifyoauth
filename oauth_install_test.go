package main_test

import (
	"testing"

	"github.com/aws/aws-lambda-go/events"
	"github.com/stretchr/testify/assert"
	"github.com/tkeech1/golambda_helper"
	main "github.com/tkeech1/shopifyoauth"
	"github.com/tkeech1/shopifyoauth/mocks"
)

func TestHandlerShopifyOauth_OauthInstall(t *testing.T) {

	tests := map[string]struct {
		Request       events.APIGatewayProxyRequest
		Response      golambda_helper.Response
		ResponseError error
		EnvApiKey     string
		EnvScope      string
		EnvCallback   string
	}{
		"redirect": {
			Request: events.APIGatewayProxyRequest{
				PathParameters: map[string]string{"shopname": "testshop.myshopify.com"},
			},
			Response: golambda_helper.Response{
				StatusCode: 302,
				Header: golambda_helper.Header{
					Location: "https://testshop.myshopify.com/admin/oauth/authorize?client_id=SOMEKEY&redirect_uri=http%3A%2F%2Fmycallback.myshopify.com&scope=scope234",
				},
			},
			ResponseError: nil,
			EnvApiKey:     "SOMEKEY",
			EnvScope:      "scope234",
			EnvCallback:   "http://mycallback.myshopify.com",
		},
		"error_response": {
			Request: events.APIGatewayProxyRequest{
				PathParameters: map[string]string{"x": "x"},
			},
			Response: golambda_helper.Response{
				StatusCode: 400,
				Body:       `{"message":"An error occurred processing the request."}`,
			},
			ResponseError: nil,
			EnvApiKey:     "SOMEKEY",
			EnvScope:      "scope234",
			EnvCallback:   "http://mycallback.myshopify.com",
		},
	}

	for name, test := range tests {
		t.Logf("Running test case: %s", name)

		envInterface := &mocks.EnvInterface{}
		envInterface.
			On("Getenv", "API_KEY").
			Return(test.EnvApiKey)
		envInterface.
			On("Getenv", "OAUTH_CALLBACK_URI").
			Return(test.EnvCallback)
		envInterface.
			On("Getenv", "SCOPES").
			Return(test.EnvScope)

		h := &main.EnvHandler{
			Env: envInterface,
		}

		response, err := h.Handler(test.Request)
		assert.Equal(t, test.ResponseError, err)
		assert.Equal(t, test.Response.StatusCode, response.StatusCode)
		// exclude checking the state since it changes with every incovation
		if response.StatusCode == 302 {
			assert.Equal(t, test.Response.Header.Location, response.Header.Location[:136])
		}
	}
}
