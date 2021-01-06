package build

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/soerenkoehler/simpson/util"
)

const artifactsParentDir = "artifacts"

// Build dates in several formats.
var (
	buildDate      = time.Now().UTC()
	buildDateLong  = buildDate.Format("2006.01.02-15:04:05")
	buildDateShort = buildDate.Format("20060102-150405")
)

// TestAndBuild performs the standard build process.
func TestAndBuild(
	packageName string,
	versionLabels []string,
	targets []TargetSpec) ([]string, []error) {

	if Test() == nil {
		return Build(packageName, versionLabels, targets)
	}

	return []string{}, []error{}
}

// Test runs 'go test' for all packages in the current module.
func Test() error {
	return util.Execute([]string{"go", "test", "./..."})
}

// Build runs 'go build' for the named package and supplied target definitions.
// The resulting binary is stored in target specific subdirectories of the
// directory 'artifacts'.
func Build(
	packageName string,
	versionLabels []string,
	targets []TargetSpec) ([]string, []error) {

	os.RemoveAll(artifactsParentDir)

	artifactList := []string{}
	errorList := []error{}

	for _, target := range targets {
		if path, err := buildArtifact(
			packageName,
			target,
			versionLabels); err == nil {
			artifactList = append(artifactList, path)
		} else {
			errorList = append(errorList, err)
		}
	}

	return artifactList, errorList
}

func buildArtifact(
	packageName string,
	target TargetSpec,
	versionLabels []string) (string, error) {

	// TODO include git ref info
	targetLabels := append(versionLabels, target.Desc())
	artifactDir := createArtifactSubdir(packageName, targetLabels)

	if err := util.Execute(
		[]string{
			"go",
			"build",
			"-a",
			"-ldflags", fmt.Sprintf(
				`-X "main._Version=%s"`,
				strings.Join(targetLabels, " ")),
			"-o",
			artifactDir,
			packageName},
		target.Env()...); err != nil {
		return "", err
	}

	return util.CreateArchive(target.archiveType, artifactDir)
}

func createArtifactSubdir(
	packageName string,
	targetLabels []string) string {

	realPackageName := packageName
	if packagePath, err := filepath.Abs(packageName); err == nil {
		realPackageName = filepath.Base(packagePath)
	}

	targetDir := path.Join(
		artifactsParentDir,
		strings.Join(append([]string{realPackageName}, targetLabels...), "-"))

	os.MkdirAll(targetDir, 0777)

	result, err := filepath.Abs(targetDir)
	if err != nil {
		panic(err)
	}

	return result
}
