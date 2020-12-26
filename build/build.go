package build

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"time"

	"github.com/soerenkoehler/simpson/util"
)

const artifactDir = "artifacts"

var buildDate = time.Now().UTC().Format("2006.01.02-15:04:05")

// TestAndBuild performs the standard build process.
func TestAndBuild(
	packageName string,
	targets []TargetSpec) ([]string, []error) {

	if Test() == nil {
		return Build(packageName, targets)
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
	targets []TargetSpec) ([]string, []error) {

	os.RemoveAll(artifactDir)

	artifactList := []string{}
	errorList := []error{}

	for _, target := range targets {
		if path, err := buildArtifact(
			packageName,
			target,
			artifactDir); err == nil {
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
	artifactDir string) (string, error) {

	// TODO include git ref info
	version := fmt.Sprintf("%s %s", buildDate, target.Desc())

	if err := util.Execute(
		[]string{
			"go",
			"build",
			"-a",
			"-ldflags", fmt.Sprintf("-X \"main._Version=%s\"", version),
			"-o",
			createArtifactSubdir(target, artifactDir),
			packageName},
		target.Env()...); err != nil {
		return "", err
	}

	return util.CreateArchive(
		target.archiveType,
		filepath.Join(
			artifactDir,
			target.Desc()))
}

func createArtifactSubdir(target TargetSpec, parent string) string {
	targetDir := path.Join(parent, target.Desc())
	os.MkdirAll(targetDir, 0777)
	result, err := filepath.Abs(targetDir)
	if err != nil {
		panic(err)
	}
	return result
}

func realPackageName(packageName string) string {
	if packagePath, err := filepath.Abs(packageName); err == nil {
		return filepath.Base(packagePath)
	}
	return packageName
}
