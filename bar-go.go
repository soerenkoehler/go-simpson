package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/docopt/docopt-go"
	"github.com/soerenkoehler/bar-go/build"
	"github.com/soerenkoehler/bar-go/github"
)

var _Version = "2020-10-31 23:44:34"
var _Usage = "Build And Release for Go (" + _Version + `)

Usage:
    bar-go PACKAGE [--all-targets | --targets TARGETS] [--release]
    bar-go (-h | --help | --version)

Options:
    --all-targets      Build all possible targets
    --targets TARGETS  Build the given targets (comma seperated list).

    --release  Creates a tagged release and uploads the current artifacts.
               - Must run in a Github action (push event).
               - Pushing a tag with a version number (like 'v0.0.0') will
                 create a production release with that version number.
               - Updating a HEAD ref creates a release named 'latest'.

    -h --help  Show help.
    --version  Show version.`

func main() {
	opts, err := docopt.ParseArgs(_Usage, nil, _Version)
	if err == nil {
		fmt.Println(opts)
		buildArtifacts(&opts)
		if optDeploy, _ := opts.Bool("--release"); optDeploy {
			releaseArtifacts(&opts)
		}
	} else {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(127)
	}
}

func buildArtifacts(opts *docopt.Opts) {
	var targets []*build.TargetSpec
	if optAllTargets, _ := opts.Bool("--all-targets"); optAllTargets {
		targets = build.AllTargets
	} else if optTargets, err := opts.String("--targets"); err == nil {
		targets = build.GetTargets(optTargets)
	} else {
		targets = build.DefaultTargets
	}
	optPackage, _ := opts.String("PACKAGE")
	build.TestAndBuild(optPackage, targets)
}

func releaseArtifacts(opts *docopt.Opts) {
	githubContext := github.NewDefaultContext()
	fmt.Println(githubContext)
	if len(githubContext.Token) > 0 {
		if isVersionTag(githubContext.Ref) {
			fmt.Println("Not implemented!")
		} else if strings.HasPrefix(githubContext.Ref, "refs/heads/") {
			githubContext.SetTag(githubContext.Sha, "latest")
			// TODO create release
		}
	} else {
		fmt.Fprintln(os.Stderr, "Error: No Github API token found.")
	}
}

func isVersionTag(ref string) bool {
	return false // TODO
}
