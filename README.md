Simpson
=======

Simpson is a **B**uild **A**nd **R**elease **T**ool for Go on Github. It runs
in a shell script task of a Github Action.

Features
--------

*   local and CI builds from the same tool
*   multi platform builds
*   automatic latest release
*   manually triggered tagged releases
*   creating archives depending on target platform:
    * Windows: zip
    * Linux: tgz

Usage
-----

*   For a local build execute the following in your module workspace:
    ```
    go run github.com/soerenkoehler/simpson MAINPACKAGE [OPTIONS...]
    ```
    Since there is no Github Actions context, Simpson will skip release related
    tasks.

    The local build creates a directory `artifacts` where it places the build artifacts (as well in sub directories and compressed archives).

*   To prepare a workflow file for Github Actions add the option `--init`:
    ```
    go run github.com/soerenkoehler/simpson MAINPACKAGE [OPTIONS...] --init
    ```

*   `go run` should normally request the most recent release version. But if
    your Go module cache already contains an older (or wrong) version you may consider using `go get` to fetch the latest (or a specific) version.
    ```
    go get -u github.com/soerenkoehler/simpson@main
    go get -u github.com/soerenkoehler/simpson@<VERSION>
    ```

*   Running Simpson locally will modify your `go.mod` and `go.sum` files. So you
    should run `go mod tidy` from time to time.

*   The current workflow template file sets `GOPROXY=direct`. If you do not want
    this, feel free to edit the workflow file.

### Creating Named Releases ###

### General Oprions ###

```
--init
```

Creates a file `.github/workflows/simpson-build-and-release-tool.yml` with the
given options except `--init`.

### Build Options ###

```
--targets TARGET-SPECS
```

`TARGET-SPECS` is a comma delimited list of target platforms. Currently
supported are:

*   windows-amd64
*   linux-amd64
*   linux-arm64
*   linux-arm

```
--all-targets
```

Shortcut to build all supported targets.

### Release Options ###

```
--latest
```

TODO: handling changing tag `latest`

```
--skip-upload
```

Considerations
--------------

### Fetching Tags ###

The tag `latest` will change with every push. Fetching/pulling tags will thus
fail with the error:
```
! [rejected]        latest     -> latest  (would clobber existing tag)
```

On the commandline you must add the option `-f`:
```
git fetch --tags -f
git pull --tags -f
```

In VSCode you should change the setting `git.pullTags` to `false` and fetch tags
manually when required.
