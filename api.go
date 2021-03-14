package main

import (
	"context"
	"encoding/json"
	"net/url"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
)

type errorResponse struct {
	Status string `json:"status"`
	Error  string `json:"error"`
}

type loginResponse struct {
	Status string `json:"status"`

	Title       string `json:"title"`
	Description string `json:"description"`

	IpPermissions []types.IpPermission `json:"ips"`
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

func routeLogin(context context.Context, request events.APIGatewayV2HTTPRequest, data url.Values, svc *ec2.Client) (events.APIGatewayV2HTTPResponse, error) {
	if data.Get("password") != currentConfig.PasswordHash {
		return jsonResponse(errorResponse{
			Status: "error",
			Error:  "Incorrect password.",
		})
	}

	result, err := svc.DescribeSecurityGroups(context, &ec2.DescribeSecurityGroupsInput{
		GroupIds: []string{
			currentConfig.SecurityGroupID,
		},
	})
	if err != nil {
		panic(err)
	}

	securityGroup := result.SecurityGroups[0]

	return jsonResponse(loginResponse{
		Status: "ok",

		Title:       currentConfig.Title,
		Description: currentConfig.Description,

		IpPermissions: securityGroup.IpPermissions,
	})
}
