package github

import (
	"encoding/json"
	"fmt"
)

// ReleaseInfo ... TODO
type ReleaseInfo struct {
	context   *Context
	ID        string `json:"id"`
	AssetsURL string `json:"assets_url"`
}

// GetRelease ... TODO
func (context *Context) GetRelease(tag string) *ReleaseInfo {
	release, err := context.getReleaseByTag(tag)
	if err == nil {
		release, err = release.updateRelease(tag)
	} else {
		release, err = context.createRelease(tag)
	}
	if err != nil {
		fmt.Println(err)
	}
	return release
}

func (context *Context) getReleaseByTag(tag string) (*ReleaseInfo, error) {
	response, err := context.APICall(APIGetReleaseByTag, nil, tag)
	if err != nil {
		return nil, err
	}
	return context.jsonToReleaseInfo(response), nil
}

func (release *ReleaseInfo) updateRelease(tag string) (*ReleaseInfo, error) {
	body, _ := json.Marshal(map[string]string{
		"tag_name": tag,
	})
	response, err := release.context.APICall(APIUpdateRelease, body, release.ID)
	fmt.Printf("Update release %s\nResult: %s\nError: %v\n", tag, response, err)
	return release.context.jsonToReleaseInfo(response), err
}

func (context *Context) createRelease(tag string) (*ReleaseInfo, error) {
	body, _ := json.Marshal(map[string]string{
		"tag_name": tag,
	})
	response, err := context.APICall(APICreateRelease, body)
	fmt.Printf("Create release %s\nResult: %s\nError: %v\n", tag, response, err)
	return context.jsonToReleaseInfo(response), err
}

func (context *Context) jsonToReleaseInfo(jsonData string) *ReleaseInfo {
	result := &ReleaseInfo{}
	result.context = context
	json.Unmarshal([]byte(jsonData), result)
	return result
}
