package main

import (
	"encoding/json"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"time"
)

const (
	currComicInfoUrl    = "https://xkcd.com/info.0.json"
	randComicInfoUrlFmt = "https://xkcd.com/%d/info.0.json"
	timeOut             = 2
	contentTypeHeader   = "Content-type"
	contentTypeJson     = "application/json"
	slackInChannel      = "in_channel"
	slackEphemeral      = "ephemeral"
)

type xkcdComicInfo struct {
	Num   int
	Title string
	Alt   string
	Img   string
}

type slackResponse struct {
	Text          string            `json:"text"`
	Response_type string            `json:"response_type"`
	Attachments   []slackAttachment `json:"attachments"`
}

type slackAttachment struct {
	Text string `json:"text"`
}

func handler(request events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	resp := &events.APIGatewayProxyResponse{
		Headers:    make(map[string]string),
		StatusCode: http.StatusOK,
		Body:       http.StatusText(http.StatusOK),
	}

	xkcdClient := http.Client{
		Timeout: time.Second * timeOut,
	}

	req, err := http.NewRequest(http.MethodGet, currComicInfoUrl, nil)
	res, err := xkcdClient.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	currXkcdComicInfo := &xkcdComicInfo{}
	err = json.Unmarshal(body, &currXkcdComicInfo)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Curr Comic: %+v", *currXkcdComicInfo)

	rand.Seed(time.Now().UnixNano())
	randComicInfoUrl := fmt.Sprintf(randComicInfoUrlFmt, rand.Intn(currXkcdComicInfo.Num)+1)
	req, err = http.NewRequest(http.MethodGet, randComicInfoUrl, nil)
	res, err = xkcdClient.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	body, err = ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	randXkcdComicInfo := &xkcdComicInfo{}
	err = json.Unmarshal(body, &randXkcdComicInfo)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Rand Comic: %+v", *randXkcdComicInfo)

	slackResp := &slackResponse{
		Text:          randXkcdComicInfo.Img + "\n" + randXkcdComicInfo.Title,
		Response_type: slackInChannel,
		Attachments: []slackAttachment{
			slackAttachment{
				Text: randXkcdComicInfo.Alt,
			},
		},
	}

	slackRespByte, err := json.Marshal(slackResp)
	if err != nil {
		log.Fatal(err)
	}

	resp.Headers[contentTypeHeader] = contentTypeJson
	resp.Body = string(slackRespByte)
	log.Printf("Response Body:'%s'\n", resp.Body)

	return resp, nil
}

func main() {
	lambda.Start(handler)
}
