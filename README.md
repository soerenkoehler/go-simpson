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
*   creating zip and tgz archives depending on target platform

Usage
-----

### Downloading ###

TODO `GOPROXY=direct`

```
go get github.com/soerenkoehler/simpson@main
```

```
go get github.com/soerenkoehler/simpson@<VERSION>
```

### Running ###

```
go run github.com/soerenkoehler/simpson MAINPACKAGE [OPTIONS...]
```

### General Oprions ###

```
--init
```

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
