package main

import (
	"fmt"
	"os"

	"github.com/docopt/docopt-go"
	"github.com/soerenkoehler/bar-go/build"
)

var _Version = "2020-10-31 23:44:34"
var _Usage = "Build And Release for Go (" + _Version + `)

Usage:
    bar-go PACKAGE [--all-targets | --targets TARGETS]
    bar-go (-h | --help | --version)

Options:
    --all-targets     Build all possible targets
    --targets TARGETS Build the given targets (comma seperated list).

    -h --help  Show help.
    --version  Show version.`

func main() {
	opts, err := docopt.ParseArgs(_Usage, nil, _Version)
	if err == nil {
		fmt.Println(opts)
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
	} else {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(127)
	}
}
