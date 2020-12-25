package github

import (
	"fmt"
	"net/http"
	"path/filepath"

	"github.com/soerenkoehler/simpson/util"
)

// UploadArtifact ... TODO
func (release ReleaseInfo) UploadArtifact(path string) error {
	_, err := release.APICallURL(
		http.MethodPost,
		fmt.Sprintf("%s?name=%s", release.UploadURL, filepath.Base(path)),
		util.BodyFromFile(path))
	return err
}
