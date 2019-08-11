package main

import (
	"encoding/json"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/gorilla/schema"
	"net/http"
	"net/url"
)

const (
	inviteLink    string = "https://jdsb-test-ground-slack-invite.herokuapp.com/"
	imageLink     string = "http://giphygifs.s3.amazonaws.com/media/5q7XU3MurT0pW/giphy.gif"
	imageText     string = "I came in like a wrecking ball"
	imageAltText  string = "Picture of a cat wrecking neatly stacked books"
	linkMsg       string = "Click on the link below to get an invite to *jdsb-wrecking-ball* workspace: "
	myJdsbSlackID        = "UGAA0BJ4R"
)

var (
	decoder = schema.NewDecoder()
)

func handler(request events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	parsedBody, err := url.ParseQuery(request.Body)
	if err != nil {
		return &events.APIGatewayProxyResponse{
			StatusCode: http.StatusBadRequest,
			Body:       fmt.Sprintf("%d: %s", http.StatusBadRequest, ErrParseRequestBody.Error()),
		}, ErrParseRequestBody
	}

	var receivedData SlackRequestBody
	err = decoder.Decode(&receivedData, parsedBody)
	if err != nil {
		return &events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       fmt.Sprintf("%d: %s", http.StatusInternalServerError, ErrDecodingParsedQuery.Error()),
		}, ErrDecodingParsedQuery
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
		Type: STypeImage,
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
			StatusCode: http.StatusInternalServerError,
			Body:       fmt.Sprintf("%d: %s", http.StatusInternalServerError, ErrJSONMarshal.Error()),
		}, ErrJSONMarshal
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
