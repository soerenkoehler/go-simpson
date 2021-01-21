package resource

// Usage : Description of arguments and options.
const Usage = `BaRT (Build and Release Tool) ${VERSION}

Usage:
    simpson PACKAGE (--all-targets | --targets TARGETS) [--release] [--latest]
    simpson --init
    simpson (-h | --help | --version)

Options:
    --all-targets      Build all possible targets
    --targets TARGETS  Build the given targets (comma seperated list).

    --latest   Tags the latest commit and creates a release named 'latest'.
    --release  Creates a named releases when pushing a tag 'vX.Y.Z'.

    --init  Creates a template Github Action.

    -h --help  Show help.
    --version  Show version.

See detailed documentation under: https://github.com/soerenkoehler/simpson/
`

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
