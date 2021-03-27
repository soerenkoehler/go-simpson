package main

import (
	_ "embed"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/alecthomas/kong"
	"github.com/soerenkoehler/simpson/build"
	"github.com/soerenkoehler/simpson/github"
	"github.com/soerenkoehler/simpson/util"
)

//go:embed resource/workflowfile.yml
var workflowFileTemplate string

//go:embed resource/description.txt
var _Description string

var _Version = "DEV"

var cli struct {
	Package      string   `arg:""`
	ArtifactName string   `name:"artifact-name" help:"An alternate base name for the artifact files."`
	AllTargets   bool     `name:"all-targets" short:"a" xor:"targets" help:"Build all possible targets"`
	Targets      []string `name:"targets" short:"t" xor:"targets" help:"Build the given targets (comma seperated list)."`
	Latest       bool     `name:"latest" short:"l" help:"Tags the latest commit and creates a release named 'latest'."`
	SkipUpload   bool     `name:"skip-upload" short:"" help:"Build artifacts but do not upload them to the release."`
	Init         bool     `name:"init" short:"i" help:"Creates a Github Action file using the current commandline."`
}

func main() {
	if err := doMain(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func doMain() error {
	kong.Parse(
		&cli,
		kong.Vars{"VERSION": _Version},
		kong.Description(_Description),
		kong.UsageOnError())

	if cli.Init {
		return initializeWorkflowFile()
	} else {
		githubContext := github.NewDefaultContext()
		artifacts, errs := build.TestAndBuild(
			cli.Package,
			cli.ArtifactName,
			getTargets(),
			githubContext.GetVersionLabels())
		if len(errs) == 0 {
			if githubContext.IsGithubAction() {
				if cli.SkipUpload {
					artifacts = []string{}
				}
				errs = githubContext.CreateRelease(artifacts, cli.Latest)
			} else {
				fmt.Fprint(
					os.Stdout,
					"Skipping release: Must run in a Github action\n")
			}
		}
		if len(errs) != 0 {
			return fmt.Errorf("multiple errors: %v", errs)
		}
		return nil
	}
}

func getTargets() []build.TargetSpec {
	if cli.AllTargets {
		return build.AllTargets
	}
	targets, unknown := build.GetTargets(cli.Targets)
	if len(unknown) > 0 {
		fmt.Fprintf(os.Stderr, "Skipping unknown targets: %v\n", unknown)
	}
	return targets
}

func initializeWorkflowFile() error {
	workflowFile := ".github/workflows/simpson-build-and-release-tool.yml"

	goInfo := util.FindInFile("go.mod", `^\s*go\s+(.+)$`)

	cmdline := strings.ReplaceAll(
		strings.Join(os.Args[1:], " "),
		" --init",
		"")

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
		util.ReplaceVariables(
			workflowFileTemplate,
			map[string]string{
				"SIMPSON_CMDLINE":   cmdline,
				"SIMPSON_GOVERSION": goInfo[1]})))

	return nil
}
