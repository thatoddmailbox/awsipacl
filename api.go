package main

import (
	"context"
	"encoding/json"
	"net/url"

	"github.com/aws/aws-lambda-go/events"
)

type errorResponse struct {
	Status string `json:"status"`
	Error  string `json:"error"`
}

type loginResponse struct {
	Status string `json:"status"`

	Title       string `json:"title"`
	Description string `json:"description"`
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

func routeLogin(context context.Context, request events.APIGatewayV2HTTPRequest, data url.Values) (events.APIGatewayV2HTTPResponse, error) {
	if data.Get("password") != currentConfig.PasswordHash {
		return jsonResponse(errorResponse{
			Status: "error",
			Error:  "Incorrect password.",
		})
	}

	return jsonResponse(loginResponse{
		Status:      "ok",
		Title:       currentConfig.Title,
		Description: currentConfig.Description,
	})
}
