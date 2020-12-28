package main

import (
	"errors"
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
	options, err := docopt.ParseArgs(
		util.ReplaceVariable(
			resource.Usage,
			"VERSION",
			_Version),
		nil,
		_Version)

	if err == nil {
		if hasOption(options, "--init") {
			if err := initializeWorkflowFile(); err != nil {
				fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			}
		} else {
			versionLabels := []string{} // TODO
			artifacts, errs := build.TestAndBuild(
				getString(options, "PACKAGE"),
				versionLabels,
				getTargets(options))
			if len(errs) == 0 {
				if hasOption(options, "--release") {
					errs = createRelease(artifacts)
				}
			} else {
				fmt.Fprintf(os.Stderr, "Errors:\n%v\n", errs)
			}
		}
	} else {
		fmt.Fprintf(os.Stderr, "Arguments: %v\nError: %v\n", options, err)
	}
}

func getTargets(options docopt.Opts) []build.TargetSpec {
	if hasOption(options, "--all-targets") {
		return build.AllTargets
	}
	return build.GetTargets(getString(options, "--targets"))
}

func hasOption(options docopt.Opts, name string) bool {
	result, err := options.Bool(name)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Checking Option: %v\nError: %v\n", name, err)
	}
	return result
}

func getString(options docopt.Opts, name string) string {
	result, err := options.String(name)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Get Value: %v\nError: %v\n", name, err)
	}
	return result
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

func createRelease(artifacts []string) []error {
	githubContext := github.NewDefaultContext()

	if len(githubContext.Token) > 0 {

		if isVersionTag(githubContext.Ref) {
			return []error{errors.New("TODO: Not implemented")}

		} else if strings.HasPrefix(githubContext.Ref, "refs/heads/") {

			githubContext.SetTag("latest", githubContext.Sha)
			if release, err := githubContext.GetRelease("latest"); err == nil {
				var errs []error
				for _, artifact := range artifacts {
					if err := release.UploadArtifact(artifact); err != nil {
						errs = append(errs, err)
					}
				}
				return errs
			}
			return []error{errors.New("Release 'latest' not found")}
		}
		return []error{errors.New("Invalid Github Context")}
	}
	return []error{errors.New("Github API token not found")}
}

func isVersionTag(ref string) bool {
	return false // TODO
}
