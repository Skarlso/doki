package cmd

import (
	"github.com/Masterminds/semver/v3"
)

var (
	currentVersion *semver.Version
)

// SetVersion must be called at bootstrap to pass the current build version
func SetVersion(releaseVersion string) {
	currentVersion = semver.MustParse(releaseVersion)
}
