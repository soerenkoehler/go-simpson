package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/docopt/docopt-go"
	"github.com/soerenkoehler/simpson/build"
	"github.com/soerenkoehler/simpson/github"
	"github.com/soerenkoehler/simpson/resource"
	"github.com/soerenkoehler/simpson/util"
)

var _Version = "DEV"

func main() {
	opts, err := docopt.ParseArgs(
		util.ReplaceVariable(resource.Usage, "VERSION", _Version),
		nil,
		_Version)

	if err == nil {
		fmt.Printf("Arguments:\n%v\n", opts)
		if optInit, _ := opts.Bool("--init"); optInit {
			err = initializeWorkflowFile()
		} else {
			artifacts := buildArtifacts(&opts)
			if optRelease, _ := opts.Bool("--release"); optRelease {
				createRelease(&opts, artifacts)
			}
		}
	}

	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(127)
	}
}

func initializeWorkflowFile() error {
	workflowFile := ".github/workflows/simpson-bart.yml"

	err := os.MkdirAll(filepath.Dir(workflowFile), 0777)
	if err != nil {
		return err
	}

	output, err := os.Create(workflowFile)
	if err != nil {
		return err
	}

	defer output.Close()

	output.Write([]byte(resource.WorkflowFile))

	return nil
}

func buildArtifacts(opts *docopt.Opts) []string {
	var targets []build.TargetSpec

	if optAllTargets, _ := opts.Bool("--all-targets"); optAllTargets {
		targets = build.AllTargets
	} else if optTargets, err := opts.String("--targets"); err == nil {
		targets = build.GetTargets(optTargets)
	}
	optPackage, _ := opts.String("PACKAGE")

	artifacts, errs := build.TestAndBuild(optPackage, targets)
	if len(errs) > 0 {
		fmt.Fprintf(os.Stderr, "Errors:\n%v\n", errs)
	}

	return artifacts
}

func createRelease(opts *docopt.Opts, artifacts []string) {
	githubContext := github.NewDefaultContext()

	if len(githubContext.Token) > 0 {

		if isVersionTag(githubContext.Ref) {
			fmt.Fprintln(os.Stderr, "TODO: Not implemented!")
		} else if strings.HasPrefix(githubContext.Ref, "refs/heads/") {
			githubContext.SetTag("latest", githubContext.Sha)
			if release, err := githubContext.GetRelease("latest"); err == nil {
				artifacts = renameArtifacts(
					fmt.Sprintf(
						"%s-%s",
						build.BuildDateShort,
						githubContext.Sha[0:7]),
					artifacts)
				uploadArtifacts(release, artifacts)
			} else {
				fmt.Fprintf(
					os.Stderr,
					"Skipping release: Release 'latest' not found.\n")
			}
		}

	} else {
		fmt.Fprintf(os.Stderr, "Skipping release: No Github API token found.\n")
	}
}

// TODO include package name and version suffix in artifact path
func renameArtifacts(
	prefix string,
	suffix string,
	artifacts []string) []string {

	var newNames []string

	for _, artifact := range artifacts {
		newName := fmt.Sprintf("%s-%s-%s", artifact, suffix)
		os.Rename(artifact, newName)
		newNames = append(newNames, newName)
	}

	return newNames
}

func uploadArtifacts(release github.ReleaseInfo, artifacts []string) {
	for _, artifact := range artifacts {
		if err := release.UploadArtifact(artifact); err != nil {
			fmt.Fprintf(
				os.Stderr,
				"Error uploading release asset %s: %v\n",
				artifact,
				err)
		}
	}
}

func isVersionTag(ref string) bool {
	return false // TODO
}
