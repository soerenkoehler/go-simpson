package github

import (
	"encoding/json"
	"fmt"

	"github.com/soerenkoehler/simpson/util"
)

// ReleaseInfo ... TODO
type ReleaseInfo struct {
	context   Context
	ID        string `json:"id"`
	AssetsURL string `json:"assets_url"`
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
	response, err := release.context.APICall(
		APIUpdateRelease,
		util.BodyFromMap(map[string]string{
			"tag_name": tag,
		}),
		release.ID)
	return release.context.jsonToReleaseInfo(response), err
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
	result := ReleaseInfo{}
	result.context = context
	json.Unmarshal([]byte(jsonData), &result)
	fmt.Printf("ReleaseInfo: %v\n", result) // TODO DEBUG
	return result
}
