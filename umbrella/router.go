package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	slackauth "github.com/phoenixcoder/slack-golang-sdk/auth"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"path"
	"strings"
	"time"
)

const (
	timeout             = 2
	forbiddenErrRespMsg = "uh uh uh...you didn't say the magic word."
	internalErrRespMsg  = "Sorry...we uh...messed up."
	// TODO Make sure internal error messages are transferred across more friendly-like.
	funcNotFoundErrMsgFmt   = "We're embarrassed for you, but we don't know a '%s'. Try these instead:\n'%s'"
	logMsg                  = "%s, %s"
	routerRegistryUrlEnvVar = "ROUTER_REGISTRY_URL_ENV_VAR"
	commandRegistryLoc      = "./registry.json"
	cmdRequestUrlRoot       = "requestUrl"
	cmdFunctionsReg         = "functions"
	contentTypeHeader       = "content-type"
	slackCmdParam           = "command"
	slackArgsParam          = "text"
	slackResponseUrlParam   = "response_url"
	// TODO Download registry during worker setup.
	tempCmdReg = `{
    "struggle" : {
        "requestUrl" : "https://admiring-meninsky-dcfbdb.netlify.com/.netlify/functions/",
        "helpKeyword" : "help",
        "functions" : {
            "functionName" : {
                "usage" : "/struggle functionname arg1 arg2 arg3...",
                "description" : "This describes what your function does when they use /struggle help functionName. This should also describe how it uses the arguments.",
                "manual" : "Optional documentation website for your command."
            },
            "randxkcd" : {
                "usage" : "/struggle randxkcd",
                "description" : "Returns a random xkcd comic.",
                "manual" : "THERE IS NO DOCUMENTATION WEBSITE. FIGURE IT OUT YOURSELF! -YOUR GRUMPY PHOENIXCODER"                
            }
        }
    }
}`
)

var (
	// TODO Convert this registry to a generic Data Access Object (DAO)
	//      that can pull from any data source.
	commandRegistry map[string]cmdInfo
	registryUrl     = os.Getenv(routerRegistryUrlEnvVar)
)

type cmdInfo struct {
	RequestUrl  string
	HelpKeyword string
	Functions   map[string]funcRoutingInfo
}

type funcRoutingInfo struct {
	Name        string
	RequestUrl  string
	ResponseUrl string
	Usage       string
	Description string
	Manual      string
}

// TODO Generalize away from Slack patterns.
func getFuncRoutingInfo(values url.Values, cmdReg *map[string]cmdInfo) (*funcRoutingInfo, error) {
	// TODO Polish the error handling here for Slack.
	cmdList, cmdOk := values[slackCmdParam]
	if !cmdOk {
		return nil, errors.New("Command argument was not received.")
	}

	if len(cmdList) != 1 {
		return nil, fmt.Errorf("Must have sent exactly 1 command. Command List: '%v'", cmdList)
	}

	cmd := strings.TrimLeft(cmdList[0], "/")
	cmdInfo, cmdInfoExists := (*cmdReg)[cmd]
	if !cmdInfoExists {
		return nil, fmt.Errorf("Command is not registered. Command Name: '%s'", cmd)
	}

	funcReg := cmdInfo.Functions
	if len(funcReg) <= 0 {
		return nil, fmt.Errorf("Function registry for the command does not exist. Command Name: '%s'", cmd)
	}

	argsList, argsOk := values[slackArgsParam]
	if !argsOk {
		return nil, errors.New("Arguments list was not received.")
	}

	if len(argsList) <= 0 {
		return nil, nil
	}

	funcName := argsList[0]
	funcRoutingInfo, funcRoutingInfoOk := funcReg[funcName]
	if !funcRoutingInfoOk {
		return nil, fmt.Errorf("Function was not registered. Function Name: '%s'", funcName)
	}

	requestUrlRoot, err := url.Parse(cmdInfo.RequestUrl)
	if err != nil {
		return nil, err
	}

	requestUrl := requestUrlRoot
	requestUrl.Path = path.Join(requestUrlRoot.Path, funcName)
	funcRoutingInfo.Name = funcName
	funcRoutingInfo.RequestUrl = requestUrl.String()
	return &funcRoutingInfo, nil
}

func loadCommandRegistryFromFile(location string, registry *map[string]cmdInfo) {
	// TODO If the file ever gets large enough, opening and unmarshalling it will cause
	//      function overhead to increase, eventually to a point where invokes fail. If that
	//      happens, it's time to move to a DB backend to retrieve function relay info.
	// TODO Add metrics on time to load.
	contents, err := ioutil.ReadFile(location)
	if err != nil {
		log.Fatalf("Could not open and read function registry. Error: %v\n", err)
	}
	log.Printf("%s\n", contents)

	loadCommandRegistryFromContents(contents, registry)
}

func loadCommandRegistryFromContents(contents []byte, registry *map[string]cmdInfo) {
	if err := json.Unmarshal(contents, registry); err != nil {
		log.Fatalf("Could not unmarshal registry contents. Error: %v\n", err)
	}
}

func loadCommandRegistryFromUrl(url string, registry *map[string]cmdInfo) {
	log.Printf("Registry URL: %s\n", url)
	resp, err := http.Get(url)
	if err != nil {
		log.Fatalf("Could not download registry contents. Error: %v\n", err)
	}

	contents, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Could not read registry contents. Error: %v\n", err)
	}

	loadCommandRegistryFromContents(contents, registry)
}

func setInternalErrCode(resp *events.APIGatewayProxyResponse, reason string) {
	setErredStatusCode(resp, internalErrRespMsg, reason, http.StatusInternalServerError)
}

func setForbiddenErrCode(resp *events.APIGatewayProxyResponse) {
	setErredStatusCode(resp, forbiddenErrRespMsg, "You're just not allowed.", http.StatusForbidden)
}

func setErredStatusCode(resp *events.APIGatewayProxyResponse, msg string, reason string, statusCode int) {
	resp.StatusCode = statusCode
	resp.Body = msg + " (" + http.StatusText(resp.StatusCode) + ")"
	log.Printf(logMsg, resp.Body, reason)
}

// Handles things...duh
// 1. Authenticate the request.
// 2. Route the request to the appropriate function.
//    * Extract the function name from the request.
//    * Check whether the function name exists.
//    * Retrieves a function URL endpoint to send a request to.
//    * Create/send request to endpoint with arguments from this request.
// 3. Send immediate status OK reponse to caller unless authN failed.
// 4. Create/send a request to response_url once response returns from endpoint.
// TODO Currently, this routing mechanism is very specific to Slack. We may want
//      to turn this into a general routing mechanism since this pattern shows up
//      often.
func handler(request events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	resp := &events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Body:       http.StatusText(http.StatusOK),
	}

	authNOk, err := slackauth.AuthenticateLambdaReq(&request)
	if err != nil {
		setInternalErrCode(resp, err.Error())
		return resp, nil
	}

	if !authNOk {
		setForbiddenErrCode(resp)
		return resp, nil
	}

	values, err := url.ParseQuery(request.Body)
	if err != nil {
		setInternalErrCode(resp, err.Error())
		return resp, nil
	}

	funcRoutingInfo, err := getFuncRoutingInfo(values, &commandRegistry)
	if err != nil {
		setInternalErrCode(resp, err.Error())
		return resp, nil
	}

	// The responseUrl extraction is not a part of the function
	// routing info extraction since it is Slack-specific.
	// TODO Generalize an alternate response-path.
	responseUrlList, responseUrlListOk := values[slackResponseUrlParam]
	// If this is not here, you want to respond right away before
	// the Slack service closes their end of the connection. It will
	// be the only way to speak to the user. Otherwise, they'll see
	// Slack's response, which may be nothing at all.
	if !responseUrlListOk || len(responseUrlList) <= 0 {
		setInternalErrCode(resp, "Could not parse response url.")
		return resp, nil
	}
	funcRoutingInfo.ResponseUrl = responseUrlList[0]

	resp, err = routeRequestWaitForResp(funcRoutingInfo, &request)

	log.Printf(logMsg, http.StatusText(http.StatusOK), resp.Body)
	return resp, nil
}

func routeRequestWaitForResp(routeInfo *funcRoutingInfo, request *events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	log.Printf("Routing Info: '%+v\n'", *routeInfo)
	log.Printf("Request to Route: '%+v\n'", *request)
	resp := &events.APIGatewayProxyResponse{
		Headers:    make(map[string]string),
		StatusCode: http.StatusOK,
		Body:       http.StatusText(http.StatusOK),
	}

	routerClient := http.Client{
		Timeout: time.Second * timeout,
	}

	req, err := http.NewRequest(http.MethodPost, routeInfo.RequestUrl, strings.NewReader(request.Body))
	if err != nil {
		setInternalErrCode(resp, err.Error())
		return resp, err
	}

	routedResp, err := routerClient.Do(req)
	if err != nil {
		setInternalErrCode(resp, err.Error())
		return resp, err
	}
	log.Printf("Routed Response: %+v\n", routedResp)

	defer routedResp.Body.Close()
	resp.StatusCode = routedResp.StatusCode
	if err != nil {
		setInternalErrCode(resp, err.Error())
		return resp, err
	}

	body, err := ioutil.ReadAll(routedResp.Body)
	if err != nil {
		setInternalErrCode(resp, err.Error())
		return resp, err
	}
	resp.Body = string(body)

	resp.Headers[contentTypeHeader] = routedResp.Header.Get(contentTypeHeader)
	log.Printf("Delegate Function Response: '%+v'\n", resp)
	return resp, err
}

func main() {
	// TODO Perform smart loading of contents for local testing. Flags > Local Variables > Environment Variable > Configuration File search.
	loadCommandRegistryFromUrl(registryUrl, &commandRegistry)
	log.Printf("Loaded Command Registry: %+v\n", commandRegistry)
	lambda.Start(handler)
}
