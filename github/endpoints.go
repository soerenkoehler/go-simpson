package github

import "net/http"

// apiEndpoint defines REST API endpoints.
type apiEndpoint struct {
	method string
	url    string
}

// Github API endpoints (not exhaustive)
var (
	apiGetRef    = apiEndpoint{http.MethodGet, "git/ref/%v"}
	apiCreateRef = apiEndpoint{http.MethodPost, "git/refs"}
	apiUpdateRef = apiEndpoint{http.MethodPatch, "git/refs/%v"}

	apiGetReleaseByTag = apiEndpoint{http.MethodGet, "releases/tags/%v"}
	apiCreateRelease   = apiEndpoint{http.MethodPost, "releases"}
	apiDeleteRelease   = apiEndpoint{http.MethodDelete, "releases/%v"}

	apiUploadReleaseAsset = apiEndpoint{http.MethodPost, "releases/%v/assets?name=%v"}
)
