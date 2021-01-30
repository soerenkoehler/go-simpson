package main

import (
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
	retcode := 0
	defer func() {
		os.Exit(retcode)
	}()
	retcode = doMain()
}

func doMain() int {
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
			// TODO check for option --latest
			githubContext := github.NewDefaultContext()
			artifacts, errs := build.TestAndBuild(
				getString(options, "PACKAGE"),
				githubContext.GetVersionLabels(),
				getTargets(options))
			if len(errs) == 0 {
				if hasOption(options, "--release") {
					errs = githubContext.CreateRelease(artifacts)
				}
			} else {
				fmt.Fprintf(os.Stderr, "Errors:\n%v\n", errs)
				return 1
			}
		}
	} else {
		fmt.Fprintf(os.Stderr, "Arguments: %v\nError: %v\n", options, err)
		return 1
	}

	return 0
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

// TODO use go.mod or such for package
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
