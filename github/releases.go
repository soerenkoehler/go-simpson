package github

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"path/filepath"
	"regexp"

	"github.com/soerenkoehler/simpson/util"
)

var uploadURLNormalizer = regexp.MustCompile(`\{\?[\w,]+\}$`)

type ReleaseInfo struct {
	Context   `json:"-"`
	ID        int64  `json:"id"`
	Name      string `json:"name"`
	UploadURL string `json:"upload_url"`
}

func (context Context) CreateRelease(
	artifacts []string,
	doLatest bool) []error {

	if len(context.Token) > 0 {
		if version, ok := context.getPushVersion(); ok {
			return context.uploadArtifacts(version, artifacts)
		} else if doLatest && context.isPushHead() {
			context.setTag("latest", context.Sha)
			return context.uploadArtifacts("latest", artifacts)
		}
		return []error{errors.New("Pushed neither version tag nor head ref")}
	}
	return []error{errors.New("Github API token not found")}
}

func (context Context) uploadArtifacts(
	releaseName string,
	artifacts []string) []error {

	if release, err := context.getRelease(releaseName); err == nil {
		var errs []error
		for _, artifact := range artifacts {
			if err := release.uploadArtifact(artifact); err != nil {
				errs = append(errs, err)
			}
		}
		return errs
	}
	return []error{fmt.Errorf("Release '%v' not found", releaseName)}
}

func (release ReleaseInfo) uploadArtifact(path string) error {
	_, err := release.apiCallURL(
		http.MethodPost,
		fmt.Sprintf(
			"%v?name=%v",
			uploadURLNormalizer.ReplaceAllString(release.UploadURL, ""),
			filepath.Base(path)),
		util.BodyFromFile(path))
	return err
}

func (context Context) getRelease(tag string) (ReleaseInfo, error) {

	release, err := context.getReleaseByTag(tag)

	if err == nil {
		return release.updateRelease(tag, release.Name)
	}

	return context.createRelease(tag, tag)
}

func (context Context) getReleaseByTag(tag string) (ReleaseInfo, error) {

	response, err := context.apiCall(apiGetReleaseByTag, util.BodyReader{}, tag)

	if err != nil {
		return ReleaseInfo{}, err
	}

	return context.jsonToReleaseInfo(response)
}

func (release ReleaseInfo) updateRelease(
	tag string,
	name string) (ReleaseInfo, error) {

	_, err := release.apiCall(apiDeleteRelease, util.BodyReader{}, release.ID)

	if err != nil {
		return ReleaseInfo{}, err
	}

	return release.createRelease(tag, name)
}

func (context Context) createRelease(
	tag string,
	name string) (ReleaseInfo, error) {

	response, err := context.apiCall(
		apiCreateRelease,
		util.BodyFromMap(map[string]string{
			"tag_name": tag,
			"name":     name,
		}))

	if err != nil {
		return ReleaseInfo{}, err
	}

	return context.jsonToReleaseInfo(response)
}

func (context Context) jsonToReleaseInfo(jsonData string) (ReleaseInfo, error) {

	result := ReleaseInfo{
		Context: context,
	}

	err := json.Unmarshal([]byte(jsonData), &result)

	return result, err
}
