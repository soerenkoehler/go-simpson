package github

import (
	"fmt"

	"github.com/soerenkoehler/simpson/util"
)

func (context Context) setTag(tag string, sha string) {
	if context.hasTag(tag) {
		context.updateTag(tag, sha)
	} else {
		context.createTag(tag, sha)
	}
}

func (context Context) hasTag(tag string) bool {
	_, err := context.apiCall(apiGetRef, util.BodyReader{}, tagPath(tag))
	return err == nil
}

func (context Context) updateTag(tag string, sha string) error {
	_, err := context.apiCall(
		apiUpdateRef,
		util.BodyFromMap(map[string]string{
			"sha": sha,
		}),
		tagPath(tag))
	return err
}

func (context Context) createTag(tag string, sha string) error {
	_, err := context.apiCall(
		apiCreateRef,
		util.BodyFromMap(map[string]string{
			"ref": fullTagPath(tag),
			"sha": sha,
		}))
	return err
}

func fullTagPath(tag string) string {
	return fmt.Sprintf("refs/%v", tagPath(tag))
}

func tagPath(tag string) string {
	return fmt.Sprintf("tags/%v", tag)
}
