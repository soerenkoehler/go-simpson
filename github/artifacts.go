package github

import (
	"fmt"
	"io/ioutil"
)

// UploadArtifacts ... TODO
func (release *ReleaseInfo) UploadArtifacts(prefix string) {
	if artifacts, err := ioutil.ReadDir("artifacts"); err == nil {
		for _, artifact := range artifacts {
			release.uploadArtifact(artifact.Name())
		}
	}
}

func (release *ReleaseInfo) uploadArtifact(name string) {
	fmt.Printf("artifacts/%s => %s\n", name, release.AssetsURL)
}
