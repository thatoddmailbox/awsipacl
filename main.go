package main

import (
	"context"
	"embed"
	"errors"
	"io/fs"
	"io/ioutil"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

//go:embed frontend/*
var frontend embed.FS

func HandleRequest(context context.Context, request events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	filePath := request.RawPath
	if request.RawPath == "/" {
		filePath = "/index.html"
	}
	filePath = "frontend" + filePath

	file, err := frontend.Open(filePath)
	if errors.Is(err, fs.ErrNotExist) {
		// file not found
		return events.APIGatewayV2HTTPResponse{
			StatusCode: 404,
			Headers: map[string]string{
				"Content-Type": "text/html",
			},
			Body: "File not found.",
		}, nil
	} else if err != nil {
		// internal server error
		// TODO: handle better
		panic(err)
	}
	defer file.Close()

	data, err := ioutil.ReadAll(file)
	if err != nil {
		panic(err)
	}

	return events.APIGatewayV2HTTPResponse{
		StatusCode: 200,
		Headers: map[string]string{
			"Content-Type": "text/html",
		},
		Body: string(data),
	}, nil
}

func main() {
	lambda.Start(HandleRequest)
}
