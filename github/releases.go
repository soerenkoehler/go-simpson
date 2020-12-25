package github

import (
	"encoding/json"
	"fmt"
	"regexp"

	"github.com/soerenkoehler/simpson/util"
)

var uploadURLNormalizer = regexp.MustCompile(`\{\?[\w,]+\}$`)

// ReleaseInfo ... TODO
type ReleaseInfo struct {
	Context
	ID        string
	UploadURL string
}

// GetRelease ... TODO
func (context Context) GetRelease(tag string) (ReleaseInfo, error) {
	release, err := context.getReleaseByTag(tag)
	if err == nil {
		release, err = release.updateRelease(tag)
	} else {
		release, err = context.createRelease(tag)
	}
	return release, err
}

func (context Context) getReleaseByTag(tag string) (ReleaseInfo, error) {
	response, err := context.APICall(APIGetReleaseByTag, util.BodyReader{}, tag)
	if err != nil {
		return ReleaseInfo{}, err
	}
	return context.jsonToReleaseInfo(response), nil
}

func (release ReleaseInfo) updateRelease(tag string) (ReleaseInfo, error) {
	response, err := release.APICall(
		APIUpdateRelease,
		util.BodyFromMap(map[string]string{
			"tag_name": tag,
		}),
		release.ID)
	return release.jsonToReleaseInfo(response), err
}

func (context Context) createRelease(tag string) (ReleaseInfo, error) {
	response, err := context.APICall(
		APICreateRelease,
		util.BodyFromMap(map[string]string{
			"tag_name": tag,
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
