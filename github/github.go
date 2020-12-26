package github

// TODO paginated responses

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/soerenkoehler/simpson/util"
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
func NewDefaultContext() Context {
	return NewContext(os.Getenv("GITHUB_CONTEXT"))
}

// NewContext ...
func NewContext(jsonContext string) Context {
	context := Context{}
	json.Unmarshal([]byte(jsonContext), &context)
	return context
}

// APICall executes an Github API call on the given context, using the provided
// endpoint, content and URL-values.
func (context Context) APICall(
	endpoint Endpoint,
	content util.BodyReader,
	values ...interface{}) (string, error) {

	return context.APICallURL(
		endpoint.method,
		fmt.Sprintf(
			"https://api.github.com/repos/%s/%s",
			context.Repository,
			fmt.Sprintf(endpoint.url, values...)),
		content)
}

// APICallURL executes an API call on the given context using the provided
// method, url and content.
func (context Context) APICallURL(
	method string,
	url string,
	content util.BodyReader) (string, error) {

	fmt.Printf("API-Call: %v %v\n", method, url) // TODO DEBUG

	request, err := http.NewRequest(method, url, &content)
	if err != nil {
		return "", err
	}

	// request.Header.Set("Content-Length", fmt.Sprint(content.Length()))
	request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", context.Token))

	fmt.Printf("Headers: %v\n", request.Header) // TODO DEBUG

	response, err := httpClient.Do(request)
	if err != nil {
		return "", err
	}

	body, err := ioutil.ReadAll(response.Body)

	fmt.Printf("Error: %v\nResponse: %v\n", err, string(body)) // TODO DEBUG

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
