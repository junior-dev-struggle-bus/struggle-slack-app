package slackauth

import (
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

const (
	testMsg          = "Test Message"
	testReqTimeStamp = "1234"
	testSSK          = "myVerySecretSigningKey"
)

func newTestReq() *events.APIGatewayProxyRequest {
	return &events.APIGatewayProxyRequest{
		Headers: map[string]string{
			reqSigHeader:       calculateSig(testSSK, msgVersionNum, testReqTimeStamp, testMsg, msgDelim),
			reqTimeStampHeader: testReqTimeStamp,
		},
		Body: testMsg,
	}
}

func calculateSig(ssk string, versionNum string, reqTimeStamp string, msg string, delim string) string {
	msgToSign := strings.Join([]string{versionNum, reqTimeStamp, msg}, delim)
	return string(hmac256(msgToSign, ssk))
}

func assertFailedAuthN(t *testing.T, req *events.APIGatewayProxyRequest, ssk string) {
	result, err := AuthenticateLambdaReqWithSSK(req, ssk)
	assert.False(t, result)
	assert.Nil(t, err)
}

func assertAuthNInternalFailure(t *testing.T, req *events.APIGatewayProxyRequest, ssk string, sigOk bool, sigEmpty bool, reqTimeStampOk bool, reqTimeStampEmpty bool) {
	result, err := AuthenticateLambdaReqWithSSK(req, ssk)
	assert.False(t, result)
	assert.NotNil(t, err)
	assert.EqualError(t, err, fmt.Sprintf(authInternalErrFmt, sigOk, sigEmpty, reqTimeStampOk, reqTimeStampEmpty))
}

func TestAuthenticateLambdaRequest(t *testing.T) {
	res, err := AuthenticateLambdaReqWithSSK(newTestReq(), testSSK)
	assert.True(t, res)
	assert.Nil(t, err)
}

func TestAuthNLambdaRequestWithDiffSSK(t *testing.T) {
	testDiffSSK := "myMoreSecretSigningKeyThanTheVerySecretSigningKey"
	testReq := newTestReq()
	testReq.Headers[reqSigHeader] = calculateSig(testDiffSSK, msgVersionNum, testReqTimeStamp, testMsg, msgDelim)
	assertFailedAuthN(t, testReq, testSSK)
}

func TestAuthNLambdaReqWithDiffReqTimeStamp(t *testing.T) {
	testDiffReqTimeStamp := "5678"
	testReq := newTestReq()
	testReq.Headers[reqSigHeader] = calculateSig(testSSK, msgVersionNum, testDiffReqTimeStamp, testMsg, msgDelim)
	assertFailedAuthN(t, testReq, testSSK)
}

func TestAuthNLambdaReqWithDiffMsg(t *testing.T) {
	testDiffMsg := "Testing Testing Testing"
	testReq := newTestReq()
	testReq.Headers[reqSigHeader] = calculateSig(testSSK, msgVersionNum, testReqTimeStamp, testDiffMsg, msgDelim)
	assertFailedAuthN(t, testReq, testSSK)
}

func TestAuthNLambdaReqWithDiffVersionNum(t *testing.T) {
	testDiffVersionNum := "v123456"
	testReq := newTestReq()
	testReq.Headers[reqSigHeader] = calculateSig(testSSK, testDiffVersionNum, testReqTimeStamp, testMsg, msgDelim)
	assertFailedAuthN(t, testReq, testSSK)
}

func TestAuthNLambdaReqEmptySignatureHeader(t *testing.T) {
	testReq := newTestReq()
	testReq.Headers[reqSigHeader] = ""
	assertAuthNInternalFailure(t, testReq, testSSK, true, true, true, false)
}

func TestAuthNLambdaReqNoSignatureHeader(t *testing.T) {
	testReq := newTestReq()
	delete(testReq.Headers, reqSigHeader)
	assertAuthNInternalFailure(t, testReq, testSSK, false, true, true, false)
}

func TestAuthNLambdaReqEmptyTimeStampHeader(t *testing.T) {
	testReq := newTestReq()
	testReq.Headers[reqTimeStampHeader] = ""
	assertAuthNInternalFailure(t, testReq, testSSK, true, false, true, true)
}

func TestAuthNLambdaReqNoTimeStampHeader(t *testing.T) {
	testReq := newTestReq()
	delete(testReq.Headers, reqTimeStampHeader)
	assertAuthNInternalFailure(t, testReq, testSSK, true, false, false, true)
}

func TestAuthNLambdaReqSigWithVersionPrefix(t *testing.T) {
	testReq := newTestReq()
	testReq.Headers[reqSigHeader] = reqSigPrefix + testReq.Headers[reqSigHeader]

	testResult, testErr := AuthenticateLambdaReqWithSSK(testReq, testSSK)
	assert.True(t, testResult)
	assert.Nil(t, testErr)
}
