package main

import (
	"context"
	"embed"
	"errors"
	"io/fs"
	"io/ioutil"
	"net/http"
	"path"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

//go:embed frontend/*
var frontend embed.FS

var mimeTypes = map[string]string{
	".html": "text/html; charset=utf-8",
	".js":   "application/javascript",
	".css":  "text/css",
}

func HandleRequest(context context.Context, request events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	// check if this looks like an API route
	if request.RequestContext.HTTP.Method == http.MethodPost {
		if request.RawPath == "/login" {
			return routeLogin(context, request)
		}
	}

	// it does not, so it's probably a file then

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

	extension := path.Ext(filePath)
	mimeType, foundMimeType := mimeTypes[extension]
	if !foundMimeType {
		mimeType = "application/octet-stream"
	}

	return events.APIGatewayV2HTTPResponse{
		StatusCode: 200,
		Headers: map[string]string{
			"Content-Type": mimeType,
		},
		Body: string(data),
	}, nil
}

func main() {
	loadConfig()
	lambda.Start(HandleRequest)
}
