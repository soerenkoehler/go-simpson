package github

import (
	"path/filepath"

	"github.com/soerenkoehler/simpson/util"
)

// UploadArtifact ... TODO
func (release ReleaseInfo) UploadArtifact(path string) error {
	_, err := release.APICall(
		APIUploadReleaseAsset,
		util.BodyFromFile(path),
		release.ID,
		filepath.Base(path))
	return err
}
