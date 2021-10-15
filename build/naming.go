package build

import (
	"path/filepath"
	"strings"
	"time"

	"github.com/soerenkoehler/simpson/util"
)

// Build dates in several formats.
var (
	buildDate         = time.Now().UTC()
	buildDateLong     = buildDate.Format("2006.01.02-15:04:05")
	buildDateShort    = buildDate.Format("20060102-150405")
	TokenBuildDate    = "${BUILD_DATE}"
	TokenArtifactName = "${ARTIFACT_NAME}"
)

type NamingSpec struct {
	packageName         string
	artifactBaseName    string
	nameParts           []string
	executableExtension string
}

func NewNamingSpec(
	packageName string,
	artifactBaseName string) NamingSpec {

	realPackageName := packageName
	if packagePath, err := filepath.Abs(packageName); err == nil {
		realPackageName = filepath.Base(packagePath)
	}

	return NamingSpec{
		packageName:         realPackageName,
		artifactBaseName:    artifactBaseName,
		nameParts:           []string{TokenBuildDate},
		executableExtension: ""}
}

func (naming NamingSpec) GetVersionInfo() string {
	return util.ReplaceMultiple(
		strings.Join(naming.nameParts, " "),
		map[string]string{
			TokenBuildDate:    buildDateLong,
			TokenArtifactName: ""})
}

func (naming NamingSpec) GetArtifactName() string {
	return util.ReplaceMultiple(
		strings.Join(naming.nameParts, "-"),
		map[string]string{
			TokenBuildDate:    buildDateShort,
			TokenArtifactName: naming.artifactBaseName})
}

func (naming NamingSpec) GetArtifactFile() string {
	filename := naming.packageName
	if len(naming.artifactBaseName) > 0 {
		filename = naming.artifactBaseName
	}
	return filename + naming.executableExtension
}

func (naming NamingSpec) WithTarget(target TargetSpec) NamingSpec {
	return NamingSpec{
		packageName:         naming.packageName,
		artifactBaseName:    naming.artifactBaseName,
		nameParts:           append(naming.nameParts, target.Desc()),
		executableExtension: target.executableExtension}
}

func (naming NamingSpec) WithNameParts(newNameParts []string) NamingSpec {
	return NamingSpec{
		packageName:         naming.packageName,
		artifactBaseName:    naming.artifactBaseName,
		nameParts:           newNameParts,
		executableExtension: naming.executableExtension}
}
