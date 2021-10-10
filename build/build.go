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

const (
	artifactsParentDir = "artifacts"
	hashFileName       = "sha256.txt"
)

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
	artifactName string,
	targets []TargetSpec,
	versionLabels []string) ([]string, []error) {

	if err := Test(); err != nil {
		return []string{}, []error{err}
	}
	return Build(packageName, artifactName, targets, versionLabels)
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
	packageName string,
	artifactName string,
	targets []TargetSpec,
	versionLabels []string) ([]string, []error) {

	os.RemoveAll(artifactsParentDir)
	os.Mkdir(artifactsParentDir, 0777)

	artifactList := []string{}
	hashList := []string{}
	errorList := []error{}

	for _, target := range targets {
		archivePath, archiveHash, err := buildArtifact(
			packageName,
			artifactName,
			target,
			versionLabels)
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
	packageName string,
	artifactName string,
	target TargetSpec,
	versionLabels []string) (string, string, error) {

	targetLabels := append(versionLabels, target.Desc())
	artifactDir, artifactFile := createArtifactSubdir(
		packageName,
		artifactName,
		target,
		targetLabels)

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
				formatTargetLabels(targetLabels, buildDateLong, " ")),
			"-o", artifactFile,
			packageName},
		target.Env()...); err != nil {
		return "", "", err
	}

	return util.CreateArchive(target.archiveType, artifactDir)
}

func createArtifactSubdir(
	packageName string,
	artifactName string,
	target TargetSpec,
	targetLabels []string) (string, string) {

	realPackageName := artifactName
	if len(realPackageName) == 0 {
		realPackageName = packageName
		if packagePath, err := filepath.Abs(packageName); err == nil {
			realPackageName = filepath.Base(packagePath)
		}
	}

	targetDir := formatTargetLabels(
		append([]string{realPackageName}, targetLabels...),
		buildDateShort,
		"-")

	targetPath := path.Join(artifactsParentDir, targetDir)

	os.Mkdir(targetPath, 0777)

	artifactDir, err := filepath.Abs(targetPath)
	if err != nil {
		panic(err)
	}

	artifactFile := realPackageName
	if len(artifactName) > 0 {
		artifactFile = artifactName
	}
	if len(target.executableExtension) > 0 {
		artifactFile += "." + target.executableExtension
	}

	return artifactDir, path.Join(artifactDir, artifactFile)
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
