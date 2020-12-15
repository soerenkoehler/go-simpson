package github

// TODO paginated responses

import (
	"encoding/json"
	"fmt"
)

// SetTag tags creates or updates the given <tag> to the commit <sha>.
func (context *Context) SetTag(tag string, sha string) {
	if context.hasTag(tag) {
		context.updateTag(tag, sha)
	} else {
		context.createTag(tag, sha)
	}
}

func (context *Context) hasTag(tag string) bool {
	_, err := context.APICall(APIGetRef, nil, tagPath(tag))
	return err == nil
}

func (context *Context) updateTag(tag string, sha string) {
	body, _ := json.Marshal(map[string]string{
		"sha": sha,
	})
	result, err := context.APICall(APIUpdateRef, body, tagPath(tag))
	fmt.Printf("Update tag %s\nResult: %s\nError: %v\n", tag, result, err)
}

func (context *Context) createTag(tag string, sha string) {
	body, _ := json.Marshal(map[string]string{
		"ref": fullTagPath(tag),
		"sha": sha,
	})
	result, err := context.APICall(APICreateRef, body)
	fmt.Printf("Create tag %s\nResult: %s\nError: %v\n", tag, result, err)
}

func fullTagPath(tag string) string {
	return fmt.Sprintf("refs/%s", tagPath(tag))
}

func tagPath(tag string) string {
	return fmt.Sprintf("tags/%s", tag)
}
