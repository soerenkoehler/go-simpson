package main

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/docopt/docopt-go"
	"github.com/soerenkoehler/simpson/build"
	"github.com/soerenkoehler/simpson/github"
	"github.com/soerenkoehler/simpson/resource"
	"github.com/soerenkoehler/simpson/util"
)

var _Version = "DEV"

func main() {
	errors := []error{}
	defer func() {
		if len(errors) > 0 {
			for _, e := range errors {
				fmt.Fprintf(os.Stderr, "Error: %v\n", e)
			}
			os.Exit(1)
		}
	}()
	errors = doMain()
}

func doMain() []error {
	options, err := docopt.ParseArgs(
		util.ReplaceVariable(
			resource.Usage,
			"VERSION",
			_Version),
		nil,
		_Version)

	if err == nil {
		if hasOption(options, "--init") {
			err = initializeWorkflowFile()
		} else {
			githubContext := github.NewDefaultContext()
			artifacts, errs := build.TestAndBuild(
				getString(options, "PACKAGE"),
				githubContext.GetVersionLabels(),
				getTargets(options))
			if len(errs) == 0 {
				if githubContext.IsGithubAction() {
					errs = githubContext.CreateRelease(
						artifacts,
						hasOption(options, "--latest"))
				} else {
					fmt.Fprint(
						os.Stdout,
						"Skipping release: Must run in a Github action\n")
				}
			}
			return errs
		}
	}

	if err != nil {
		return []error{err}
	}

	return []error{}
}

func getTargets(options docopt.Opts) []build.TargetSpec {
	if hasOption(options, "--all-targets") {
		return build.AllTargets
	}
	targets, unknown := build.GetTargets(getString(options, "--targets"))
	if len(unknown) > 0 {
		fmt.Fprintf(os.Stderr, "Skipping unknown targets: %v\n", unknown)
	}
	return targets
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
	workflowFile := ".github/workflows/simpson-build-and-release-tool.yml"
	moduleInfo := util.FindInFile("go.mod", `^\s*module\s+(.+)$`)
	goInfo := util.FindInFile("go.mod", `^\s*go\s+(.+)$`)

	if len(moduleInfo) < 1 {
		return errors.New("go.mod: no module found")
	}

	if len(goInfo) < 1 {
		return errors.New("go.mod: no go version found")
	}

	err := os.MkdirAll(filepath.Dir(workflowFile), 0777)
	if err != nil {
		return err
	}

	output, err := os.Create(workflowFile)
	if err != nil {
		return err
	}

	defer output.Close()

	output.Write([]byte(
		util.ReplaceVariable(
			util.ReplaceVariable(
				resource.WorkflowFile,
				"SIMPSON_MODULE",
				moduleInfo[1]),
			"SIMPSON_GOVERSION",
			goInfo[1])))

	return nil
}
