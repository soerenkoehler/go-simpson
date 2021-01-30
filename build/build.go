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
	TokenBuildDate = "$BUILDDATE"
)

// TestAndBuild performs the standard build process.
func TestAndBuild(
	packageName string,
	versionLabels []string,
	targets []TargetSpec) ([]string, []error) {

	if err := Test(); err != nil {
		return []string{}, []error{err}
	}
	return Build(packageName, versionLabels, targets)
}

// Test runs 'go test' for all packages in the current module.
func Test() error {
	if err := util.Execute([]string{"go", "vet", "./..."}); err != nil {
		return err
	}
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
				formatTargetLabels(targetLabels, buildDateLong, " ")),
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

	targetDir := formatTargetLabels(
		append([]string{realPackageName}, targetLabels...),
		buildDateShort,
		"-")

	targetPath := path.Join(artifactsParentDir, targetDir)

	os.MkdirAll(targetPath, 0777)

	result, err := filepath.Abs(targetPath)
	if err != nil {
		panic(err)
	}

	return result
}

func formatTargetLabels(
	targetLabels []string,
	date string,
	separator string) string {

	result := []string{}
	for _, label := range targetLabels {
		result = append(result, replaceBuildDate(label, date))
	}

	return strings.Join(result, separator)
}

func replaceBuildDate(label string, date string) string {
	if label == TokenBuildDate {
		return date
	}
	return label
}
