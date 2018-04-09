FROM golang:1.10.1 as golang_build

RUN go get github.com/aws/aws-lambda-go/lambda \
    && go get github.com/aws/aws-lambda-go/events \
    && go get github.com/stretchr/testify/assert \
    && go get github.com/aws/aws-sdk-go \
    && go get github.com/satori/go.uuid \
    && go get github.com/aws/aws-sdk-go/aws \
    && go get github.com/aws/aws-sdk-go/aws/session \
    && go get github.com/aws/aws-sdk-go/service/dynamodb \
    && go get github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute \
    && go get github.com/tkeech1/golambda_helper \
    && go get github.com/tkeech1/goshopify

COPY oauth_install.go handler_types.go ./
RUN env GOOS=linux go build -ldflags="-s -w" -o bin/oauth_install oauth_install.go handler_types.go

#--------------------------------

FROM node:9.11.1-alpine as serverless

RUN npm install -g serverless --unsafe-perm=true
RUN npm install --save serverless-python-requirements
RUN mkdir /app
WORKDIR /app
COPY --from=golang_build /go/bin/oauth_install .
COPY serverless.yml .