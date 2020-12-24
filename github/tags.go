package github

// TODO paginated responses

import (
	"fmt"

	"github.com/soerenkoehler/simpson/util"
)

// SetTag tags creates or updates the given <tag> to the commit <sha>.
func (context Context) SetTag(tag string, sha string) {
	if context.hasTag(tag) {
		context.updateTag(tag, sha)
	} else {
		context.createTag(tag, sha)
	}
}

func (context Context) hasTag(tag string) bool {
	_, err := context.APICall(APIGetRef, util.BodyReader{}, tagPath(tag))
	return err == nil
}

func (context Context) updateTag(tag string, sha string) error {
	_, err := context.APICall(
		APIUpdateRef,
		util.BodyFromMap(map[string]string{
			"sha": sha,
		}),
		tagPath(tag))
	return err
}

func (context Context) createTag(tag string, sha string) error {
	_, err := context.APICall(
		APICreateRef,
		util.BodyFromMap(map[string]string{
			"ref": fullTagPath(tag),
			"sha": sha,
		}))
	return err
}

func fullTagPath(tag string) string {
	return fmt.Sprintf("refs/%s", tagPath(tag))
}

func tagPath(tag string) string {
	return fmt.Sprintf("tags/%s", tag)
}
