package main

import (
	"context"
	"embed"
	"encoding/base64"
	"errors"
	"io/fs"
	"io/ioutil"
	"net/http"
	"net/url"
	"path"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"

	awsConfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
)

//go:embed frontend/*
var frontend embed.FS

var mimeTypes = map[string]string{
	".html": "text/html; charset=utf-8",
	".js":   "application/javascript",
	".css":  "text/css",
}

func messageResponse(status int, message string) (events.APIGatewayV2HTTPResponse, error) {
	return events.APIGatewayV2HTTPResponse{
		StatusCode: status,
		Headers: map[string]string{
			"Content-Type": "text/html",
		},
		Body: message,
	}, nil
}

func HandleRequest(context context.Context, request events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	// check if this looks like an API route
	if request.RequestContext.HTTP.Method == http.MethodPost {
		requestType, found := request.Headers["content-type"]
		if !found || requestType != "application/x-www-form-urlencoded" {
			return messageResponse(http.StatusBadRequest, "Bad request Content-Type.")
		}

		body := request.Body
		if request.IsBase64Encoded {
			bodyBytes, err := base64.StdEncoding.DecodeString(request.Body)
			if err != nil {
				panic(err)
			}
			body = string(bodyBytes)
		}

		data, err := url.ParseQuery(body)
		if err != nil {
			panic(err)
		}

		// set up the aws sdk
		cfg, err := awsConfig.LoadDefaultConfig(context, awsConfig.WithRegion(currentConfig.Region))
		if err != nil {
			panic(err)
		}
		svc := ec2.NewFromConfig(cfg)

		if request.RawPath == "/login" {
			return routeLogin(context, request, data, svc)
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
