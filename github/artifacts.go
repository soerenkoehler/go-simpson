package github

import (
	"fmt"
)

// UploadArtifact ... TODO
func (release ReleaseInfo) UploadArtifact(name string) {
	fmt.Printf("artifacts/%s => %s\n", name, release.AssetsURL)
}
