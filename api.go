package main

import (
	"context"
	"encoding/json"
	"net/url"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
)

type errorResponse struct {
	Status string `json:"status"`
	Error  string `json:"error"`
}

type ipEntry struct {
	IP          string `json:"ip"`
	Description string `json:"description"`
}

type loginResponse struct {
	Status string `json:"status"`

	Title       string `json:"title"`
	Description string `json:"description"`

	ClientIP string `json:"clientIP"`

	IPs []ipEntry `json:"ips"`
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

	ipList := []ipEntry{}

	for _, permission := range securityGroup.IpPermissions {
		if permission.FromPort != int32(currentConfig.Port) || permission.ToPort != int32(currentConfig.Port) {
			continue
		}

		for _, ipRange := range permission.IpRanges {
			cidrIP := *ipRange.CidrIp
			if !strings.HasSuffix(cidrIP, "/32") {
				continue
			}

			description := ""
			if ipRange.Description != nil {
				description = *ipRange.Description
			}

			ipList = append(ipList, ipEntry{
				IP:          strings.Split(cidrIP, "/")[0],
				Description: description,
			})
		}
	}

	return jsonResponse(loginResponse{
		Status: "ok",

		Title:       currentConfig.Title,
		Description: currentConfig.Description,

		ClientIP: request.RequestContext.HTTP.SourceIP,

		IPs: ipList,
	})
}
