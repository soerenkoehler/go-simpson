package github

// TODO paginated responses

import (
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

// SetTag tags creates or updates the given <tag> to the commit <sha>.
func (context *Context) SetTag(tag string, sha string) {
	// TODO
	body := map[string]string{
		"ref": fmt.Sprintf("refs/tags/%s", tag),
		"sha": sha,
	}
	if _, error := context.APICall(APIGetRef, "", fmt.Sprintf("tags/%s", tag)); error == nil {
		context.APICall(APIUpdateRef, body)
	} else {
		context.APICall(APICreateRef, body)
	}
}

// APICall executes an Github API on the given context, using the provided
// endpoint and values.
func (context *Context) APICall(
	endpoint *Endpoint,
	values ...interface{}) (string, error) {

	request, err := http.NewRequest(
		endpoint.method,
		fmt.Sprintf(
			"https://api.github.com/repos/%s/%s",
			context.Repository,
			fmt.Sprintf(endpoint.url, values...)),
		nil)
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

	return string(body), nil
}
