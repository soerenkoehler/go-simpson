package main

var _Version = "DEV"
var _Usage = "BaRT (Build and Release Tool) " + _Version + `

Usage:
    simpson PACKAGE (--all-targets | --targets TARGETS) [--release]
    simpson --init
    simpson (-h | --help | --version)

Options:
    --all-targets      Build all possible targets
    --targets TARGETS  Build the given targets (comma seperated list).

    --release  Creates a tagged release and uploads the current artifacts.
               - Must run in a Github action (push event).
               - Pushing a tag with a version number (like 'v0.0.0') will
                 create a production release with that version number.
               - Updating a HEAD ref creates a release named 'latest'.

    --init     Create a template Github Action.

    -h --help  Show help.
    --version  Show version.`
