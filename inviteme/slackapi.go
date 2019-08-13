package main

import "errors"

// SlackRequestBody contains events.APIGatewayProxyRequest.Body received from slack
type slackRequestBody struct {
	Token        string   `schema:"token"`
	TeamID       string   `schema:"team_id"`
	TeamDomain   string   `schema:"team_domain"`
	ChannelID    string   `schema:"channel_id"`
	ChannelName  string   `schema:"channel_name"`
	UserID       string   `schema:"user_id"`
	UserName     string   `schema:"user_name"`
	Command      string   `schema:"command"`
	CommandTexts []string `schema:"text"`
	ResponseURL  string   `schema:"response_url"`
	TriggerID    string   `schema:"trigger_id"`
}

// SlackMessagePayload contains metadata to be JSONed for slack API to publish message
// Message Payload: https://api.slack.com/reference/messaging/payload
// Block Types: https://api.slack.com/reference/messaging/blocks
type slackMessagePayload struct {
	ResponseType string        `json:"response_type"`
	ChannelName  string        `json:"channel"`
	Token        string        `json:"token"`
	Blocks       []interface{} `json:"blocks"`
}

// SlackTypeText used for nesting objects inside a block
type slackTypeText struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

// SlackBlockSection is a block that will display simplest of texts
type slackBlockSection struct {
	Type string        `json:"type"`
	Text slackTypeText `json:"text"`
}

// SlackBlockDivider is a content divider to split up different blocks inside a message
type slackBlockDivider struct {
	Type string `json:"type"`
}

// SlackBlockImage is a block that displays images and its associated texts
type slackBlockImage struct {
	Type     string        `json:"type"`
	Title    slackTypeText `json:"title"`
	ImageURL string        `json:"image_url"`
	AltText  string        `json:"alt_text"`
}

// SlackBlockContext is a block that displays message context
type slackBlockContext struct {
	Type     string          `json:"type"`
	Elements []slackTypeText `json:"elements"`
}

// consts related to slack: header types and creating message layouts
const (
	headerContentType     = "Content-type"
	headerContentTypeJSON = "application/json"
	sBlockSection         = "section"
	sBlockDivider         = "divider"
	sBlockImage           = "image"
	sBlockContext         = "context"
	sTypeMarkDown         = "mrkdwn"
	sTypePlain            = "plain_text"
	sTypeResponse         = "in_channel"
)

// vars related to slack: headers and error messages
var (
	jSONHeader             = map[string]string{headerContentType: headerContentTypeJSON}
	errParseRequestBody    = errors.New("unable to parse APIGatewayProxyRequest.Body")
	errDecodingParsedQuery = errors.New("unable to decode parsed query into a struct")
	errJSONMarshal         = errors.New("unable to JSON encode for response")
	errAuthLambdaReq       = errors.New("unable to authenticate lambda request")
	userFriendlyErrorMsg   = "Uh oh, the service screwed up. Sorry about that! Please try again later."
)
