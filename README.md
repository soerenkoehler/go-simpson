Simpson
=======

*Simpson* is a **B**uild **A**nd **R**elease **T**ool for Go on Github. It runs
in a shell script task of a Github Action.

Features
--------

*   local and CI builds from the same tool
*   automatic latest release
*   manually triggered tagged releases
*   multi platform builds
*   creating archives depending on target platform:
    *   Windows: zip
    *   Linux: tgz
*   create sha256-digests of generated artifacts

Usage
-----

### Invocation and Module Aware Mode ###

You run *Simpson* like this:

```
go run github.com/soerenkoehler/simpson@main [MAINPACKAGE] [OPTIONS...]
```

You may replace `@main` with a different version if desired. But you must
provide a version to run in module aware mode (ref. [Documentation of `go
run`][go-docs-run]).

### Local Build ###

Just run *Simpson* **without** the option `--init`. Then the following Go
commands will run:

*   `go vet` & `go test` for all packages in your module
*   `go build` for the given main package

*Simpson* then creates a directory `artifacts` with:

*   artifact directories and archive files for all specified targets
*   a SHA256 file with checksums for all archives

### Github Build ###

1.  Prepare the workflow file
    
    Run *Simpson* locally **with** option `--init`:

    This will create a Github workflow file which will call *Simpson* in a
    Github action to build and release your module.

    The current workflow template file runs Go with `GOPROXY=direct`. If you do
    not want this, feel free to edit the workflow file.

2.  Push a branch containing the workflow file to Github. *Simpson* will run and
    create a release `latest`.

    **Warning:** The tag `latest` will change with every push. While fetching or
    pulling tags your local Git client will fail with the error:

    ```
    ! [rejected] latest -> latest (would clobber existing tag)
    ```

    On the commandline you must add the option `-f` to update the changed tags:
    ```
    git fetch --tags -f
    git pull --tags -f
    ```

    Your IDE may automagically try to pull or fetch tags in the background when
    synchronizing your local repository. You should switch off such behaviour
    and fetch tags manually:

    *   In VSCode change the setting `git.pullTags` to `false`.
    *   In other IDEs consult your friendly IDE manual.

3.  Push a tag in the semver format `v<MAJOR>.<MINOR>.<PATCH>` (or create it in
    the Github UI). *Simpson* will run and create a release for this tag.

### Options ###

*   If MAINPACKAGE is omitted the current directory `.` will be used.

*   `--artifact-name NAME` Changes the package name part of the created
    artifacts.

*   `--targets TARGET-SPEC, ...` Specifies the target platforms to build.
    Currently supported are:

    *   windows-amd64
    *   linux-amd64
    *   linux-arm64
    *   linux-arm

*   `--skip-upload` Suppresses the artifact upload (e.g. if you build a library
     which will be imported rather than an application executable).

*   `--init` Creates (or replaces) a Github workflow file with the given options
    except `--init` itself.

### Artifact Naming ###

Artifacts will be named as follows:

Release Type  | Artifact Name
--------------|-------------------------------------
Local         | MAINPACKAGE-DATE-TIME-PLATFORM
Github Latest | MAINPACKAGE-DATE-TIME-HASH-PLATFORM
Github Named  | MAINPACKAGE-VERSION

*Simpson* also injects the version string as `main._Version` into the build.

[go-docs-run]: https://pkg.go.dev/cmd/go#hdr-Compile_and_run_Go_program
