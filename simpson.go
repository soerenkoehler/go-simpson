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
		fmt.Printf("Arguments:\n%v\n", opts)
		artifacts := buildArtifacts(&opts)
		if optRelease, _ := opts.Bool("--release"); optRelease {
			releaseArtifacts(&opts, artifacts)
		}
	} else {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(127)
	}
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

func releaseArtifacts(opts *docopt.Opts, artifacts []string) {
	githubContext := github.NewDefaultContext()

	if len(githubContext.Token) > 0 {

		if isVersionTag(githubContext.Ref) {
			fmt.Fprintln(os.Stderr, "TODO: Not implemented!")
		} else if strings.HasPrefix(githubContext.Ref, "refs/heads/") {
			githubContext.SetTag("latest", githubContext.Sha)
			if release, err := githubContext.GetRelease("latest"); err == nil {
				uploadArtifacts(release, artifacts)
			} else {
				fmt.Fprintf(os.Stderr, "Skipping release: Release 'latest' not found.\n")
			}
		}

	} else {
		fmt.Fprintf(os.Stderr, "Skipping release: No Github API token found.\n")
	}
}

func uploadArtifacts(release github.ReleaseInfo, artifacts []string) {
	for _, artifact := range artifacts {
		if err := release.UploadArtifact(artifact); err != nil {
			fmt.Fprintf(os.Stderr, "Error uploading release asset %s: %v\n", artifact, err)
		}
	}
}

func isVersionTag(ref string) bool {
	return false // TODO
}
