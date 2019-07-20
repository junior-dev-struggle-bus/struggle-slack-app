package slackauth

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"os"
	"strings"
)

const (
	reqSigHeader       = "x-slack-signature"
	reqSigPrefix       = "v0="
	reqTimeStampHeader = "x-slack-request-timestamp"
	msgVersionNum      = "v0"
	msgDelim           = ":"
	authInternalErrFmt = "\nInternal Error, sigExists: %t, sigEmpty: %t, timeStampExists: %t, timeStampEmpty: %t\n"
	slackAppSSKEnvVar  = "SLACK_APP_SIGNING_SECRET_KEY"
)

var (
	appSigningSecretKey = os.Getenv(slackAppSSKEnvVar)
)

// AuthenticateLambdaReq authenticates the Lambda request using the signing secret
// key from the environment variable by default. It is a wrapper around the method
// AuthenticateLambdaReqWithSSK. Its returns patterns are the same as that method.
func AuthenticateLambdaReq(request *events.APIGatewayProxyRequest) (bool, error) {
	return AuthenticateLambdaReqWithSSK(request, appSigningSecretKey)
}

// AuthenticateLambdaReqWithSSK authenticates the Lambda request using the given signing secret key.
// It uses HMAC SHA256 to sign the message with the key. If the authentication succeeds, the function
// returns true and a nil error. If there's a missing or empty signature, or a missing or empty request
// timestamp, it returns false and an error indicating an internal error. Otherwise, it returns false
// and nil error to indicate the request is unauthentic.
func AuthenticateLambdaReqWithSSK(request *events.APIGatewayProxyRequest, ssk string) (bool, error) {
	reqSig, reqSigOk := request.Headers[reqSigHeader]
	reqTimeStamp, reqTimeStampOk := request.Headers[reqTimeStampHeader]

	if !reqSigOk || !reqTimeStampOk || reqSig == "" || reqTimeStamp == "" {
		return false, errors.New(fmt.Sprintf(authInternalErrFmt, reqSigOk, reqSig == "", reqTimeStampOk, reqTimeStamp == ""))
	}

	reqBody := request.Body
	msg := strings.Join([]string{msgVersionNum, reqTimeStamp, reqBody}, msgDelim)
	// The Slack AuthN token comes prefixed with 'v0=', which is not used in the HMAC equality procedure.
	reqSigWOutVersion := ""
	if strings.HasPrefix(reqSig, reqSigPrefix) {
		reqSigWOutVersion = reqSig[3:len(reqSig)]
	} else {
		reqSigWOutVersion = reqSig
	}

	if ok := authenticate(msg, reqSigWOutVersion, ssk); !ok {
		return false, nil
	}

	return true, nil
}

// Authenticate uses HMAC SHA256 to hash the message with the key. It
// then compares the resulting hash with the target message authentication
// code. Note, this does not use standard equality for the comparison.
// It instead uses an HMAC-safe equality check to avoid timesliding issues
// that could affect the resulting hash.
func authenticate(message string, targetMac string, key string) bool {
	expectedMac := hmac256(message, key)
	return hmac.Equal([]byte(targetMac), expectedMac)
}

func hmac256(message string, key string) []byte {
	mac := hmac.New(sha256.New, []byte(key))
	mac.Write([]byte(message))
	return []byte(hex.EncodeToString(mac.Sum(nil)))
}
