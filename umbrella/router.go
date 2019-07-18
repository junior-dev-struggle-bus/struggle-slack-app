package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"log"
	"net/http"
	"os"
	"strings"
)

const (
	slackAuthNTokenHeader        = "x-slack-signature"
	slackRequestTimestampHeader  = "x-slack-request-timestamp"
	slackVersionNumber           = "v0"
	slackMessageDelimiter        = ":"
	statusForbiddenMsgFmt        = "uh uh uh...you didn't say the magic word. (%s)"
	statusInternalServerErrorFmt = "Sorry...we uh...messed up. (%s)"
	logResponseStatement         = "%s:\n%+v\n"
	logAuthInfoStatement         = "\nauthNToken: %s\nauthNOk: %t\nreqTimeStamp: %s\nreqTimeStampOk: %t\n"
	logHmacInfoStatement         = "\nauthNToken: %s\nexpectedAuthNToken: %s\n"
)

var (
	// Configured in Netlify's environment variables.
	slackAppSigningSecretKey = os.Getenv("SLACK_SIGNING_SECRET_KEY")
)

// Handles things...duh
// 1. Authenticate the request.
// 2. Route the request to the appropriate function.
//   * Extract the function name from the request.
//   * Check whether the function name exists.
//   * Retrieves a function URL endpoint to send a request to.
//   * Create/send request to endpoint with arguments from this request.
// 3. Send immediate status OK reponse to caller unless authN failed.
// 4. Create/send a request to response_url once response returns from endpoint.
func handler(request events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	authNToken, authNOk := request.Headers[slackAuthNTokenHeader]
	reqTimeStamp, reqTimeStampOk := request.Headers[slackRequestTimestampHeader]

	log.Printf(logAuthInfoStatement, authNToken, authNOk, string(reqTimeStamp), reqTimeStampOk)

	body := request.Body
	resp := &events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Body:       http.StatusText(http.StatusOK),
	}

	if !authNOk || !reqTimeStampOk {
		resp.StatusCode = http.StatusInternalServerError
		resp.Body = fmt.Sprintf(statusInternalServerErrorFmt, http.StatusText(http.StatusInternalServerError))
		log.Printf(logResponseStatement, http.StatusText(http.StatusInternalServerError), resp)
		return resp, nil
	}

	message := strings.Join([]string{slackVersionNumber, reqTimeStamp, body}, slackMessageDelimiter)
	// The Slack AuthN token comes prefixed with 'v0=', which is not used in the HMAC equality procedure.
	authNTokenWoutVersion := strings.Split(authNToken, "=")[1]
	if ok := authenticate(message, authNTokenWoutVersion, slackAppSigningSecretKey); !ok {
		resp.StatusCode = http.StatusForbidden
		resp.Body = fmt.Sprintf(statusForbiddenMsgFmt, http.StatusText(http.StatusForbidden))
		log.Printf(logResponseStatement, http.StatusText(http.StatusForbidden), resp)
		return resp, nil
	}
	//TODO The following is where the request would be routed to the proper function
	//     on a separate thread.
	//go routeRequest(request)

	log.Printf(logResponseStatement, http.StatusText(http.StatusOK), resp)
	return resp, nil
}

// Authenticate things...duh
// Uses a secure hashing function to combine the message and the key to create
// a secured hash unique to the message. This is dependent on the timestamp so
// we use the built-in HMAC equal function to verify equality rather than regular
// equals. This is to avoid timesliding issues.
func authenticate(message string, authNToken string, key string) bool {
	mac := hmac.New(sha256.New, []byte(key))
	mac.Write([]byte(message))
	expectedAuthNToken := hex.EncodeToString(mac.Sum(nil))
	log.Printf(logHmacInfoStatement, authNToken, expectedAuthNToken)
	return hmac.Equal([]byte(authNToken), []byte(expectedAuthNToken))
}

func main() {
	lambda.Start(handler)
}
