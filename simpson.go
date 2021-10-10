package main

import (
	_ "embed"
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

type commandLine struct {
	Package      string `arg:"" default:"." help:"The package to compile. Default: '.'"`
	ArtifactName string `name:"artifact-name" help:"An alternate base name for the artifact files."`

	AllTargets bool     `name:"all-targets" short:"a" xor:"targets" help:"Build all possible targets"`
	Targets    []string `name:"targets" short:"t" xor:"targets" help:"Build the given targets."`

	Latest     bool `name:"latest" short:"l" help:"Tags the latest commit and creates a release named 'latest'."`
	SkipUpload bool `name:"skip-upload" short:"" help:"Build artifacts but do not upload them to the release."`

	Init bool `name:"init" short:"i" help:"Creates a Github Action file using the current commandline."`
}

func (cli commandLine) Validate() error {
	if len(cli.Targets) == 0 && !cli.AllTargets && !cli.SkipUpload {
		return fmt.Errorf("requires --skip-upload or one of --all-targets, --targets")
	}
	return nil
}

func (cli commandLine) getTargets() []build.TargetSpec {
	if cli.AllTargets {
		return build.AllTargets
	}
	targets, unknown := build.GetTargets(cli.Targets)
	if len(unknown) > 0 {
		logInfo("Skipping unknown targets: %v", unknown)
	}
	return targets
}

func main() {
	if err := doMain(); err != nil {
		logError("", err)
		os.Exit(1)
	}
}

func doMain() error {
	cli := commandLine{}
	kong.Parse(
		&cli,
		kong.Vars{"VERSION": _Version},
		kong.Description(_Description))

	if cli.Init {
		return initializeWorkflowFile()
	}

	githubContext := github.NewDefaultContext()

	artifacts, errs := build.TestAndBuild(
		cli.Package,
		cli.ArtifactName,
		cli.getTargets(),
		githubContext.GetVersionLabels())

	if len(errs) == 0 {
		if githubContext.IsGithubAction() {
			if cli.SkipUpload {
				artifacts = []string{}
				logInfo("found option --skip-upload: skipping artifact upload", nil)
			}
			errs = githubContext.CreateRelease(artifacts, cli.Latest)
		} else {
			logInfo("missing Github action context: skipping release creation", nil)
		}
	}

	if len(errs) > 1 {
		return fmt.Errorf("multiple errors: %v", errs)
	} else if len(errs) == 1 {
		return errs[0]
	}

	return nil
}

func initializeWorkflowFile() error {
	workflowFile := ".github/workflows/simpson-build-and-release-tool.yml"

	goInfo := util.FindInFile("go.mod", `^\s*go\s+(.+)$`)

	cmdline := strings.ReplaceAll(
		strings.Join(os.Args[1:], " "),
		" --init",
		"")

	if len(goInfo) < 1 {
		return fmt.Errorf("go.mod: no go version found")
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

func logInfo(message string, params interface{}) {
	logOutput(os.Stdout, "INFO", message, params)
}

func logError(message string, params interface{}) {
	logOutput(os.Stderr, "ERROR", message, params)
}

func logOutput(
	output *os.File,
	category string,
	message string,
	params interface{}) {
	fmt.Fprintf(output, "[%s] %s %v\n", category, message, params)
}
