package github

// TODO paginated responses

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"strings"

	"github.com/soerenkoehler/simpson/build"
	"github.com/soerenkoehler/simpson/util"
)

var pushVersionExtractor = regexp.MustCompile(`^refs/tags/(v\d+\.\d+\.\d+)`)
var httpClient *http.Client = &http.Client{}

// Context of current Github Actions workflow call.
type Context struct {
	Token      string
	Repository string
	Ref        string
	Sha        string
}

// NewDefaultContext ... TODO
func NewDefaultContext() Context {
	return NewContext(os.Getenv("GITHUB_CONTEXT"))
}

// NewContext ... TODO
func NewContext(jsonContext string) Context {
	context := Context{}
	json.Unmarshal([]byte(jsonContext), &context)
	return context
}

// IsGithubAction ... TODO
func (context Context) IsGithubAction() bool {
	return len(context.Token) > 0
}

// GetVersionLabels ... TODO
func (context Context) GetVersionLabels() []string {
	if pushVersion, ok := context.getPushVersion(); ok {
		return []string{pushVersion}
	} else if context.isPushHead() {
		return []string{build.TokenBuildDate, context.Sha[0:8]}
	}
	return []string{build.TokenBuildDate}
}

func (context Context) getPushVersion() (string, bool) {
	matches := pushVersionExtractor.FindStringSubmatch(context.Ref)
	if len(matches) == 2 {
		return matches[1], true
	}
	return "", false
}

func (context Context) isPushHead() bool {
	return strings.HasPrefix(context.Ref, "refs/heads/")
}

func (context Context) apiCall(
	endpoint apiEndpoint,
	content util.BodyReader,
	values ...interface{}) (string, error) {

	return context.apiCallURL(
		endpoint.method,
		fmt.Sprintf(
			"https://api.github.com/repos/%s/%s",
			context.Repository,
			fmt.Sprintf(endpoint.url, values...)),
		content)
}

func (context Context) apiCallURL(
	method string,
	url string,
	content util.BodyReader) (string, error) {

	request, err := http.NewRequest(method, url, &content)
	if err != nil {
		return "", err
	}

	request.ContentLength = content.Length()
	// request.Header.Set("Content-Length", fmt.Sprint(content.Length()))
	request.Header.Set("Content-Type", "application/octet-stream")
	request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", context.Token))

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
