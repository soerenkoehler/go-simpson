Simpson
=======

Simpson is a **B**uild **A**nd **R**elease **T**ool for Go on Github. It runs
in a shell script task of a Github Action.

Features
--------

* local and CI builds from the same tool
* multi platform builds
* automatic latest release
* manually triggered tagged releases
* creating archives depending on target platform:
  * Windows: zip
  * Linux: tgz
* create sha256-digests of generated release artifacts

Usage
-----

### Local ###

* For a local build execute the following in your module workspace:
  ```
  go run github.com/soerenkoehler/simpson [MAINPACKAGE] [OPTIONS...]
  ```

  The local build creates a directory `artifacts` where it places the build
  artifacts (as well in sub directories and compressed archives). The artifacts
  names will contain date, time and target platform.

  Since there is no Github Actions context, Simpson will skip release related
  tasks.

* To prepare a workflow file for Github Actions add the option `--init`:
  ```
  go run github.com/soerenkoehler/simpson [MAINPACKAGE] [OPTIONS...] --init
  ```

* `go run` should normally request the most recent release version. But if your
  Go module cache already contains an older (or wrong) version you may consider
  using `go get` to fetch the latest (or a specific) version.
  ```
  go get -u github.com/soerenkoehler/simpson@main
  go get -u github.com/soerenkoehler/simpson@<VERSION>
  ```

* Running Simpson will modify your `go.mod` and `go.sum` files. So you should
  run `go mod tidy` from time to time.

### On Github ###

Once you have created and pushed a workflow file, Simpson will build releases on
pushes to Github:

* When using option `--latest`, updating branch heads creates a _latest
  release_.
* Pushing tags with a semantic version number `v<MAJOR>.<MINOR.<PATCH>` creates
  a release with that version number. Conveniently, this happens also, when
  manually creating a release in the Github webapp.

### Artifact Version and Naming ###

If MAINPACKAGE is not provided the basename of the current directory `./.` will
be used.

Release Type    | Artifact Name
----------------|-------------------------------------
Local           | MAINPACKAGE-DATE-TIME-PLATFORM
Workflow Latest | MAINPACKAGE-DATE-TIME-HASH-PLATFORM
Workflow Named  | MAINPACKAGE-VERSION

Simpson also injects the version string as `main._Version` into the build.

### General Options ###

```
-i, --init
```

Creates (or replaces) the file
`.github/workflows/simpson-build-and-release-tool.yml` with the given Simpson
options except `--init`.

The current workflow template file runs Go with `GOPROXY=direct`. If you do not
want this, feel free to edit the workflow file.

### Build Options ###

```
-t, --targets TARGET-SPEC,...
```

`TARGET-SPEC,...` is a comma delimited list of target platforms. Currently
supported are:

* windows-amd64
* linux-amd64
* linux-arm64
* linux-arm

```
-a, --all-targets
```

Shortcut to build all supported targets.

### Release Options ###

```
-l, --latest
```

When this option is given and the Github Action is triggered by a push onto a
branch head, Simpson will tag, build and create a release named `latest`.

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
synchronizing your local repository. When using the option `--latest` you should
switch off such behaviour and fetch tags manually when required:

* In VSCode change the setting `git.pullTags` to `false`.
* In other IDEs consult your friendly IDE manual.

```
--skip-upload
```

Add this option if you want to run `go vet`, `go test` and `go build` but you do
not need the created binary artifacts. When you use `--skip-upload` you may omit
both `--targets` and `--all-targets`
