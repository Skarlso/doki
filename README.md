# Dōki

A companion tool to work with versioning and keeping in sync [profiles](https://github.com/weaveworks/profiles) and [pctl](https://github.com/weaveworks/pctl).

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

## Development

### Optional Token

## Releasing
