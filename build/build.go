package build

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/soerenkoehler/simpson/util"
)

const (
	artifactsParentDir = "artifacts"
	hashFileName       = "sha256.txt"
)

// TestAndBuild performs the standard build process.
func TestAndBuild(
	naming NamingSpec,
	targets []TargetSpec) ([]string, []error) {

	if err := Test(); err != nil {
		return []string{}, []error{err}
	}
	return Build(naming, targets)
}

// Test runs 'go test' for all packages in the current module.
func Test() error {
	if err := util.Execute([]string{"go", "vet", "./..."}); err != nil {
		return err
	}
	return util.Execute([]string{"go", "test", "--cover", "./..."})
}

// Build runs 'go build' for the named package and supplied target definitions.
// The resulting binary is stored in target specific subdirectories of the
// directory 'artifacts'.
func Build(
	naming NamingSpec,
	targets []TargetSpec) ([]string, []error) {

	os.RemoveAll(artifactsParentDir)
	os.Mkdir(artifactsParentDir, 0777)

	artifactList := []string{}
	hashList := []string{}
	errorList := []error{}

	for _, target := range targets {
		archivePath, archiveHash, err := buildArtifact(naming, target)
		if err == nil {
			artifactList = append(artifactList, archivePath)
			hashList = append(
				hashList,
				fmt.Sprintf("%s *%s", archiveHash, path.Base(archivePath)))
		} else {
			errorList = append(errorList, err)
		}
	}

	hashFilePath, err := writeHashFile(hashList)
	if err == nil {
		artifactList = append(artifactList, hashFilePath)
	} else {
		errorList = append(errorList, err)
	}

	return artifactList, errorList
}

func writeHashFile(hashList []string) (string, error) {
	hashFilePathRel := path.Join(artifactsParentDir, hashFileName)

	hashFilePathAbs, err := filepath.Abs(hashFilePathRel)
	if err != nil {
		return "", err
	}

	hashFile, err := os.Create(hashFilePathAbs)
	if err != nil {
		return "", err
	}

	defer hashFile.Close()
	_, err = hashFile.WriteString(strings.Join(hashList, "\n"))
	if err != nil {
		return "", err
	}

	return hashFilePathAbs, nil
}

func buildArtifact(
	naming NamingSpec,
	target TargetSpec) (string, string, error) {

	targetNaming := naming.WithTarget(target)
	artifactDir, artifactFile := createArtifactSubdir(targetNaming)

	if err := util.Execute(
		[]string{
			"go",
			"build",
			"-a",
			"-ldflags", fmt.Sprintf(
				// -s  omit symbols
				// -w  omit DWARF
				// -X  set string value
				`-s -w -X "main._Version=%v"`,
				targetNaming.GetVersionInfo()),
			"-o", artifactFile,
			targetNaming.packageName},
		target.Env()...); err != nil {
		return "", "", err
	}

	return util.CreateArchive(target.archiveType, artifactDir)
}

func createArtifactSubdir(naming NamingSpec) (string, string) {

	targetPath := path.Join(artifactsParentDir, naming.GetArtifactName())

	os.Mkdir(targetPath, 0777)

	artifactDir, err := filepath.Abs(targetPath)
	if err != nil {
		panic(err)
	}

	return artifactDir, path.Join(artifactDir, naming.GetArtifactFile())
}
