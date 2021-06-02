# Dōki

A companion tool to work with versioning and keeping in sync [profiles](https://github.com/weaveworks/profiles) and [pctl](https://github.com/weaveworks/pctl).

<!--
To update the TOC, install https://github.com/kubernetes-sigs/mdtoc
and run: mdtoc -inplace README.md
-->

<!-- toc -->
- [Features](#features)
  - [go mod](#go-mod)
    - [latest](#latest)
    - [replace](#replace)
  - [dev tags](#dev-tags)
  - [self update-check](#self-update-check)
  - [optional Token](#optional-token)
- [Development](#development)
- [Releasing](#releasing)
<!-- /toc -->

## Features

### go mod

Dōki, wraps some `go mod` functionality for convenience and for usage in Makefiles.

#### latest
 
It provides a list of `go get`-able urls for a list of modules, by adding the latest available tag to that module.

Example:
```console
doki go mod latest github.com/weaveworks/pctl github.com/weaveworks/profiles
github.com/weaveworks/pctl@v0.0.3
github.com/weaveworks/profiles@v0.0.5
```

This output is best consumed by a `make` target such as:

```Makefile
.PHONY: update-modules
	go get \
		$(shell doki mod latest \
			github.com/weaveworks/profiles \
			github.com/weaveworks/pctl \
			<whatever> \
		)
	go mod tidy
```

Calling this `make` target will result in the latest pins in `go.mod` file to the list of these modules.

#### replace

Similar to `latest`, it provides a list of `replace` statements for modules.

Example:

```console
doki go mod replace --replacements github.com/weaveworks/pctl=github.com/weaveworks/pctl@specific-version,github.com/weaveworks/profiles=github.com/weaveworks/profiles@specific-version
-replace github.com/weaveworks/pctl@<specific-version> -replace github.com/weaveworks/profiles@<specific-version>
```

This output is best consumed by another make target, or the same `update-modules` target as above to make sure
all the modules always have the correct replacements.

```Makefile
.PHONY: update-modules
	go mod edit \
   		$(shell doki mod replacements)
	go get \
		$(shell doki mod latest \
			github.com/weaveworks/profiles \
			github.com/weaveworks/pctl \
			<whatever> \
		)
	go mod tidy
```

### dev tags

Convenient check to retrieve a dev tag for a branch.

Run:

```console
➜  doki git:(my-new-branch) ✗ doki get dev tag
v0.0.1-my-new-branch
```

It will automatically fetch the latest released tag and append the branch to it. This is the format used by `profiles`.
These tags are pushed into on new commits so they are always up-to-date.

### self update-check

This is a convenient function to check if your installation of Doki is up-to-date.

### Optional Token

If Dōki is used to access a non-public repository, calls to determine the latest version might require a token.
For that, it provides two options:

- through `DOKI_TOKEN`
- flag `--token`

## Development

Running tests: `make tests`.

## Releasing

There are some manual steps right now, should be streamlined soon.

Steps:

1. Create a new release notes file:
   ```sh
   touch docs/release_notes/<version>.md
   ```

1. Copy-and paste the release notes from the draft on the releases page into this file.
   _Note: sometimes the release drafter is a bit of a pain, verify that the notes are
   correct by doing something like: `git log --first-parent tag1..tag2`._

1. PR the release notes into main.

1. Create and push a tag with the new version:
   ```sh
   git tag <version>
   git push origin <version>
   ```

1. The `Create release` action should run. Verify that:
1. The release has been created in Github
   1. With the correct assets
   1. With the correct release notes
1. The image has been pushed to docker
1. The image can be pulled and used in a deployment

_Note_ that `<version>` must be in the following format: `v0.0.1`. 
