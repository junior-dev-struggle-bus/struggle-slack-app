package main

import (
	"encoding/json"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/gorilla/schema"
	slackauth "github.com/phoenixcoder/slack-golang-sdk/auth"
	"net/http"
	"net/url"
)

const (
	inviteLink    = "https://jdsb-test-ground-slack-invite.herokuapp.com/"
	imageLink     = "http://giphygifs.s3.amazonaws.com/media/5q7XU3MurT0pW/giphy.gif"
	imageText     = "I came in like a wrecking ball"
	imageAltText  = "Picture of a cat wrecking neatly stacked books"
	linkMsg       = "Click on the link below to get an invite to *jdsb-wrecking-ball* workspace: "
	myJdsbSlackID = "UGAA0BJ4R"
)

var (
	decoder = schema.NewDecoder()
)

func handler(request events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	authOK, err := slackauth.AuthenticateLambdaReq(&request)
	if err != nil {
		return &events.APIGatewayProxyResponse{
			StatusCode: http.StatusOK,
			Body:       userFriendlyErrorMsg,
		}, fmt.Errorf("%v: %v", errAuthLambdaReq, err)
	}

	if !authOK {
		return &events.APIGatewayProxyResponse{
			StatusCode: http.StatusOK,
			Body:       userFriendlyErrorMsg,
		}, fmt.Errorf("%v", errAuthLambdaReq)
	}

	parsedBody, err := url.ParseQuery(request.Body)
	if err != nil {
		return &events.APIGatewayProxyResponse{
			StatusCode: http.StatusOK,
			Body:       userFriendlyErrorMsg,
		}, fmt.Errorf("%v: %v", errParseRequestBody, err)
	}

	var receivedData slackRequestBody
	err = decoder.Decode(&receivedData, parsedBody)
	if err != nil {
		return &events.APIGatewayProxyResponse{
			StatusCode: http.StatusOK,
			Body:       userFriendlyErrorMsg,
		}, fmt.Errorf("%v: %v", errDecodingParsedQuery, err)
	}

	blockSection := slackBlockSection{
		Type: sBlockSection,
		Text: slackTypeText{
			Type: sTypeMarkDown,
			Text: fmt.Sprintf("Hello <@%s>, fellow 🤓! \n%s\n\n🔗 %s 🔗", receivedData.UserID, linkMsg, inviteLink),
		},
	}

	blockDivider := slackBlockDivider{Type: sBlockDivider}

	blockImage := slackBlockImage{
		Type: sBlockImage,
		Title: slackTypeText{
			Type: sTypePlain,
			Text: imageText,
		},
		ImageURL: imageLink,
		AltText:  imageAltText,
	}

	blockContext := slackBlockContext{
		Type: sBlockContext,
		Elements: []slackTypeText{
			{
				Type: sTypeMarkDown,
				Text: fmt.Sprintf("*Function by:* <@%s>", myJdsbSlackID),
			},
		},
	}

	slackPayload := slackMessagePayload{
		ResponseType: sTypeResponse,
		ChannelName:  receivedData.ChannelName,
		Token:        receivedData.Token,
		Blocks: []interface{}{
			blockSection,
			blockDivider,
			blockImage,
			blockContext,
		},
	}

	slackPayloadByte, err := json.Marshal(slackPayload)
	if err != nil {
		return &events.APIGatewayProxyResponse{
			StatusCode: http.StatusOK,
			Body:       userFriendlyErrorMsg,
		}, fmt.Errorf("%v: %v", errJSONMarshal, err)
	}

	return &events.APIGatewayProxyResponse{
		Headers:    jSONHeader,
		StatusCode: http.StatusOK,
		Body:       string(slackPayloadByte),
	}, nil
}

func main() {
	// Make the handler available for Remote Procedure Call by AWS Lambda
	lambda.Start(handler)
}
