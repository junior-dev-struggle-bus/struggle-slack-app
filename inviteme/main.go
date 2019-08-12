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

// Bad request is a 400 error, which blames the user for doing something wrong
// Internal error is a 500 error, the service screwed up
// ie: If no request.Body comes in, that meant Slack didn't send the body over (service error)
func handler(request events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	authOK, err := slackauth.AuthenticateLambdaReq(&request)
	if err != nil {
		return &events.APIGatewayProxyResponse{
			StatusCode: http.StatusOK,
			Body:       ErrAuthLambdaReq.Error(),
		}, fmt.Errorf("%v: %v", ErrAuthLambdaReq, err)
	}

	if !authOK {
		return &events.APIGatewayProxyResponse{
			StatusCode: http.StatusOK,
			Body:       ErrAuthLambdaReq.Error(),
		}, fmt.Errorf("%v", ErrAuthLambdaReq)
	}

	parsedBody, err := url.ParseQuery(request.Body)
	if err != nil {
		return &events.APIGatewayProxyResponse{
			StatusCode: http.StatusOK,
			Body:       ErrParseRequestBody.Error(),
		}, fmt.Errorf("%v: %v", ErrParseRequestBody, err)
	}

	var receivedData SlackRequestBody
	err = decoder.Decode(&receivedData, parsedBody)
	if err != nil {
		return &events.APIGatewayProxyResponse{
			StatusCode: http.StatusOK,
			Body:       ErrDecodingParsedQuery.Error(),
		}, fmt.Errorf("%v: %v", ErrDecodingParsedQuery, err)
	}

	blockSection := SlackBlockSection{
		Type: SBlockSection,
		Text: SlackTypeText{
			Type: STypeMarkDown,
			Text: fmt.Sprintf("Hello <@%s>, fellow 🤓! \n%s\n\n🔗 %s 🔗", receivedData.UserID, linkMsg, inviteLink),
		},
	}

	blockDivider := SlackBlockDivider{Type: SBlockDivider}

	blockImage := SlackBlockImage{
		Type: SBlockImage,
		Title: SlackTypeText{
			Type: STypePlain,
			Text: imageText,
		},
		ImageURL: imageLink,
		AltText:  imageAltText,
	}

	blockContext := SlackBlockContext{
		Type: SBlockContext,
		Elements: []SlackTypeText{
			{
				Type: STypeMarkDown,
				Text: fmt.Sprintf("*Function by:* <@%s>", myJdsbSlackID),
			},
		},
	}

	slackPayload := SlackMessagePayload{
		ResponseType: STypeResponse,
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
			Body:       ErrJSONMarshal.Error(),
		}, fmt.Errorf("%v: %v", ErrJSONMarshal, err)
	}

	return &events.APIGatewayProxyResponse{
		Headers:    JSONHeader,
		StatusCode: http.StatusOK,
		Body:       string(slackPayloadByte),
	}, nil
}

func main() {
	// Make the handler available for Remote Procedure Call by AWS Lambda
	lambda.Start(handler)
}
