package main

import "errors"

// SlackRequestBody contains events.APIGatewayProxyRequest.Body received from slack
type SlackRequestBody struct {
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
type SlackMessagePayload struct {
	ResponseType string        `json:"response_type"`
	ChannelName  string        `json:"channel"`
	Token        string        `json:"token"`
	Blocks       []interface{} `json:"blocks"`
}

// SlackTypeText used for nesting objects inside a block
type SlackTypeText struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

// SlackBlockSection is a block that will display simplest of texts
type SlackBlockSection struct {
	Type string        `json:"type"`
	Text SlackTypeText `json:"text"`
}

// SlackBlockDivider is a content divider to split up different blocks inside a message
type SlackBlockDivider struct {
	Type string `json:"type"`
}

// SlackBlockImage is a block that displays images and its associated texts
type SlackBlockImage struct {
	Type     string        `json:"type"`
	Title    SlackTypeText `json:"title"`
	ImageURL string        `json:"image_url"`
	AltText  string        `json:"alt_text"`
}

// SlackBlockContext is a block that displays message context
type SlackBlockContext struct {
	Type     string          `json:"type"`
	Elements []SlackTypeText `json:"elements"`
}

// consts related to slack: header types and creating message layouts
const (
	HeaderContentType     = "Content-type"
	HeaderContentTypeJSON = "application/json"
	SBlockSection         = "section"
	SBlockDivider         = "divider"
	SBlockImage           = "image"
	SBlockContext         = "context"
	STypeMarkDown         = "mrkdwn"
	STypePlain            = "plain_text"
	STypeResponse         = "in_channel"
)

// vars related to slack: headers and error messages
var (
	JSONHeader             = map[string]string{HeaderContentType: HeaderContentTypeJSON}
	ErrParseRequestBody    = errors.New("unable to parse APIGatewayProxyRequest.Body")
	ErrDecodingParsedQuery = errors.New("unable to decode parsed query into a struct")
	ErrJSONMarshal         = errors.New("unable to JSON encode for response")
	ErrAuthLambdaReq       = errors.New("unable to authenticate lambda request")
)
