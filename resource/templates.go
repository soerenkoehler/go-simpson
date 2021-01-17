package resource

// Usage : Description of arguments and options.
const Usage = `BaRT (Build and Release Tool) ${VERSION}

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

    --init  Creates a template Github Action.

    -h --help  Show help.
    --version  Show version.`

// WorkflowFile : Template File for Github Actions.
const WorkflowFile = `name: Build And Release Tool

on:
  push

jobs:
  build-and-release:
    name: Build And Release
    runs-on: ubuntu-latest

    steps:
      - name: Setup Go environment
        uses: actions/setup-go@v2
        with:
          go-version: 1.15

      - name: Checkout
        uses: actions/checkout@v2

      - name: Setup Environment
        run: |
          echo 'GITHUB_CONTEXT<<EOF' >>$GITHUB_ENV
          echo '${{toJson(github)}}' >>$GITHUB_ENV
          echo 'EOF' >>$GITHUB_ENV
          echo 'GOPROXY=direct' >>$GITHUB_ENV

      - name: Build
        run: |
          go get ${MODULE}@dev
          go run ${MODULE} . --all-targets --release`
