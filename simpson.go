package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/docopt/docopt-go"
	"github.com/soerenkoehler/simpson/build"
	"github.com/soerenkoehler/simpson/github"
)

func main() {
	opts, err := docopt.ParseArgs(_Usage, nil, _Version)
	if err == nil {
		fmt.Println(opts) // TODO debug or info?
		buildArtifacts(&opts)
		if optRelease, _ := opts.Bool("--release"); optRelease {
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
	}
	optPackage, _ := opts.String("PACKAGE")
	build.TestAndBuild(optPackage, targets)
}

func releaseArtifacts(opts *docopt.Opts) {
	githubContext := github.NewDefaultContext()
	fmt.Println(githubContext)
	if len(githubContext.Token) > 0 {
		if isVersionTag(githubContext.Ref) {
			fmt.Println("TODO: Not implemented!")
		} else if strings.HasPrefix(githubContext.Ref, "refs/heads/") {
			githubContext.SetTag("latest", githubContext.Sha)
			release := githubContext.GetRelease("latest")
			optPackage, _ := opts.String("PACKAGE")
			release.UploadArtifacts(optPackage)
		}
	} else {
		fmt.Fprintln(os.Stderr, "Error: No Github API token found.")
	}
}

func isVersionTag(ref string) bool {
	return false // TODO
}
