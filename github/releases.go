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

// ReleaseInfo ... TODO
type ReleaseInfo struct {
	Context
	ID        string
	Name      string
	UploadURL string
}

// CreateRelease ... TODO
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

// GetRelease ... TODO
func (context Context) getRelease(tag string) (ReleaseInfo, error) {
	release, err := context.getReleaseByTag(tag)
	if err == nil {
		release, err = release.updateRelease(tag, release.Name)
	} else {
		release, err = context.createRelease(tag, tag)
	}
	return release, err
}

func (context Context) getReleaseByTag(tag string) (ReleaseInfo, error) {
	response, err := context.apiCall(apiGetReleaseByTag, util.BodyReader{}, tag)
	if err != nil {
		return ReleaseInfo{}, err
	}
	fmt.Printf("Release:\n%v\n", context.jsonToReleaseInfo(response)) // DEBUG
	return context.jsonToReleaseInfo(response), nil
}

func (release ReleaseInfo) updateRelease(
	tag string,
	name string) (ReleaseInfo, error) {

	fmt.Printf("Update Release:\n%v %v\n", tag, name) // DEBUG
	if _, err := release.apiCall(
		apiDeleteRelease,
		util.BodyReader{},
		release.ID); err != nil {
		return ReleaseInfo{}, err
	}
	return release.createRelease(tag, name)
}

func (context Context) createRelease(
	tag string,
	name string) (ReleaseInfo, error) {

	fmt.Printf("Create Release:\n%v %v\n", tag, name) // DEBUG
	response, err := context.apiCall(
		apiCreateRelease,
		util.BodyFromMap(map[string]string{
			"tag_name": tag,
			"name":     name,
		}))
	return context.jsonToReleaseInfo(response), err
}

func (context Context) jsonToReleaseInfo(jsonData string) ReleaseInfo {
	var result map[string]interface{}
	json.Unmarshal([]byte(jsonData), &result)
	return ReleaseInfo{
		Context: context,
		ID:      fmt.Sprintf("%.f", result["id"]),
		UploadURL: uploadURLNormalizer.ReplaceAllString(
			result["upload_url"].(string), "")}
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
		fmt.Sprintf("%s?name=%s", release.UploadURL, filepath.Base(path)),
		util.BodyFromFile(path))
	return err
}
