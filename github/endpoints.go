package github

import "net/http"

// apiEndpoint defines REST API endpoints.
type apiEndpoint struct {
	method string
	url    string
}

// Github API endpoints (not exhaustive)
var (
	apiGetRef    = apiEndpoint{http.MethodGet, "git/ref/%s"}
	apiCreateRef = apiEndpoint{http.MethodPost, "git/refs"}
	apiUpdateRef = apiEndpoint{http.MethodPatch, "git/refs/%s"}

	apiGetReleaseByTag = apiEndpoint{http.MethodGet, "releases/tags/%s"}
	apiCreateRelease   = apiEndpoint{http.MethodPost, "releases"}
	apiDeleteRelease   = apiEndpoint{http.MethodDelete, "releases/%s"}

	apiUploadReleaseAsset = apiEndpoint{http.MethodPost, "releases/%s/assets?name=%s"}
)
