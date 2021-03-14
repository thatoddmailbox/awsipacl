package main

import (
	"context"
	"encoding/json"
	"net/url"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
)

type statusResponse struct {
	Status string `json:"status"`
}

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

func routeAdd(context context.Context, request events.APIGatewayV2HTTPRequest, data url.Values, svc *ec2.Client) (events.APIGatewayV2HTTPResponse, error) {
	if data.Get("password") != currentConfig.PasswordHash {
		return jsonResponse(errorResponse{
			Status: "error",
			Error:  "Incorrect password.",
		})
	}

	ip := data.Get("ip") + "/32"
	description := data.Get("description")

	_, err := svc.AuthorizeSecurityGroupIngress(context, &ec2.AuthorizeSecurityGroupIngressInput{
		GroupId: &currentConfig.SecurityGroupID,

		IpPermissions: []types.IpPermission{
			{
				FromPort:   currentConfig.Port,
				ToPort:     currentConfig.Port,
				IpProtocol: &currentConfig.Protocol,
				IpRanges: []types.IpRange{
					{
						CidrIp:      &ip,
						Description: &description,
					},
				},
			},
		},
	})
	if err != nil {
		return jsonResponse(errorResponse{
			Status: "error",
			Error:  err.Error(),
		})
	}

	return jsonResponse(statusResponse{"ok"})
}

func routeDelete(context context.Context, request events.APIGatewayV2HTTPRequest, data url.Values, svc *ec2.Client) (events.APIGatewayV2HTTPResponse, error) {
	if data.Get("password") != currentConfig.PasswordHash {
		return jsonResponse(errorResponse{
			Status: "error",
			Error:  "Incorrect password.",
		})
	}

	ip := data.Get("ip") + "/32"

	_, err := svc.RevokeSecurityGroupIngress(context, &ec2.RevokeSecurityGroupIngressInput{
		GroupId: &currentConfig.SecurityGroupID,

		IpPermissions: []types.IpPermission{
			{
				FromPort:   currentConfig.Port,
				ToPort:     currentConfig.Port,
				IpProtocol: &currentConfig.Protocol,
				IpRanges: []types.IpRange{
					{
						CidrIp:      &ip,
						Description: nil,
					},
				},
			},
		},
	})
	if err != nil {
		return jsonResponse(errorResponse{
			Status: "error",
			Error:  err.Error(),
		})
	}

	return jsonResponse(statusResponse{"ok"})
}
