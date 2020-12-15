package github

// TODO paginated responses

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

// Context of current Github Actions workflow call.
type Context struct {
	Token      string
	Repository string
	Ref        string
	Sha        string
}

var httpClient *http.Client = &http.Client{}

// NewDefaultContext ...
func NewDefaultContext() *Context {
	return NewContext(os.Getenv("GITHUB_CONTEXT"))
}

// NewContext ...
func NewContext(jsonContext string) *Context {
	context := &Context{}
	json.Unmarshal([]byte(jsonContext), context)
	return context
}

// APICall executes an Github API on the given context, using the provided
// endpoint and values.
func (context *Context) APICall(
	endpoint *Endpoint,
	requestBody []byte,
	values ...interface{}) (string, error) {

	url := fmt.Sprintf(
		"https://api.github.com/repos/%s/%s",
		context.Repository,
		fmt.Sprintf(endpoint.url, values...))
	fmt.Printf("Request: %s\n%s\n", url, string(requestBody))

	request, err := http.NewRequest(
		endpoint.method,
		url,
		bytes.NewBuffer(requestBody))
	if err != nil {
		return "", err
	}

	request.Header.Add("authorization", fmt.Sprintf("Bearer %s", context.Token))
	response, err := httpClient.Do(request)
	if err != nil {
		return "", err
	}

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return "", err
	}

	bodyStr := string(body)

	if isHTTPSuccess(response) {
		return bodyStr, nil
	}

	return bodyStr, fmt.Errorf("Status: %d", response.StatusCode)
}

func isHTTPSuccess(response *http.Response) bool {
	return response.StatusCode >= 200 && response.StatusCode < 300
}
