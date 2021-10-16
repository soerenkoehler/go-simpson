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

const tagLatest = "latest"

var uploadURLNormalizer = regexp.MustCompile(`\{\?[\w,]+\}$`)

type ReleaseInfo struct {
	Context   `json:"-"`
	ID        int64  `json:"id"`
	Name      string `json:"name"`
	UploadURL string `json:"upload_url"`
}

func (context Context) CreateRelease(artifacts []string) []error {
	var errs []error

	if len(context.Token) == 0 {
		//lint:ignore ST1005 Github is a proper noun
		errs = append(errs, errors.New("Github API token not found"))
	} else if release, err := context.getRelease(); err != nil {
		errs = append(errs, err)
	} else {
		for _, artifact := range artifacts {
			if err := release.uploadArtifact(artifact); err != nil {
				errs = append(errs, err)
			}
		}
	}

	return errs
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

func (context Context) getRelease() (ReleaseInfo, error) {

	tag := ""
	if version, ok := context.getPushVersion(); ok {
		tag = version
	} else if context.isPushHead() {
		context.setTag(tagLatest, context.Sha)
		tag = tagLatest
	} else {
		return ReleaseInfo{}, errors.New("pushed neither version tag nor head ref")
	}

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
