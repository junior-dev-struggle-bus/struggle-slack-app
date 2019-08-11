package main

import (
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func handler(request events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	fmt.Println("lisa was here")
	fmt.Println("Body--", request.Body)
	fmt.Println("Headers--", request.Headers)
	fmt.Println("HTTPMethod--", request.HTTPMethod)
	fmt.Println("IsBase64Encoded--", request.IsBase64Encoded)
	fmt.Println("MultiValueHeaders--", request.MultiValueHeaders)
	fmt.Println("MultiValueQueryStringParameters--", request.MultiValueQueryStringParameters)
	fmt.Println("Path--", request.Path)
	fmt.Println("PathParameters--", request.PathParameters)
	fmt.Println("QueryStringParameters--", request.QueryStringParameters)
	fmt.Println("RequestContext--", request.RequestContext)
	fmt.Println("Resource--   ", request.Resource)
	fmt.Println("StageVariables--", request.StageVariables)

	return &events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       "Hello, World",
	}, nil
}

func main() {
	// Make the handler available for Remote Procedure Call by AWS Lambda
	lambda.Start(handler)
}