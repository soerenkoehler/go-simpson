package github

import "net/http"

// Endpoint defines REST API endpoints.
type Endpoint struct {
	method string
	url    string
}

// Github API endpoints (not exhaustive)
var (
	APIGetRef    = &Endpoint{http.MethodGet, "git/ref/%s"}
	APICreateRef = &Endpoint{http.MethodPost, "git/refs"}
	APIUpdateRef = &Endpoint{http.MethodPatch, "git/refs/%s"}

	APIGetReleaseByTag = &Endpoint{http.MethodGet, "releases/tags/%s"}
	APICreateRelease   = &Endpoint{http.MethodPost, "releases"}
	APIUpdateRelease   = &Endpoint{http.MethodPatch, "releases/%s"}
)
