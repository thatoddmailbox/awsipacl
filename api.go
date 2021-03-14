package main

import (
	"context"
	"encoding/json"

	"github.com/aws/aws-lambda-go/events"
)

type loginResponse struct {
	Status string `json:"status"`
}

func jsonResponse(data interface{}) (events.APIGatewayV2HTTPResponse, error) {
	dataJSON, err := json.Marshal(data)
	if err != nil {
		// TODO: handle better
		panic(err)
	}

	return events.APIGatewayV2HTTPResponse{
		StatusCode: 200,
		Headers: map[string]string{
			"Content-Type": "text/json",
		},
		Body: string(dataJSON),
	}, nil
}

func routeLogin(context context.Context, request events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	return jsonResponse(loginResponse{"ok"})
}
