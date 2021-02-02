package resource

// Usage : Description of arguments and options.
const Usage = `BaRT (Build and Release Tool) ${VERSION}

Usage:
    simpson PACKAGE (--all-targets | --targets TARGETS) [--latest] [--skip-upload] [--init]
    simpson (-h | --help | --version)

Options:
    --all-targets      Build all possible targets
    --targets TARGETS  Build the given targets (comma seperated list).

    --latest       Tags the latest commit and creates a release named 'latest'.

    --skip-upload  Build artifacts but do not upload them to the release.

    --init         Creates a Github Action file using the current commandline.

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
          go-version: ${SIMPSON_GOVERSION}

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
          go get github.com/soerenkoehler/simpson@main
          go run github.com/soerenkoehler/simpson ${SIMPSON_CMDLINE}
`
